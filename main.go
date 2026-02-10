package main

import (
	"comparei-servico-proconf/config"
	"comparei-servico-proconf/internal/app"
	"comparei-servico-proconf/internal/infrastructure/messaging/subscriber"
	"comparei-servico-proconf/internal/infrastructure/repository"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
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

	logRepo := repository.NewLogsRepository(
		mongoClient,
		os.Getenv("MONGO_DB_NAME"),
		"logs_confirmacao",
	)

	userService := app.NewUserService(userRepo)
	logService := app.NewLogsService(logRepo)

	subscriber.SetUserService(userService)
	subscriber.SetLogsConfirmacaoService(logService)

	// Iniciar o subscriber (rodar ouvindo eventos)
	go func() {
		fmt.Println("📡 Inicializando subscriber...")
		subscriber.Run()
	}()

	// --- CRON ---
	crn := cron.New()

	// Função auxiliar para logar execuções de cron
	addCronJob := func(spec string, cmd func(), description string) {
		_, err := crn.AddFunc(spec, func() {
			log.Printf("CRON: Iniciando tarefa '%s' (%s)", description, spec)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("CRON: Panic recuperado na tarefa '%s': %v", description, r)
					}
				}()
				cmd()
				log.Printf("CRON: Tarefa '%s' concluída.", description)
				printMemoryAndGoroutineUsage()
			}()
		})
		if err != nil {
			log.Fatalf("Erro ao agendar tarefa cron '%s': %v", description, err)
		}
	}

	// Agendar tarefas cron
	addCronJob("0 */6 * * *", func() {
		proconfService.CalculateConfidenceScores()
	}, "A cada 6 horas")

	crn.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("Encerrando serviço...")
}
