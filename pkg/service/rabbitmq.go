package service

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/radekska/outbox/pkg/relay"
	log "github.com/sirupsen/logrus"
)

// RabbitMQPublisher is a struct that represents a RabbitMQ publisher
type RabbitMQPublisher struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

// NewRabbitMQPublisher is a constructor for the RabbitMQPublisher
func NewRabbitMQPublisher(brokerURL string) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(brokerURL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQPublisher{
		Channel: ch,
		Queue:   queue,
	}, nil
}

// Publish is responsible for publishing the message to the message broker
func (r *RabbitMQPublisher) Publish(ctx context.Context, msg *relay.Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return r.Channel.PublishWithContext(ctx, "", r.Queue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        b,
	})
}

// RabbitMQRetrier is a struct that represents a RabbitMQ retrier
type RabbitMQRetrier struct {
	logger *log.Logger
}

// Retry is a method that retries the RabbitMQ connection
func (r *RabbitMQRetrier) Retry(_ context.Context, err error) error {
	e, ok := err.(*amqp.Error)
	if !ok {
		return nil
	}
	r.logger.Errorf("broker network error: %v", e)
	return nil
}
