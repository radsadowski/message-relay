package service

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/radekska/outbox/pkg/relay"
	log "github.com/sirupsen/logrus"
)

// Service is a struct that represents the service
type Service struct {
	messageRelay *relay.Relay
	healthCh     chan struct{}
	signalCh     chan os.Signal
	logger       *log.Logger
	config       *Config
}

// NewService is a constructor for the Service
func NewService(messageRelay *relay.Relay, config *Config, logger *log.Logger) *Service {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	healthCh := make(chan struct{})

	go func() {
		startHealthCheckServer(healthCh, logger)
	}()
	go func() {
		startMetricsServer(logger)
	}()
	return &Service{
		messageRelay: messageRelay,
		signalCh:     signalCh,
		logger:       logger,
		config:       config,
		healthCh:     healthCh,
	}
}

// Run is a method that runs the service
func (s *Service) Run(ctx context.Context, cancel context.CancelFunc) error {
	s.logger.Info("starting the relay service...")
	ticker := time.NewTicker(s.config.RelayInterval)
	for {
		select {
		case signal := <-s.signalCh:
			s.shutdown(cancel, signal)
		case <-s.healthCh:
			s.logger.Info("health check passed")
		case <-ticker.C:
			s.relay(ctx)
		}
	}
}

func (s *Service) relay(ctx context.Context) {
	s.logger.Info("relaying message...")
	err := s.messageRelay.RelayMessage(ctx)
	if err != nil {
		s.messageRelay.Retry(ctx, err)
	} else {
		s.logger.Info("message relayed successfully")
		relayedMessagesCounter.Inc()
	}
}

func (s *Service) shutdown(cancel context.CancelFunc, signal os.Signal) {
	s.logger.Warnf("received signal: %v", signal)
	s.logger.Warn("quitting gracefully...")
	cancel()
	os.Exit(0)
}

// Start is a function that starts the service
func Start(ctx context.Context, conf *Config) error {
	logger := log.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	initializer := &MySQLWithRabbitMQRelayInitializer{DBURL: conf.DBURL, BrokerURL: conf.BrokerURL, Logger: logger}
	messageRelay, err := initializer.Init()
	if err != nil {
		logger.Errorf("failed to initialize the relay: %v", err)
		return err
	}
	service := NewService(messageRelay, conf, logger)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return service.Run(ctx, cancel)
}
