package service

import (
	"github.com/radekska/outbox/pkg/relay"
	log "github.com/sirupsen/logrus"
)

// MySQLWithRabbitMQRelayInitializer is a struct that represents the initializer for the MySQL with RabbitMQ relay
type MySQLWithRabbitMQRelayInitializer struct {
	DBURL     string
	BrokerURL string
	Logger    *log.Logger
}

// Init is a method that initializes the MySQL with RabbitMQ relay
func (m *MySQLWithRabbitMQRelayInitializer) Init() (*relay.Relay, error) {
	publisher, err := NewRabbitMQPublisher(m.BrokerURL)
	if err != nil {
		return nil, err
	}
	repo, err := NewMySQLMessageOutboxRepository(m.DBURL)
	if err != nil {
		return nil, err
	}
	retriers := []relay.Retrier{&MySQLRetrier{DB: repo.DB, Logger: m.Logger}}
	relay, err := relay.NewRelay(repo, publisher, retriers)
	if err != nil {
		return nil, err
	}
	return relay, nil
}
