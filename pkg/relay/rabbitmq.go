package relay

import "github.com/radekska/outbox/pkg/outbox"

type RabbitMQMessageRelay struct {
	outboxRepository outbox.MessageOutboxRepository
}

// Consume is responsible for consuming the message from the outbox
func (r *RabbitMQMessageRelay) Consume() ([]outbox.MessageOutbox, error) {
	messages, err := r.outboxRepository.FetchMessages()
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// Publish is responsible for publishing the message to the message broker
func (r *RabbitMQMessageRelay) Publish() error {
	messages, err := r.Consume()
	if err != nil {
		return err
	}
	// publish to RabbitMQ
	return nil
}
