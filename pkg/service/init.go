package service

import (
	"github.com/radekska/outbox/pkg/relay"
)

type RelayInitializer interface {
	InitRelay() (*relay.Relay, error)
}

type MySQLWithRabbitMQRelayInitializer struct {
	DBURL     string
	BrokerURL string
}

func (m *MySQLWithRabbitMQRelayInitializer) InitRelay() (*relay.Relay, error) {
	repo, err := NewMySQLMessageOutboxRepository(m.DBURL)
	if err != nil {
		return nil, err
	}
	publisher, err := NewRabbitMQPublisher(m.BrokerURL)
	if err != nil {
		return nil, err
	}
	// &MySQLRetrier{DB: repo.DB, Logger: m.Logger}
	retriers := []relay.Retrier{}
	return &relay.Relay{
		Repo:      repo,
		Publisher: publisher,
		Retriers:  retriers,
	}, nil
}
