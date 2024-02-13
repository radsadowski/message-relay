package pkg

import (
	"context"
	"testing"

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
	m := &Message{ID: "test-msg", Payload: "encoded-json", IsSent: false, CreatedAt: "2022-01-01 10:00:00"}
	repo := &DummyOutboxRepository{m: m}
	pub := &DummyPublisher{}
	ctx := context.Background()
	relay := NewRelay(repo, pub)

	err := relay.RelayMessage(ctx)

	assert.NoError(t, err)
	assert.Equal(t, true, m.IsSent)
}
