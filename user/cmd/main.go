package main

import (
	"context"
	"fmt"
	"net/http"
	"user/internal/config"
	"user/internal/database"
	"user/internal/handlers"
	"user/internal/helpers"
	"user/internal/logger"
	"user/internal/repository"
	"user/internal/routes"
	"user/internal/server"
	"user/internal/service"
	"user/migrations"

	"github.com/alimoharrami/go-micro/pkg/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	fx.New(
		fx.Provide(
			// Config
			config.LoadConfig,

			// Logger
			func(cfg *config.Config) *zap.Logger {
				logger.InitLogger(cfg.Server.Env)
				return logger.Log
			},

			// Database Config
			func(cfg *config.Config) database.PostgresConfig {
				return database.PostgresConfig{
					Host:     cfg.Database.Host,
					Port:     cfg.Database.Port,
					User:     cfg.Database.User,
					Password: cfg.Database.Password, 
					DBName:   cfg.Database.Name,
					SSLMode:  cfg.Database.SSLMode,
				}
			},

			// Database Connection
			database.NewPostgresConnection,

			// RabbitMQ Config
			func(cfg *config.Config) *rabbitmq.RabbitMQConfig {
				return &rabbitmq.RabbitMQConfig{
					Host:     cfg.RabbitMQ.Host,
					Port:     cfg.RabbitMQ.Port,
					User:     cfg.RabbitMQ.User,
					Password: cfg.RabbitMQ.Password,
				}
			},

			// RabbitMQ Connection
			func(cfg *rabbitmq.RabbitMQConfig) (*amqp091.Connection, error) {
				return rabbitmq.NewRabbitMQConn(cfg, context.Background())
			},

			// RabbitMQ Publisher
			func(conn *amqp091.Connection) rabbitmq.IPublisher {
				return rabbitmq.NewPublisher(conn)
			},

			// Layers
			repository.NewUserRepository,
			service.NewUserService,
			handlers.NewUserController,
			routes.NewRouter,
			server.NewServer,
		),
		fx.Invoke(
			// Migrations
			func(db *gorm.DB) {
				migrations.AutoMigrate(db)
			},

			// Legacy GRPC Helper
			helpers.InitilaizeGRPC,

			// Lifecycle Management
			func(lifecycle fx.Lifecycle, srv *server.Server, cfg *config.Config, logger *zap.Logger, db *gorm.DB, rabbitConn *amqp091.Connection) {
				lifecycle.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						logger.Info("Starting User Microservice")
						logger.Info(fmt.Sprintf("Server port: %s", cfg.Server.Port))

						go func() {
							if err := srv.Start(cfg.Server.Port); err != nil && err != http.ErrServerClosed {
								logger.Fatal(fmt.Sprintf("Failed to start service: %v", err))
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						logger.Info("Shutting down server...")
						
						// Close DB
						sqlDb, err := db.DB()
						if err == nil {
							_ = sqlDb.Close()
						}

						// Close RabbitMQ
						_ = rabbitConn.Close()

						return srv.Shutdown(ctx)
					},
				})
			},
		),
	).Run()
}
