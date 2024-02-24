/*
Copyright © 2024 Radoslaw Skalbania <ra.skalbania@gmail.com>
*/
package cmd

import (
	"context"

	"github.com/radekska/outbox/pkg/service"
	"github.com/urfave/cli/v2"
)

var App = &cli.App{
	Name:  "message-relay",
	Usage: "Relay messages from the outbox to the broker",
	Flags: []cli.Flag{databaseURLFlag, brokerURLFlag, relayIntervalFlag, batchSizeFlag, retriesCountFlag},
	Action: func(cCtx *cli.Context) error {
		// TODO: argument validation - we should validate DB and broker URLs based on their format
		return service.Start(
			context.Background(),
			&service.Config{
				DBURL:         cCtx.String("database-url"),
				BrokerURL:     cCtx.String("broker-url"),
				BatchSize:     cCtx.Int("batch-size"),
				RelayInterval: cCtx.Duration("relay-interval"),
				RetriesCount:  cCtx.Int("retries"),
			},
		)
	},
}
