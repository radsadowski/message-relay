package relay

type MessageRelay interface {
	// Consume is responsible for consuming the message from the outbox
	Consume() (any, error)
	// Publish is responsible for publishing the message to the message broker
	Publish() error
}
