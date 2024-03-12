// Package relay provides the logic for relaying messages from the outbox to the broker
package relay

import (
	"context"
	"fmt"
)

// MessageAlreadySentError is an error that represents the message already sent error
type MessageAlreadySentError struct {
	MessageID int
}

func (e *MessageAlreadySentError) Error() string {
	return fmt.Sprintf("message with id %d has already been sent", e.MessageID)
}

// Retrier is an interface that represents the retrier
type Retrier interface {
	Retry(ctx context.Context, err error) error
}

// Relay is a struct that represents the relay
type Relay struct {
	Repo      MessageOutboxRepository
	Publisher Publisher
	Retriers  []Retrier
}

// RelayMessage is a method that relays the message
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

// Retry is a method that retries the relay
func (r *Relay) Retry(ctx context.Context, err error) error {
	for _, retrier := range r.Retriers {
		err := retrier.Retry(ctx, err)
		if err != nil {
			return err
		}
	}
	return nil
}
