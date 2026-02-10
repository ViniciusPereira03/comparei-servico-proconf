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

