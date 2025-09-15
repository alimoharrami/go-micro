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

func StartApp(p AppParams) {

	migrations.AutoMigrate(p.DB)

	connection := p.RabbitConn

	p.RabbitCon.ConsumeMessage(context.Background(), connection, "notification")

	//Initialize Redis
	// redisClient := database.GetRedis()
	// defer redisClient.Close()

	// p.Router.GET("/ws", wsHandler)

	// Start broadcaster
	// go handleBroadcast()

	// Send a notification every 5 seconds
	// go func() {
	// 	for {
	// 		broadcast <- fmt.Sprintf("Notification at %s", time.Now().Format(time.RFC1123))
	// 		time.Sleep(5 * time.Second)
	// 	}
	// }()

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
			hub.NewNotificationHelper,
			NewRabbitConfig,
			NewRabbitConn,
			repository.NewChannelRepository,
			service.NewChannelService,
			handlers.NewNotificationController,
			handlers.NewNotificationBoradcastHandler,
			routes.SetRouter,
		),
		fx.Invoke(StartApp),
	).Run()
}
