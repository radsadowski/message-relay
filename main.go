package main

import (
	"math/rand"
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

	db.AutoMigrate(pkg.SQLMessage{})

	for {
		t := time.Unix(rand.Int63n(1708162571), 0)
		tx := db.Create(&pkg.SQLMessage{Payload: "Hello, world!", CreatedAt: t, IsSent: false})
		if tx.Error != nil {
			panic(tx.Error)
		}
	}

	// brokerURL, ok := os.LookupEnv("BROKER_URL")
	// if !ok {
	// 	panic("BROKER_URL is not set")
	// }
	// conn, err := amqp.Dial(brokerURL)
	// if err != nil {
	// 	panic(err)
	// }
	// ch, err := conn.Channel()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// queue, err = ch.QueueDeclare(
	// 	"hello", // name
	// 	false,   // durable
	// 	false,   // delete when unused
	// 	false,   // exclusive
	// 	false,   // no-wait
	// 	nil,     // arguments
	// )
	// if err != nil {
	// 	panic(err)
	// }
}

func main() {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	//
	// relay := pkg.NewRelay(&pkg.MySQLMessageOutboxRepository{DB: db}, &pkg.RabbitMQPublisher{
	// 	Channel: ch,
	// 	Queue:   queue,
	// })
	// for {
	// 	log.Println("Relaying message")
	// 	err := relay.RelayMessage(ctx)
	// 	failOnError(err, "Failed to relay message")
	// 	time.Sleep(100 * time.Millisecond)
	// }
}
