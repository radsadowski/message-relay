package pkg

import (
	"context"
	"fmt"

	"github.com/radekska/outbox/pkg/outbox"
)

type ErrMessageAlreadySent struct {
	MessageID string
}

func (e *ErrMessageAlreadySent) Error() string {
	return fmt.Sprintf("message with id %s has already been sent", e.MessageID)
}

type Relay struct {
	outboxRepository outbox.MessageOutboxRepository
	publisher        Publisher
}

func NewRelay(outboxRepo outbox.MessageOutboxRepository, publisher Publisher) *Relay {
	return &Relay{
		outboxRepository: outboxRepo,
		publisher:        publisher,
	}
}

func (r *Relay) RelayMessage(ctx context.Context) error {
	m, err := r.outboxRepository.FetchMessage(ctx)
	if err != nil {
		return err
	}
	if m.IsSent {
		return &ErrMessageAlreadySent{MessageID: m.ID}
	}
	err = r.publisher.Publish(ctx, m)
	if err != nil {
		return err
	}
	return r.outboxRepository.SetMessageAsSent(ctx, m)
}
