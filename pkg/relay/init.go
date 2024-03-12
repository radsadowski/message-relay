package relay

// Initializer is an interface that represents the relay initializer
type Initializer interface {
	Init() (*Relay, error)
}
