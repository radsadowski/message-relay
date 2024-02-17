package pkg

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/radekska/outbox/pkg/relay"
	log "github.com/sirupsen/logrus"
)

type MySQLWithRabbitMQRelayInitializer struct {
	DBURL     string
	BrokerURL string
	Logger    *log.Logger
}

func (m *MySQLWithRabbitMQRelayInitializer) Init() (*relay.Relay, error) {
	publisher, err := NewRabbitMQPublisher(m.BrokerURL)
	if err != nil {
		return nil, err
	}
	repo, err := NewMySQLMessageOutboxRepository(m.DBURL)
	if err != nil {
		return nil, err
	}
	retriers := []relay.Retrier{&MySQLRetrier{DB: repo.DB, Logger: m.Logger}}
	relay, err := relay.NewRelay(repo, publisher, retriers)
	if err != nil {
		return nil, err
	}
	return relay, nil
}

type Config struct {
	DBURL         string
	BrokerURL     string
	BatchSize     int
	RelayInterval time.Duration
}

type Service struct {
	messageRelay *relay.Relay
	signalCh     chan os.Signal
	logger       *log.Logger
	config       *Config
}

func NewService(messageRelay *relay.Relay, config *Config, logger *log.Logger) *Service {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	return &Service{messageRelay: messageRelay, signalCh: signalCh, logger: logger, config: config}
}

func (s *Service) Run(ctx context.Context, cancel context.CancelFunc) error {
	s.logger.Info("starting the relay service...")
	ticker := time.NewTicker(s.config.RelayInterval)
	for {
		select {
		case signal := <-s.signalCh:
			s.shutdown(cancel, signal)
		case <-ticker.C:
			s.relay(ctx)
		}
	}
}

func (s *Service) relay(ctx context.Context) {
	s.logger.Info("relaying message")
	err := s.messageRelay.RelayMessage(ctx)
	if err != nil {
		s.messageRelay.Retry(ctx, err)
	}
}

func (s *Service) shutdown(cancel context.CancelFunc, signal os.Signal) {
	s.logger.Warnf("received signal: %v", signal)
	s.logger.Warn("quitting gracefully...")
	cancel()
	os.Exit(0)
}

func Run(ctx context.Context, conf *Config) error {
	logger := log.New()
	initializer := &MySQLWithRabbitMQRelayInitializer{DBURL: conf.DBURL, BrokerURL: conf.BrokerURL, Logger: logger}
	messageRelay, err := initializer.Init()
	if err != nil {
		return err
	}
	service := NewService(messageRelay, conf, logger)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return service.Run(ctx, cancel)
}
