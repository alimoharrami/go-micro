package main

import (
	"context"
	"log"
	"net/http"
	"notification/internal/database"
	"notification/internal/handlers"
	"notification/internal/helpers"
	"notification/internal/hub"
	"notification/internal/migrations"
	"notification/internal/repository"
	"notification/internal/service"

	"notification/internal/config"
	"notification/internal/routes"
	"notification/internal/server"

	"github.com/alimoharrami/go-micro/pkg/rabbitmq"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

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

// RunMigrations performs database migrations
func RunMigrations(db *gorm.DB) {
	migrations.AutoMigrate(db)
}

func RegisterHooks(
	lc fx.Lifecycle,
	cfg *config.Config,
	db *gorm.DB,
	rabbitConn *amqp.Connection,
	rabbitCon *helpers.RabbitConsumer,
	router *gin.Engine,
) {
	srv := server.NewServer(router)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Starting Notification Microservice on port %s", cfg.Server.Port)

			// Start RabbitMQ consumer
			go rabbitCon.ConsumeMessage(context.Background(), rabbitConn, "notification")

			// Start HTTP server
			go func() {
				if err := srv.Start(cfg.Server.Port); err != nil && err != http.ErrServerClosed {
					log.Printf("Server start error: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down Notification Microservice...")

			// RabbitMQ connection is closed by its own hook in NewRabbitConn if needed,
			// or we can explicitly close it here.
			if rabbitConn != nil {
				_ = rabbitConn.Close()
			}

			return srv.Shutdown(ctx)
		},
	})
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
			hub.NewNotificationHelper,
			NewRabbitConfig,
			NewRabbitConn,
			repository.NewChannelRepository,
			service.NewChannelService,
			handlers.NewNotificationController,
			handlers.NewNotificationBroadcastHandler,
			routes.SetRouter,
		),
	fx.Invoke(
		RunMigrations,
		RegisterHooks,
		func(
			lc fx.Lifecycle,
			consumer *helpers.RabbitConsumer,
			conn *amqp.Connection,
		) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go consumer.ConsumeMessage(ctx, conn, "notification_queue")
					return nil
				},
			})
		},
	),
).Run()
}
