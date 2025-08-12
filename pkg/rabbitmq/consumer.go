package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

//go:generate mockery --name IConsumer
type IConsumer interface {
	ConsumeMessage(queue string, msg string) (*string, error)
}

type Consumer struct {
	conn *amqp091.Connection
}

func (c Consumer) ConsumeMessage(queue string, msg string) (*string, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		log.Printf("Error in opening channel to consume message")
		return nil, err
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue, // queue name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	if err != nil {
		return nil, err
	}

	message := ""

	for msg := range msgs {
		message += string(msg.Body)
	}

	return &message, nil
}

func NewConsumer(conn *amqp091.Connection) IConsumer {
	return &Consumer{conn: conn}
}
