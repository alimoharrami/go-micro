package helpers

import (
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ() {
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

	listenRabbitMQ(ch)

}

func ConsumeMessage(conn *amqp091.Connection, queue string) {
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

			var payload map[string]interface{}
			if err := json.Unmarshal(d.Body, &payload); err == nil {
				log.Printf("UserID: %v", payload["user_id"])
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

func listenRabbitMQ(ch *amqp091.Channel) {
	q, err := ch.QueueDeclare(
		"notification",
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
		q.Name,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Waiting for messages...")
	go func() {
		for d := range msgs {
			log.Println("Received a message:")
			log.Printf("Body: %s", d.Body)

			var payload map[string]interface{}
			if err := json.Unmarshal(d.Body, &payload); err == nil {
				log.Printf("UserID: %v, Message: %v", payload["userID"], payload["message"])
			}

			// manual ack
			if err := d.Ack(false); err != nil {
				log.Printf("Failed to ack message: %v", err)
			}
		}
	}()
}
