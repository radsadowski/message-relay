package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/radekska/outbox/pkg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dsn := "admin:secret123@tcp(127.0.0.1:3306)/outbox?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		failOnError(err, "Failed to connect to database")
	}
	err = db.AutoMigrate(&pkg.SQLMessage{})
	if err != nil {
		failOnError(err, "Failed to migrate database")
	}

	if db.Create(
		&pkg.SQLMessage{ID: 1, Payload: "Hello, World!", IsSent: false, CreatedAt: time.Now().Format(time.RFC3339)},
	).Error != nil {
		failOnError(err, "Failed to create message")
	}

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

	relay := pkg.NewRelay(&pkg.MySQLMessageOutboxRepository{DB: db}, &pkg.RabbitMQPublisher{
		Channel: ch,
		Queue:   q,
	})
	err = relay.RelayMessage(ctx)
	failOnError(err, "Failed to relay message")
}
