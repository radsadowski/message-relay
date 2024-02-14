package pkg

import (
	"context"

	"github.com/radekska/outbox/pkg/outbox"
)

type Publisher interface {
	Publish(ctx context.Context, msg *outbox.Message) error
}
