package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// TODO: - add a defer statement to close the database connection and the RabbitMQ connection
	// TODO: - tests
	// TODO: - metrics
	// TODO: - health checks
	// TODO: - add a retry mechanism in case of a failure
	// TODO: - add a way to configure the number of retries

	publisher, err := pkg.NewRabbitMQPublisher(ch, queue)
	if err != nil {
		log.Fatalf("error while creating publisher: %v", err)
	}
	relay := pkg.NewRelay(&pkg.MySQLMessageOutboxRepository{DB: db}, publisher)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	for {
		select {
		case s := <-sigCh:
			log.Warnf("received signal: %v", s)
			log.Warn("quitting gracefully...")
			cancel()
			os.Exit(0)
		default:
			log.Debug("relaying message")
			err := relay.RelayMessage(ctx)
			if err != nil {
				log.Errorf("error while relaying message: %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
