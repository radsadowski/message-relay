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

var retriesCountFlag = &cli.IntFlag{
	Name:     "retries",
	Usage:    "Number of retries for each message",
	EnvVars:  []string{"RETRIES"},
	Value:    100,
	Required: false,
}

var loggingLevelFlag = &cli.StringFlag{
	Name:    "log-level",
	Usage:   "Logging level. Can be one of [panic, fatal, error, warn, warning, info, debug, trace]",
	EnvVars: []string{"LOG_LEVEL"},
	Value:   "info",
	Action: func(ctx *cli.Context, s string) error {
		levels := []string{"panic", "fatal", "error", "warn", "warning", "info", "debug", "trace"}
		for _, l := range levels {
			if s == l {
				return nil
			}
		}
		return cli.Exit("invalid log level", 1)
	},
}
