package service

import "time"

// Config is a struct that represents the configuration for the service
type Config struct {
	DBURL         string
	BrokerURL     string
	BatchSize     int
	RelayInterval time.Duration
	RetriesCount  int
}
