package service

import (
	"github.com/radekska/outbox/pkg/relay"
)

// RelayInitializer is an interface that represents the relay initializer
type RelayInitializer interface {
	InitRelay() (*relay.Relay, error)
}

// MySQLWithRabbitMQRelayInitializer is a struct that represents the MySQL with RabbitMQ relay initializer
type MySQLWithRabbitMQRelayInitializer struct {
	DBURL     string
	BrokerURL string
}

// InitRelay is a method that initializes the relay
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
