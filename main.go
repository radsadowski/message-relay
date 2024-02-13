package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/radekska/outbox/pkg"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	publisher := &pkg.RabbitMQPublisher{
		Channel: ch,
		Queue:   q,
	}
	err = publisher.Publish(
		ctx,
		&pkg.Message{ID: "test-msg", Payload: "encoded-json", IsSent: false, CreatedAt: "2022-01-01 10:00:00"},
	)
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", "encoded-json")
	// start relay
	// connect to database
	// connect to message broker
	// consume message from outbox that are not sent
}
