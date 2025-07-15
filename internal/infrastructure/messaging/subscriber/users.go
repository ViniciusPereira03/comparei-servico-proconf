package subscriber

import (
	"comparei-servico-promer/config"
	"comparei-servico-promer/internal/app"
	"comparei-servico-promer/internal/domain/users"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var user_service *app.UserService

func init() {
	config.LoadConfig()
	host := fmt.Sprintf("%v:%v", os.Getenv("REDIS_MESSAGING_HOST"), os.Getenv("REDIS_MESSAGING_PORT"))
	rdb = redis.NewClient(&redis.Options{
		Addr: host,
	})
}

// Função para injetar o user_service
func SetUserService(service *app.UserService) {
	user_service = service
}

func Run() {
	go subCreateUser()
}

func subCreateUser() error {
	log.Println("EXEC: subCreateUser")
	ctx := context.Background()

	sub := rdb.Subscribe(ctx, "user_created")
	ch := sub.Channel()

	for msg := range ch {
		var user users.User
		err := json.Unmarshal([]byte(msg.Payload), &user)
		if err != nil {
			fmt.Println("[ERRO] Erro ao decodificar payload de mensageria:", err)
			continue
		}

		err_create := user_service.CreateUser(&user)
		if err_create != nil {
			fmt.Println("[ERRO] Erro ao criar user nos logs:", err_create)
		}
	}

	return nil
}
