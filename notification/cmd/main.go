package main

import (
	"context"
	"fmt"
	"log"
	"notification/external/protos/userpb"
	"notification/internal/database"
	"notification/internal/handlers"
	"notification/internal/helpers"
	"notification/internal/migrations"
	"notification/internal/repository"
	"notification/internal/service"
	"os"
	"os/signal"
	"time"

	"notification/internal/config"
	"notification/internal/routes"
	"notification/internal/server"

	"github.com/alimoharrami/go-micro/pkg/rabbitmq"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Config struct {
	Rabbit *amqp.Connection
}

type AppParams struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Cfg        *config.Config
	DB         *gorm.DB
	Client     userpb.UserServiceClient
	EmailSvc   *service.EmailService
	RabbitCfg  rabbitmq.RabbitMQConfig
	RabbitCon  *helpers.RabbitConsumer
	RabbitConn *amqp.Connection
	Router     *gin.Engine
}

func NewRabbitConfig(cfg *config.Config) rabbitmq.RabbitMQConfig {
	return rabbitmq.RabbitMQConfig{
		Host:     "rabbitmq",
		Port:     5672,
		User:     "guest",
		Password: "guest",
	}
}

func NewRabbitConn(cfg rabbitmq.RabbitMQConfig, lc fx.Lifecycle) *amqp.Connection {
	ctx := context.Background()
	counter := 0
	for {
		conn, err := rabbitmq.NewRabbitMQConn(&cfg, ctx)
		if err != nil {
			counter++
		} else {
			return conn
		}
		if counter > 5 {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					return conn.Close()
				},
			})
			return nil
		}
	}
}

func StartApp(p AppParams) {

	migrations.AutoMigrate(p.DB)

	connection := p.RabbitConn

	p.RabbitCon.ConsumeMessage(context.Background(), connection, "notification")

	//Initialize Redis
	// redisClient := database.GetRedis()
	// defer redisClient.Close()

	srv := server.NewServer(p.Router)
	port := p.Cfg.Server.Port

	go func() {
		if err := srv.Start(port); err != nil {
			log.Fatalf("Failed to start service: %v", err)
		}
	}()
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	fmt.Println("Server gracefully stopped")
}

func main() {
	fx.New(
		fx.Provide(
			config.LoadConfig,
			database.NewPostgresConfig,
			database.InitPostgres,
			helpers.InitGRPC,
			service.NewEmailService,
			helpers.NewRabbitConsumer,
			NewRabbitConfig,
			NewRabbitConn,
			repository.NewChannelRepository,
			service.NewChannelService,
			handlers.NewNotificationController,
			routes.SetRouter,
		),
		fx.Invoke(StartApp),
	).Run()
}
