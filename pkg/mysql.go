package pkg

import (
	"context"

	"github.com/radekska/outbox/pkg/outbox"
	"gorm.io/gorm"
)

type SQLMessage struct {
	ID        int `gorm:"primaryKey;autoIncrement:true"`
	Payload   string
	IsSent    bool   `gorm:"index"`
	CreatedAt string `gorm:"index"`
}

type MySQLMessageOutboxRepository struct {
	// DB is a pointer to the database connection
	DB *gorm.DB
}

func (m *MySQLMessageOutboxRepository) FetchMessage(ctx context.Context) (*outbox.Message, error) {
	var r SQLMessage
	tx := m.DB.WithContext(ctx).Model(&SQLMessage{IsSent: false}).First(&r)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &outbox.Message{ID: r.ID, Payload: r.Payload, IsSent: r.IsSent, CreatedAt: r.CreatedAt}, nil
}

func (m *MySQLMessageOutboxRepository) SetMessageAsSent(ctx context.Context, msg *outbox.Message) error {
	tx := m.DB.WithContext(ctx).Model(&SQLMessage{ID: msg.ID}).Update("is_sent", true)
	return tx.Error
}
