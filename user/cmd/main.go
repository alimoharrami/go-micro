package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user/internal/database"
	"user/internal/helpers"
	"user/migrations"

	"user/internal/config"
	"user/internal/routes"
	"user/internal/server"

	"github.com/alimoharrami/go-micro/pkg/rabbitmq"
)

type NotificationData struct {
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}

type Notification struct {
	Type string           `json:"type"`
	Data NotificationData `json:"data"`
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println(cfg.Server.Port)

	ctx := context.Background()

	rabbitCfg := rabbitmq.RabbitMQConfig{
		Host:     "rabbitmq",
		Port:     5672,
		User:     "guest",
		Password: "guest",
	}
	rabbitconn, err := rabbitmq.NewRabbitMQConn(&rabbitCfg, ctx)

	if err != nil {
		log.Printf("Error connecting rabbitmq %v:", err)
	}

	payloadData := NotificationData{
		UserID:  1,
		Message: "this is message",
	}
	payload := Notification{
		Type: "user_notif",
		Data: payloadData,
	}

	if err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}

	rabbitPublisher := rabbitmq.NewPublisher(rabbitconn)
	rabbitPublisher.PublishMessage("notification", payload)

	// init dbs
	_ = database.InitDatabases(database.NewPostgresConfig(), database.RedisConfig(cfg.Redis))
	db := database.GetPostgres()
	sqlDb, err := db.DB()

	if err != nil {
		log.Fatalf("Failed to get DB connection: %v", err)
	}

	migrations.AutoMigrate(db)

	defer sqlDb.Close()

	//Initialize Redis
	// redisClient := database.GetRedis()
	// defer redisClient.Close()

	helpers.InitilaizeGRPC(db)

	log.Println("its in creating router")
	router := routes.SetRouter(db)

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
