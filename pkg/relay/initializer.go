package relay

type Initializer interface {
	Init() (*Relay, error)
}
