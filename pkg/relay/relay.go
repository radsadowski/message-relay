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
	Repo      MessageOutboxRepository
	Publisher Publisher
	Retriers  []Retrier
}

func (r *Relay) RelayMessage(ctx context.Context) error {
	m, err := r.Repo.FetchMessage(ctx)
	if err != nil {
		return err
	}
	if m.IsSent {
		return &MessageAlreadySentError{MessageID: m.ID}
	}
	err = r.Publisher.Publish(ctx, m)
	if err != nil {
		return err
	}
	return r.Repo.SetMessageAsSent(ctx, m)
}

func (r *Relay) Retry(ctx context.Context, err error) error {
	for _, retrier := range r.Retriers {
		err := retrier.Retry(ctx, err)
		if err != nil {
			return err
		}
	}
	return nil
}
