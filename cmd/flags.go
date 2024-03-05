package cmd

import (
	"time"

	"github.com/urfave/cli/v2"
)

var databaseURLFlag = &cli.StringFlag{
	Name:     "database-url",
	Aliases:  []string{"d"},
	Usage:    "URL to the database",
	EnvVars:  []string{"DATABASE_URL"},
	Required: true,
}

var databaseTypeFlag = &cli.StringFlag{
	Name:     "database-type",
	Usage:    "Type of the database",
	EnvVars:  []string{"DATABASE_TYPE"},
	Required: false,
	Value:    "mysql",
}

var brokerURLFlag = &cli.StringFlag{
	Name:     "broker-url",
	Aliases:  []string{"b"},
	Usage:    "URL to the broker",
	EnvVars:  []string{"BROKER_URL"},
	Required: true,
}

var brokerTypeFlag = &cli.StringFlag{
	Name:     "broker-type",
	Usage:    "Type of the broker",
	EnvVars:  []string{"BROKER_TYPE"},
	Required: false,
	Value:    "rabbitmq",
}

var relayIntervalFlag = &cli.DurationFlag{
	Name:     "relay-interval",
	Usage:    "Interval between relaying messages",
	EnvVars:  []string{"RELAY_INTERVAL"},
	Value:    100 * time.Millisecond,
	Required: false,
}

var batchSizeFlag = &cli.IntFlag{
	Name:     "batch-size",
	Usage:    "Number of messages to relay in a batch",
	EnvVars:  []string{"BATCH_SIZE"},
	Value:    100,
	Required: false,
}

var retriesCountFlag = &cli.IntFlag{
	Name:     "retries",
	Usage:    "Number of retries for each message",
	EnvVars:  []string{"RETRIES"},
	Value:    100,
	Required: false,
}
