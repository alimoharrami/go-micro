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
	"user/migrations"

	"user/internal/config"
	"user/internal/helpers"
	"user/internal/routes"
	"user/internal/server"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	connectRabbitMQ()

	log.Println(cfg.Server.Port)

	helpers.InitilaizeGRPC()

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

func connectRabbitMQ() {
	conn, err := amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"notification", // queue name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	body := "Notify this user!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	log.Println("Sent notification message")
}
