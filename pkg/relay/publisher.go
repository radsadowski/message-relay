package relay

import (
	"context"
)

// Publisher is an interface that represents the publisher
type Publisher interface {
	Publish(ctx context.Context, msg *Message) error
}
