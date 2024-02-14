package outbox

import (
	"context"
)

type Message struct {
	ID        string `json:"id"`
	Payload   string `json:"payload"`
	IsSent    bool   `json:"is_sent"`
	CreatedAt string `json:"created_at"`
}

type MessageOutboxRepository interface {
	// FetchMessage is responsible for fetching the unsent messages from the outbox
	FetchMessage(ctx context.Context) (*Message, error)
	// SetMessageAsSent is responsible for setting the message as sent
	SetMessageAsSent(ctx context.Context, msg *Message) error
}
