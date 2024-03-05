package service

import (
	"errors"

	"github.com/radekska/outbox/pkg/relay"
)

type RelayInitializer interface {
	initMessageOutboxRepository() (relay.MessageOutboxRepository, error)
	initPublisher() (relay.Publisher, error)
	InitRelay() (*relay.Relay, error)
}

type MySQLWithRabbitMQRelayInitializer struct {
	DBURL     string
	BrokerURL string
}

func (m *MySQLWithRabbitMQRelayInitializer) initMessageOutboxRepository() (relay.MessageOutboxRepository, error) {
	return NewMySQLMessageOutboxRepository(m.DBURL)
}

func (m *MySQLWithRabbitMQRelayInitializer) initPublisher() (relay.Publisher, error) {
	return NewRabbitMQPublisher(m.BrokerURL)
}

func (m *MySQLWithRabbitMQRelayInitializer) InitRelay() (*relay.Relay, error) {
	repo, err := m.initMessageOutboxRepository()
	if err != nil {
		return nil, err
	}
	publisher, err := m.initPublisher()
	if err != nil {
		return nil, err
	}
	// &MySQLRetrier{DB: repo.DB, Logger: m.Logger}
	retriers := []relay.Retrier{}
	relay, err := relay.NewRelay(repo, publisher, retriers)
	if err != nil {
		return nil, err
	}
	return relay, nil
}

type GetInitializerFunc func(*Config) (RelayInitializer, error)

func GetDefaultInitialazer(conf *Config) (RelayInitializer, error) {
	switch {
	case conf.DatabaseType == "mysql" && conf.BrokerType == "rabbitmq":
		return &MySQLWithRabbitMQRelayInitializer{DBURL: conf.DBURL, BrokerURL: conf.BrokerURL}, nil
	default:
		return nil, errors.New("unsupported database and broker type")
	}
}
