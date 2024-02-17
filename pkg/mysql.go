package pkg

import (
	"context"
	"net"
	"time"

	"github.com/radekska/outbox/pkg/relay"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// SQLMessage is a struct that represents a message in the database
type SQLMessage struct {
	ID        int `gorm:"primaryKey;autoIncrement:true"`
	Payload   string
	IsSent    bool      `gorm:"index"`
	CreatedAt time.Time `gorm:"index"`
}

// MessageOutboxRepository is an interface implementation that defines the methods for interacting with the message outbox
type MySQLMessageOutboxRepository struct {
	// DB is a pointer to the database connection
	DB *gorm.DB
}

// NewMySQLMessageOutboxRepository is a constructor for the MySQLMessageOutboxRepository
func NewMySQLMessageOutboxRepository(dbURL string) (*MySQLMessageOutboxRepository, error) {
	db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}
	return &MySQLMessageOutboxRepository{DB: db}, nil
}

// FetchMessage is a method that fetches a message from the outbox
func (m *MySQLMessageOutboxRepository) FetchMessage(ctx context.Context) (*relay.Message, error) {
	var r SQLMessage

	tx := m.DB.WithContext(ctx).Where(map[string]interface{}{"is_sent": false}).Limit(1).Find(&r)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &relay.Message{ID: r.ID, Payload: r.Payload, CreatedAt: r.CreatedAt}, nil
}

// SetMessageAsSent is a method that sets a message as sent in the outbox
func (m *MySQLMessageOutboxRepository) SetMessageAsSent(ctx context.Context, msg *relay.Message) error {
	tx := m.DB.WithContext(ctx).Model(&SQLMessage{ID: msg.ID}).Update("is_sent", true)
	return tx.Error
}

type MySQLRetrier struct {
	DB     *gorm.DB
	Logger *log.Logger
}

func (d *MySQLRetrier) Retry(ctx context.Context, err error) error {
	e, ok := err.(*net.OpError)
	if !ok {
		return nil
	}
	d.Logger.Errorf("database network error: %v", e)
	d.Logger.Warn("reconnecting to the database...")
	db, err := d.DB.DB()
	if err != nil {
		return err
	}
	return db.PingContext(ctx)
}
