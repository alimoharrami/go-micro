package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

func NewRabbitMQConn(cfg *RabbitMQConfig, ctx context.Context) (*amqp091.Connection, error) {
	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second // Maximum time to retry
	maxRetries := 5                      // Number of retries (including the initial attempt)

	var conn *amqp091.Connection
	var err error

	err = backoff.Retry(func() error {

		conn, err = amqp091.Dial(connAddr)
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %v. Connection information: %s", err, connAddr)
			return err
		}

		return nil
	}, backoff.WithMaxRetries(bo, uint64(maxRetries-1)))

	log.Println("Connected to RabbitMQ")

	go func() {
		select {
		case <-ctx.Done():
			err := conn.Close()
			if err != nil {
				log.Println("Failed to close RabbitMQ connection")
			}
			log.Println("RabbitMQ connection is closed")
		}
	}()

	return conn, err
}
