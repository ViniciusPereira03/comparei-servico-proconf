package subscriber

import (
	"comparei-servico-proconf/config"
	"comparei-servico-proconf/internal/app"
	"comparei-servico-proconf/internal/domain/logs"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var logs_confirmacao_service *app.LogsService

func init() {
	config.LoadConfig()
	host := fmt.Sprintf("%v:%v", os.Getenv("REDIS_MESSAGING_HOST"), os.Getenv("REDIS_MESSAGING_PORT"))
	rdb = redis.NewClient(&redis.Options{
		Addr: host,
	})
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

		err_create := logs_confirmacao_service.CreateLogsConfirmacao(&logConfirmacao)
		if err_create != nil {
			fmt.Println("[ERRO] Erro ao criar log de confirmação:", err_create)
		}
	}

	return nil
}
