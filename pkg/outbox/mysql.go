package outbox

type MySQLMessageOutboxRepository struct{}

func (m *MySQLMessageOutboxRepository) FetchMessages() ([]MessageOutbox, error) {
	return []MessageOutbox{}, nil
}
