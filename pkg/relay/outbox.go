package relay

import (
	"context"
	"time"
)

// Message is a structure that represents a message that is going to be relayed
type Message struct {
	ID        int       `json:"id"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
	IsSent    bool      `json:"is_sent"`
}

// MessageOutboxRepository is an interface that represents the outbox repository
type MessageOutboxRepository interface {
	// FetchMessage is responsible for fetching the unsent messages from the outbox
	FetchMessage(ctx context.Context) (*Message, error)
	// SetMessageAsSent is responsible for setting the message as sent
	SetMessageAsSent(ctx context.Context, msg *Message) error
}
