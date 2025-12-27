package helpers

import (
	"context"
	"encoding/json"
	"log"
	"notification/internal/service"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitConsumer struct {
	EmailService   *service.EmailService
	ChannelService *service.ChannelService
}

type NotificationData struct {
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}

type Notification struct {
	Type string           `json:"type"`
	Data NotificationData `json:"data"`
}

func NewRabbitConsumer(EmailService *service.EmailService, ChannelService *service.ChannelService) *RabbitConsumer {
	return &RabbitConsumer{
		EmailService:   EmailService,
		ChannelService: ChannelService,
	}
}

func (c *RabbitConsumer) ConsumeMessage(ctx context.Context, conn *amqp091.Connection, queue string) {
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Error in opening channel to consume message")
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
		log.Printf("there is a problem in cunsuming this queue %v", err)
	}

	go func() {
		for d := range msgs {
			log.Println("Received a message:")
			log.Printf("Body: %s", d.Body)

			var payload Notification
			if err := json.Unmarshal(d.Body, &payload); err == nil {
				log.Printf("Type: %v", payload.Type)

				if payload.Type == "user_notif" {
					c.EmailService.SendEmailUserID(ctx, payload.Data.UserID, payload.Data.Message)
				} else if payload.Type == "blog_created" {
					log.Printf("Broadcasting blog creation notification: %s", payload.Data.Message)
					c.ChannelService.SendNotif(ctx, "main", payload.Data.Message)
				}
			} else {
				log.Printf("there is a problem here %v", err)
			}

			// manual ack
			if err := d.Ack(false); err != nil {
				log.Printf("Failed to ack message: %v", err)
			}
		}
	}()

}
