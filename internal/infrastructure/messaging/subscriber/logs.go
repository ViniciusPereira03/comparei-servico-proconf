package subscriber

import (
	"comparei-servico-proconf/config"
	"comparei-servico-proconf/internal/app"
	"comparei-servico-proconf/internal/domain/logs"
	"comparei-servico-proconf/internal/domain/proconf"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	logs_confirmacao_service *app.LogsService
	proconf_service_subs     *app.ProconfService
)

func init() {
	config.LoadConfig()
	host := fmt.Sprintf("%v:%v", os.Getenv("REDIS_MESSAGING_HOST"), os.Getenv("REDIS_MESSAGING_PORT"))
	rdb = redis.NewClient(&redis.Options{
		Addr: host,
	})
}

// Função para injetar dependências (Adicione ou Atualize esta função)
func SetServices(logsService *app.LogsService, proconfService *app.ProconfService) {
	logs_confirmacao_service = logsService
	proconf_service_subs = proconfService
}

// Função para injetar o logs_confirmacao_service
func SetLogsConfirmacaoService(service *app.LogsService) {
	logs_confirmacao_service = service
}

func subCreateLogsConfirmacao() error {
	log.Println("EXEC: subCreateLogsConfirmacao")
	ctx := context.Background()

	logConf := rdb.Subscribe(ctx, "confirma_valor_mercado_produto")
	ch := logConf.Channel()

	type payload_log struct {
		Id     int    `json:"id"`
		UserID string `json:"user_id"`
	}

	for msg := range ch {
		// var logEntry logs.LogsConfirmacao
		var logEntry payload_log
		err := json.Unmarshal([]byte(msg.Payload), &logEntry)
		if err != nil {
			fmt.Println("[ERRO] Erro ao decodificar payload de mensageria:", err)
			continue
		}

		var logConfirmacao logs.LogsConfirmacao
		logConfirmacao.UsuarioID = logEntry.UserID
		logConfirmacao.MercadoProdutoID = logEntry.Id
		user, err := user_service.GetUser(logEntry.UserID)
		if err != nil {
			fmt.Println("[ERRO] Erro ao buscar usuário:", err)
			continue
		}
		logConfirmacao.NivelUsuario = user.Level

		// Buscar valor atual do produto para log
		var mercadoProduto *proconf.Proconf
		mercadoProduto, err = proconf_service.GetMercadoProdutoByID(logEntry.Id)
		if err != nil {
			fmt.Println("[ERRO] Erro ao buscar mercado produto:", err)
			continue
		}
		logConfirmacao.PrecoUnitario = mercadoProduto.PrecoUnitario

		logConfirm, err_create := logs_confirmacao_service.CreateLogsConfirmacao(&logConfirmacao)
		if err_create != nil {
			fmt.Println("[ERRO] Erro ao criar log de confirmação:", err_create)
			continue // Se falhou ao criar log, não recalculamos confiança
		}
		log.Println("[LOG] confirma_valor_mercado_produto criado com ID:", logConfirm.ID)

		// 3. NOVO: Recalcula a confiança IMEDIATAMENTE para este produto específico
		log.Printf("Recalculando confiança para produto ID: %d", logEntry.Id)
		err_recalc := proconf_service_subs.CalculateConfidenceScoreForProduct(logEntry.Id)
		if err_recalc != nil {
			fmt.Printf("[ERRO] Falha ao recalcular confiança para ID %d: %v\n", logEntry.Id, err_recalc)
		}

	}

	return nil
}
