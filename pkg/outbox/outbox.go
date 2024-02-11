package outbox

type MessageOutbox struct {
	ID        string
	Payload   string
	IsSent    bool
	CreatedAt string
}

type MessageOutboxRepository interface {
	// FetchMessages is responsible for fetching the unsent messages from the outbox
	FetchMessages() ([]MessageOutbox, error)
}
