package pkg

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/radekska/outbox/pkg/outbox"
)

type RabbitMQPublisher struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

// Publish is responsible for publishing the message to the message broker
func (r *RabbitMQPublisher) Publish(ctx context.Context, msg *outbox.Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return r.Channel.PublishWithContext(ctx, "", r.Queue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        b,
	})
}
