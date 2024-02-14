package pkg

import (
	"context"
	"testing"

	"github.com/radekska/outbox/pkg/outbox"
	"github.com/stretchr/testify/assert"
)

type DummyOutboxRepository struct {
	m   *outbox.Message
	err error
}

func (d *DummyOutboxRepository) FetchMessage(ctx context.Context) (*outbox.Message, error) {
	if d.err != nil {
		return nil, d.err
	}
	return d.m, nil
}

func (d *DummyOutboxRepository) SetMessageAsSent(ctx context.Context, msg *outbox.Message) error {
	if d.err != nil {
		return d.err
	}
	msg.IsSent = true
	return nil
}

type DummyPublisher struct {
	err error
}

func (d *DummyPublisher) Publish(ctx context.Context, msg *outbox.Message) error {
	if d.err != nil {
		return d.err
	}
	return nil
}

func TestRelay(t *testing.T) {
	ctx := context.Background()

	t.Run("test relay unsent message", func(t *testing.T) {
		m := &outbox.Message{ID: "test-msg", Payload: "encoded-json", IsSent: false, CreatedAt: "2022-01-01 10:00:00"}
		repo := &DummyOutboxRepository{m: m}
		pub := &DummyPublisher{}

		relay := NewRelay(repo, pub)

		err := relay.RelayMessage(ctx)

		assert.NoError(t, err)
		assert.Equal(t, true, m.IsSent)
	})

	t.Run("test do not relay sent message", func(t *testing.T) {
		m := &outbox.Message{ID: "test-msg", Payload: "encoded-json", IsSent: true, CreatedAt: "2022-01-01 10:00:00"}
		repo := &DummyOutboxRepository{m: m}
		pub := &DummyPublisher{}

		relay := NewRelay(repo, pub)

		err := relay.RelayMessage(ctx)

		assert.Error(t, err)
		assert.Equal(t, "message with id test-msg has already been sent", err.Error())
		assert.Equal(t, true, m.IsSent)
	})
}
