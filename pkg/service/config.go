package service

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Config is a struct that represents the configuration for the service
type Config struct {
	BatchSize     int
	RelayInterval time.Duration
	RetriesCount  int
	LogLevel      string
}

func (c *Config) MustParseLogLevel() log.Level {
	level, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		log.Fatalf("failed to parse log level: %s", c.LogLevel)
	}
	return level
}
