package relay

import (
	"context"
	"fmt"
)

type MessageAlreadySentError struct {
	MessageID int
}

func (e *MessageAlreadySentError) Error() string {
	return fmt.Sprintf("message with id %d has already been sent", e.MessageID)
}

type Retrier interface {
	Retry(ctx context.Context, err error) error
}

type Relay struct {
	repo      MessageOutboxRepository
	publisher Publisher
	retriers  []Retrier
}

func NewRelay(repo MessageOutboxRepository, publisher Publisher, retriers []Retrier) (*Relay, error) {
	return &Relay{
		repo:      repo,
		publisher: publisher,
		retriers:  retriers,
	}, nil
}

func (r *Relay) RelayMessage(ctx context.Context) error {
	m, err := r.repo.FetchMessage(ctx)
	if err != nil {
		return err
	}
	if m.IsSent {
		return &MessageAlreadySentError{MessageID: m.ID}
	}
	err = r.publisher.Publish(ctx, m)
	if err != nil {
		return err
	}
	return r.repo.SetMessageAsSent(ctx, m)
}

func (r *Relay) Retry(ctx context.Context, err error) error {
	for _, retrier := range r.retriers {
		err := retrier.Retry(ctx, err)
		if err != nil {
			return err
		}
	}
	return nil
}
