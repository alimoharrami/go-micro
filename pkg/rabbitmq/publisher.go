package rabbitmq

import (
	"log"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/rabbitmq/amqp091-go"
	uuid "github.com/satori/go.uuid"
)

//go:generate mockery --name IPublisher
type IPublisher interface {
	PublishMessage(queue string, msg any) error
}

type Publisher struct {
	conn *amqp091.Connection
}

func (p Publisher) PublishMessage(queue string, msg any) error {

	data, err := jsoniter.Marshal(msg)

	if err != nil {
		log.Fatal("Error in marshalling message to publish message")
		return err
	}

	channel, err := p.conn.Channel()
	if err != nil {
		log.Fatal("Error in opening channel to consume message")
		return err
	}

	defer channel.Close()

	publishingMsg := amqp091.Publishing{
		Body:         data,
		ContentType:  "application/json",
		DeliveryMode: amqp091.Persistent,
		MessageId:    uuid.NewV4().String(),
		Timestamp:    time.Now(),
	}

	err = channel.Publish("", queue, false, false, publishingMsg)

	if err != nil {
		log.Println("Error in publishing message")
		return err
	}

	log.Printf("Published message: %s", publishingMsg.Body)

	return nil
}

func NewPublisher(conn *amqp091.Connection) IPublisher {
	return &Publisher{conn: conn}
}
