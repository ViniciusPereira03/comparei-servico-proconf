package repository

import (
	"comparei-servico-proconf/internal/domain/proconf"
	"context"
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProconfRepository struct {
	promer *mongo.Collection
	logs   *mongo.Collection
}

// Estrutura auxiliar para receber o resultado da agrupação do Mongo
type AggregationResult struct {
	MercadoProdutoID int       `bson:"_id"`
	LastUpdate       time.Time `bson:"last_update"`
	TotalUserPoints  int       `bson:"total_user_points"`
	PrecoMedio       float64   `bson:"preco_medio"`
}

func NewProconfRepository(client *mongo.Client, dbName, collectionNameProconf, collectionNameLogs string) *ProconfRepository {
	collProconf := client.Database(dbName).Collection(collectionNameProconf)
	collLogs := client.Database(dbName).Collection(collectionNameLogs)
	return &ProconfRepository{promer: collProconf, logs: collLogs}
}

func (r *ProconfRepository) Create(u *proconf.Proconf) (*proconf.Proconf, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u.CreatedAt = time.Now()
	u.ModifiedAt = time.Now()

	_, err := r.promer.InsertOne(ctx, u)
	if err != nil {
		log.Printf("Erro ao inserir Proconf: %v", err)
		return nil, err
	}
	return u, nil
}

func (r *ProconfRepository) GetMercadoProdutoByID(id int) (*proconf.Proconf, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"mercado_produto_id": id}
	var proconfEntry proconf.Proconf
	err := r.promer.FindOne(ctx, filter).Decode(&proconfEntry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Retorna nil se não encontrar
		}
		log.Printf("Erro ao buscar Proconf por ID: %v", err)
		return nil, err
	}
	return &proconfEntry, nil
}

// NOVO MÉTODO: Calcula apenas para um ID específico
func (r *ProconfRepository) CalculateConfidenceScoreForProduct(mercadoProdutoID int) error {
	// Filtro adicional para o ID específico
	matchFilter := bson.E{Key: "mercado_produto_id", Value: mercadoProdutoID}
	return r.calculateConfidence(true, matchFilter)
}

// MÉTODO REFATORADO: Calcula para todos (chamado pelo Cron)
func (r *ProconfRepository) CalculateConfidenceScores() error {
	// Sem filtro adicional de ID, processa tudo
	return r.calculateConfidence(false)
}

// Lógica central extraída e privada (Reutilizável)
// Aceita filtros opcionais (variadic) para adicionar ao $match
func (r *ProconfRepository) calculateConfidence(preventDecrease bool, additionalFilters ...bson.E) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	collectionLogs := r.logs
	collectionProconf := r.promer

	// Monta o filtro de data (últimos 30 dias)
	matchStage := bson.D{
		{Key: "created_at", Value: bson.D{{Key: "$gte", Value: time.Now().AddDate(0, 0, -30)}}},
	}

	// Se houver filtros adicionais (ex: ID específico), adiciona ao $match
	if len(additionalFilters) > 0 {
		matchStage = append(matchStage, additionalFilters...)
	}

	// 1. Pipeline de Agregação
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchStage}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$mercado_produto_id"},
			{Key: "last_update", Value: bson.D{{Key: "$max", Value: "$created_at"}}},
			{Key: "total_user_points", Value: bson.D{{Key: "$sum", Value: "$nivel_usuario"}}},
			{Key: "preco_medio", Value: bson.D{{Key: "$avg", Value: "$preco"}}},
		}}},
	}

	cursor, err := collectionLogs.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Erro ao agregar logs: %v", err)
		return err
	}
	defer cursor.Close(ctx)

	var operations []mongo.WriteModel

	for cursor.Next(ctx) {
		var result AggregationResult
		if err := cursor.Decode(&result); err != nil {
			log.Printf("Erro ao decodificar cursor: %v", err)
			continue
		}

		const pointsMultiplier = 3.0
		baseScore := float64(result.TotalUserPoints) * pointsMultiplier

		// Teto máximo de 100
		if baseScore > 100 {
			baseScore = 100
		}

		// 2. Aplica o Decaimento Temporal
		// A confiança calculada acima é válida para o momento "agora".
		// Quanto mais tempo passa, mais subtraímos desse valor.
		hoursSinceUpdate := time.Since(result.LastUpdate).Hours()

		// Taxa de perda: 0.5 pontos por hora (12 pontos por dia).
		// Se a confiança era 10 (apenas 1 usuário nível 1), zera em 20 horas.
		// Se a confiança era 100, zera em ~8 dias sem novas atualizações.
		const decayRatePerHour = 0.5
		decayPenalty := hoursSinceUpdate * decayRatePerHour

		finalScore := baseScore - decayPenalty

		// Limites (0 a 100)
		if finalScore < 0 {
			finalScore = 0
		}

		confiancaInt := int(math.Round(finalScore))

		// --- REGRA DE PROTEÇÃO (Subscriber) ---
		// Se preventDecrease for true (Execução por usuário), verificamos o valor atual no banco.
		// Se o novo valor for MENOR que o atual, pulamos a atualização ("nada pode ser feito").
		if preventDecrease {
			// Busca documento atual para comparar
			existingDoc, err := r.GetMercadoProdutoByID(result.MercadoProdutoID)

			// Se o documento existe e não houve erro...
			if err == nil && existingDoc != nil {
				if confiancaInt < existingDoc.NivelConfianca {
					log.Printf("Proteção de Confiança Ativada: ID %d. Calculado: %d < Atual: %d. Update ignorado.",
						result.MercadoProdutoID, confiancaInt, existingDoc.NivelConfianca)
					continue // Pula para a próxima iteração, ignorando o update
				}
			}
		}
		// ---------------------------------------

		// Prepara o Update
		filter := bson.M{"mercado_produto_id": result.MercadoProdutoID}
		update := bson.M{
			"$set": bson.M{
				"nivel_confianca": confiancaInt,
				"modified_at":     time.Now(),
				"last_log_at":     result.LastUpdate,
			},
			"$setOnInsert": bson.M{
				"created_at": time.Now(),
			},
		}

		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		operations = append(operations, model)

		if len(operations) >= 500 {
			_, err := collectionProconf.BulkWrite(ctx, operations)
			if err != nil {
				log.Printf("Erro no BulkWrite: %v", err)
			}
			operations = nil
		}
	}

	if len(operations) > 0 {
		_, err := collectionProconf.BulkWrite(ctx, operations)
		if err != nil {
			return err
		}
	}

	return nil
}
