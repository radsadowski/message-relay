/*
Package cmd provides the command line interface for the message-relay application
Copyright © 2024 Radoslaw Skalbania <ra.skalbania@gmail.com>
*/
package cmd

import (
	"context"

	"github.com/radekska/outbox/pkg/service"
	"github.com/urfave/cli/v2"
)

// App is the main entry point for the application
var App = &cli.App{
	Name:  "message-relay",
	Usage: "Relay messages from the outbox to the broker",
	Flags: []cli.Flag{
		databaseURLFlag,
		brokerURLFlag,
		relayIntervalFlag,
		batchSizeFlag,
		retriesCountFlag,
		loggingLevelFlag,
	},
	Action: func(cCtx *cli.Context) error {
		return service.Start(
			context.Background(),
			&service.Config{
				BatchSize:     cCtx.Int("batch-size"),
				RelayInterval: cCtx.Duration("relay-interval"),
				RetriesCount:  cCtx.Int("retries"),
			},
			&service.MySQLWithRabbitMQRelayInitializer{
				DBURL:     cCtx.String("database-url"),
				BrokerURL: cCtx.String("broker-url"),
			},
		)
	},
}
