package relay

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type DummyOutboxRepository struct {
	m   *Message
	err error
}

func (d *DummyOutboxRepository) FetchMessage(ctx context.Context) (*Message, error) {
	if d.err != nil {
		return nil, d.err
	}
	return d.m, nil
}

func (d *DummyOutboxRepository) SetMessageAsSent(ctx context.Context, msg *Message) error {
	if d.err != nil {
		return d.err
	}
	msg.IsSent = true
	return nil
}

type DummyPublisher struct {
	err error
}

func (d *DummyPublisher) Publish(ctx context.Context, msg *Message) error {
	if d.err != nil {
		return d.err
	}
	return nil
}

func TestRelay(t *testing.T) {
	ctx := context.Background()

	t.Run("test relay unsent message", func(t *testing.T) {
		m := &Message{ID: 1, Payload: "encoded-json", IsSent: false, CreatedAt: time.Unix(111, 0)}
		repo := &DummyOutboxRepository{m: m}
		pub := &DummyPublisher{}
		relay, err := NewRelay(repo, pub)
		assert.NoError(t, err)

		err = relay.RelayMessage(ctx)

		assert.NoError(t, err)
		assert.Equal(t, true, m.IsSent)
	})

	t.Run("test do not relay sent message", func(t *testing.T) {
		m := &Message{ID: 1, Payload: "encoded-json", IsSent: true, CreatedAt: time.Unix(111, 0)}
		repo := &DummyOutboxRepository{m: m}
		pub := &DummyPublisher{}
		relay, err := NewRelay(repo, pub)
		assert.NoError(t, err)

		err = relay.RelayMessage(ctx)

		assert.Error(t, err)
		assert.Equal(t, "message with id 1 has already been sent", err.Error())
		assert.Equal(t, true, m.IsSent)
	})
}
