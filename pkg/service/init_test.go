package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
