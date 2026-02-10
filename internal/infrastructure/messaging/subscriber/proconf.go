package subscriber

import (
	"comparei-servico-proconf/config"
	"comparei-servico-proconf/internal/app"
	"comparei-servico-proconf/internal/domain/proconf"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var proconf_service *app.ProconfService

func init() {
	config.LoadConfig()
	host := fmt.Sprintf("%v:%v", os.Getenv("REDIS_MESSAGING_HOST"), os.Getenv("REDIS_MESSAGING_PORT"))
	rdb = redis.NewClient(&redis.Options{
		Addr: host,
	})
}

func SetProconfService(service *app.ProconfService) {
	proconf_service = service
}

func subCreateMercadoProdutos() error {
	log.Println("EXEC: subCreateMercadoProdutos")
	ctx := context.Background()

	sub := rdb.Subscribe(ctx, "new_product")
	ch := sub.Channel()

	for msg := range ch {
		var promer proconf.ProconfConfirmValue
		err := json.Unmarshal([]byte(msg.Payload), &promer)
		if err != nil {
			fmt.Println("[ERRO] Erro ao decodificar payload de mensageria:", err)
			continue
		}

		promerParsed := promer.ParseToProconf()

		_, err_create := proconf_service.Create(promerParsed)
		if err_create != nil {
			fmt.Println("[ERRO] Erro ao criar proconf:", err_create)
		}
	}

	return nil
}
