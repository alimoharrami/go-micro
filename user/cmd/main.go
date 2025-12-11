package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user/internal/config"
	"user/internal/database"
	"user/internal/helpers"
	"user/internal/logger"
	"user/internal/routes"
	"user/internal/server"
	"user/migrations"

	"github.com/alimoharrami/go-micro/pkg/rabbitmq"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize Logger
	logger.InitLogger(cfg.Server.Env)
	logger.Info("Starting User Microservice")
	logger.Info(fmt.Sprintf("Server port: %s", cfg.Server.Port))

	// Init DB
	_ = database.InitDatabases(database.NewPostgresConfig(), database.RedisConfig(cfg.Redis))
	db := database.GetPostgres() // db is *gorm.DB
	sqlDb, err := db.DB()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to get DB connection: %v", err))
	}
	defer sqlDb.Close()

	// Run Migrations
	migrations.AutoMigrate(db)

	// Initialize RabbitMQ
	rabbitCfg := &rabbitmq.RabbitMQConfig{
		Host:     cfg.RabbitMQ.Host,
		Port:     cfg.RabbitMQ.Port,
		User:     cfg.RabbitMQ.User,
		Password: cfg.RabbitMQ.Password,
	}
	rabbitConn, err := rabbitmq.NewRabbitMQConn(rabbitCfg, context.Background())
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to connect to RabbitMQ: %v", err))
	}
	defer rabbitConn.Close()
	
	publisher := rabbitmq.NewPublisher(rabbitConn)

	// Initialize GRPC Helper
	helpers.InitilaizeGRPC(db, publisher)

	logger.Info("Initializing router")
	router := routes.SetRouter(db, publisher)

	srv := server.NewServer(router)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.Start(cfg.Server.Port); err != nil {
			logger.Fatal(fmt.Sprintf("Failed to start service: %v", err))
		}
	}()

	<-quit
	logger.Info("Shutting down server...")

	// Create shutdown context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown services gracefully
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(fmt.Sprintf("Server shutdown failed: %v", err))
	}

	logger.Info("Server gracefully stopped")
}
