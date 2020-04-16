package store

type DataAgents struct {
	TransactionRequest TransactionRequestAgent
	TransactionJob     TransactionJobAgent
}

// Interfaces data agents
type TransactionRequestAgent interface {
	// TODO
}

type TransactionJobAgent interface {
	// TODO
}
