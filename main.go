package main

import (
	"comparei-servico-promer/config"
	"comparei-servico-promer/internal/app"
	"comparei-servico-promer/internal/infrastructure/messaging/subscriber"
	"comparei-servico-promer/internal/infrastructure/repository"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Erro ao carregar configurações:", err)
	}

	// Testar conexão com Redis de mensageria
	redisMessaging := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_MESSAGING_HOST") + ":" + os.Getenv("REDIS_MESSAGING_PORT"),
	})
	ctx := context.Background()
	_, err := redisMessaging.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Não foi possível conectar ao Redis de mensageria:", err)
	}

	// --- MongoDB ---
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("Erro ao criar cliente MongoDB:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mongoClient.Connect(ctx); err != nil {
		log.Fatal("Erro ao conectar no MongoDB:", err)
	}
	// opcional: ping para certificar
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatal("Ping no MongoDB falhou:", err)
	}

	userRepo := repository.NewUserRepository(
		mongoClient,
		os.Getenv("MONGO_DB_NAME"),
		"usuarios",
	)

	userService := app.NewUserService(userRepo)

	subscriber.SetUserService(userService)

	// Iniciar o subscriber (rodar ouvindo eventos)
	go func() {
		fmt.Println("📡 Inicializando subscriber...")
		subscriber.Run()
	}()

	// Aguardar sinal de término
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("Encerrando serviço...")
}
