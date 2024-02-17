package main

import (
	"context"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/radekska/outbox/pkg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db    *gorm.DB
	ch    *amqp.Channel
	queue amqp.Queue
)

func init() {
	dsn, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		panic("DATABASE_URL is not set")
	}

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	brokerURL, ok := os.LookupEnv("BROKER_URL")
	if !ok {
		panic("BROKER_URL is not set")
	}
	conn, err := amqp.Dial(brokerURL)
	if err != nil {
		panic(err)
	}
	ch, err = conn.Channel()
	if err != nil {
		panic(err)
	}

	queue, err = ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx := context.Background()
	// TODO: - add a defer statement to close the database connection and the RabbitMQ connection
	// TODO: - gracefully handle the termination of the program
	// TODO: - tests
	// TODO: - better logging
	// TODO: - metrics
	// TODO: - health checks
	// TODO: - add a retry mechanism in case of a failure
	// TODO: - add a way to configure the number of retries
	// TODO: - add a way to control the loop

	publisher, err := pkg.NewRabbitMQPublisher(ch, queue)
	if err != nil {
		log.Fatalf("error while creating publisher: %v", err)
	}
	relay := pkg.NewRelay(&pkg.MySQLMessageOutboxRepository{DB: db}, publisher)
	for {
		log.Println("Relaying message")
		err := relay.RelayMessage(ctx)
		if err != nil {
			log.Printf("error while relaying message: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
