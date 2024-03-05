package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultInitializer(t *testing.T) {
	t.Run("unsupported database and broker type", func(t *testing.T) {
		conf := &Config{
			DatabaseType: "unsupported",
			BrokerType:   "unsupported",
		}

		initializer, err := GetDefaultInitialazer(conf)

		assert.Error(t, err)
		assert.Nil(t, initializer)
		assert.Equal(t, errors.New("unsupported database and broker type"), err)
	})
	t.Run("supported database and broker type", func(t *testing.T) {
		conf := &Config{
			DatabaseType: "mysql",
			BrokerType:   "rabbitmq",
		}

		initializer, err := GetDefaultInitialazer(conf)

		assert.NoError(t, err)
		assert.NotNil(t, initializer)
		assert.IsType(t, &MySQLWithRabbitMQRelayInitializer{}, initializer)
	})
}

func TestMySQLWithRabbitMQRelayInitializer(t *testing.T) {
	t.Run("test init relay", func(t *testing.T) {
		initializer := &MySQLWithRabbitMQRelayInitializer{
			DBURL:     "admin:secret123@tcp(localhost:3306)/outbox?charset=utf8mb4&parseTime=True&loc=Local",
			BrokerURL: "amqp://guest:guest@localhost:5672/",
		}

		relay, err := initializer.InitRelay()

		assert.NoError(t, err)
		assert.NotNil(t, relay)
	})

	t.Run("test init relay with db error", func(t *testing.T) {
		initializer := &MySQLWithRabbitMQRelayInitializer{
			DBURL:     "invalid",
			BrokerURL: "amqp://guest:guest@localhost:5672/",
		}

		relay, err := initializer.InitRelay()

		assert.Error(t, err)
		assert.Nil(t, relay)
	})

	t.Run("test init relay with broker error", func(t *testing.T) {
		initializer := &MySQLWithRabbitMQRelayInitializer{
			DBURL:     "admin:secret123@tcp(localhost:3306)/outbox?charset=utf8mb4&parseTime=True&loc=Local",
			BrokerURL: "invalid",
		}
		relay, err := initializer.InitRelay()

		assert.Error(t, err)
		assert.Nil(t, relay)
	})
}
