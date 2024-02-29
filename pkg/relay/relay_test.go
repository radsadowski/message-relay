package relay

import (
	"context"
	"errors"
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

func TestRelayMessage(t *testing.T) {
	ctx := context.Background()

	t.Run("test relay unsent message", func(t *testing.T) {
		m := &Message{ID: 1, Payload: "encoded-json", IsSent: false, CreatedAt: time.Unix(111, 0)}
		repo := &DummyOutboxRepository{m: m}
		pub := &DummyPublisher{}
		relay, err := NewRelay(repo, pub, []Retrier{})
		assert.NoError(t, err)

		err = relay.RelayMessage(ctx)

		assert.NoError(t, err)
		assert.Equal(t, true, m.IsSent)
	})

	t.Run("test do not relay sent message", func(t *testing.T) {
		m := &Message{ID: 1, Payload: "encoded-json", IsSent: true, CreatedAt: time.Unix(111, 0)}
		repo := &DummyOutboxRepository{m: m}
		pub := &DummyPublisher{}
		relay, err := NewRelay(repo, pub, []Retrier{})
		assert.NoError(t, err)

		err = relay.RelayMessage(ctx)

		assert.Error(t, err)
		assert.Equal(t, "message with id 1 has already been sent", err.Error())
		assert.Equal(t, true, m.IsSent)
	})

	t.Run("test relay message with error during fetch", func(t *testing.T) {
		m := &Message{ID: 1, Payload: "encoded-json", IsSent: false, CreatedAt: time.Unix(111, 0)}
		repo := &DummyOutboxRepository{m: m, err: errors.New("error fetching message")}
		pub := &DummyPublisher{}
		relay, err := NewRelay(repo, pub, []Retrier{})
		assert.NoError(t, err)

		err = relay.RelayMessage(ctx)
		assert.Error(t, err)
		assert.Equal(t, "error fetching message", err.Error())
		assert.Equal(t, false, m.IsSent)
	})

	t.Run("test relay message with error during publish", func(t *testing.T) {
		m := &Message{ID: 1, Payload: "encoded-json", IsSent: false, CreatedAt: time.Unix(111, 0)}
		repo := &DummyOutboxRepository{m: m}
		pub := &DummyPublisher{err: errors.New("error publishing message")}
		relay, err := NewRelay(repo, pub, []Retrier{})
		assert.NoError(t, err)

		err = relay.RelayMessage(ctx)

		assert.Error(t, err)
		assert.Equal(t, "error publishing message", err.Error())
		assert.Equal(t, false, m.IsSent)
	})
}

type DummyRetrier struct {
	err error
}

func (d *DummyRetrier) Retry(ctx context.Context, err error) error {
	if d.err != nil {
		return d.err
	}
	return nil
}

func TestRetry(t *testing.T) {
	ctx := context.Background()

	t.Run("test retry with retriers", func(t *testing.T) {
		repo := &DummyOutboxRepository{}
		pub := &DummyPublisher{}
		relay, err := NewRelay(repo, pub, []Retrier{&DummyRetrier{}})
		assert.NoError(t, err)

		err = relay.Retry(ctx, errors.New("based on the error type, we can retry or not"))

		assert.NoError(t, err)
	})

	t.Run("test retry with retriers with error", func(t *testing.T) {
		repo := &DummyOutboxRepository{}
		pub := &DummyPublisher{}
		relay, err := NewRelay(repo, pub, []Retrier{&DummyRetrier{err: errors.New("error retrying")}})
		assert.NoError(t, err)

		err = relay.Retry(ctx, errors.New("based on the error type, we can retry or not"))

		assert.Error(t, err)
	})
}
