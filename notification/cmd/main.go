package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"notification/external/protos/userpb"
	"notification/internal/database"
	"notification/internal/handlers"
	"notification/internal/helpers"
	"notification/internal/hub"
	"notification/internal/migrations"
	"notification/internal/repository"
	"notification/internal/service"
	"os"
	"os/signal"
	"sync"
	"time"

	"notification/internal/config"
	"notification/internal/routes"
	"notification/internal/server"

	"github.com/alimoharrami/go-micro/pkg/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

// Upgrade HTTP -> WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Store connected clients safely
var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex

// Broadcast channel
var broadcast = make(chan string)

func wsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Add client
	mu.Lock()
	clients[conn] = true
	mu.Unlock()
	fmt.Println("Client connected")

	// Remove on disconnect
	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		fmt.Println("Client disconnected")
	}()

	// Read messages from this client
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		// Send message to broadcast channel
		broadcast <- string(msg)
	}
}

func handleBroadcast() {
	for {
		msg := <-broadcast
		mu.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

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

			// Start broadcast loop
			go handleBroadcast()

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
			handlers.NewNotificationBoradcastHandler,
			routes.SetRouter,
		),
		fx.Invoke(
			RunMigrations,
			RegisterHooks,
		),
	).Run()
}
