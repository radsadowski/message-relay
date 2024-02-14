package pkg

import (
	"github.com/radekska/outbox/pkg/outbox"
	"gorm.io/gorm"
)

type SQLMessage struct {
	ID        string `gorm:"primaryKey"`
	Payload   string
	IsSent    bool   `gorm:"index"`
	CreatedAt string `gorm:"index"`
}

type MySQLMessageOutboxRepository struct {
	// DB is a pointer to the database connection
	DB *gorm.DB
}

func (m *MySQLMessageOutboxRepository) FetchMessage() (*outbox.Message, error) {
	var r SQLMessage
	tx := m.DB.Model(SQLMessage{IsSent: false}).First(&r)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &outbox.Message{ID: r.ID, Payload: r.Payload, IsSent: r.IsSent, CreatedAt: r.CreatedAt}, nil
}
