package main

import (
	"context"
	"fmt"
	"log"
	"notification/internal/database"
	"notification/internal/helpers"
	"notification/internal/migrations"
	"notification/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	"notification/internal/config"
	"notification/internal/routes"
	"notification/internal/server"

	"github.com/alimoharrami/go-micro/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println(cfg.Server.Port)

	ctx := context.Background()

	// init dbs
	_ = database.InitDatabases(database.NewPostgresConfig(), database.RedisConfig(cfg.Redis))
	db := database.GetPostgres()
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get DB connection: %v", err)
	}

	migrations.AutoMigrate(db)

	defer sqlDb.Close()

	//grpc service
	client := helpers.InitGRPC()

	EmailService := service.NewEmailService(client)

	rabbitCfg := rabbitmq.RabbitMQConfig{
		Host:     "rabbitmq",
		Port:     5672,
		User:     "guest",
		Password: "guest",
	}
	rabbitconn, err := rabbitmq.NewRabbitMQConn(&rabbitCfg, ctx)

	rabbitConsumer := helpers.NewRabbitConsumer(EmailService)

	if err != nil {
		log.Printf("Error connecting rabbitmq %v:", err)
	} else {
		rabbitConsumer.ConsumeMessage(ctx, rabbitconn, "notification")

	}
	//Initialize Redis
	// redisClient := database.GetRedis()
	// defer redisClient.Close()

	router := routes.SetRouter(db, client)

	srv := server.NewServer(router)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		fmt.Println("Shutting down server...")

		// Create shutdown context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Shutdown services gracefully
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}

		// redisClient.Close()
		sqlDb.Close()
		fmt.Println("Server gracefully stopped")
	}()

	//start server
	port := cfg.Server.Port
	if err := srv.Start(port); err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}
}
