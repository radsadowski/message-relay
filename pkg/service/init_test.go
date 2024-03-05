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
