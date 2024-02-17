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

var brokerURLFlag = &cli.StringFlag{
	Name:     "broker-url",
	Aliases:  []string{"b"},
	Usage:    "URL to the broker",
	EnvVars:  []string{"BROKER_URL"},
	Required: true,
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
