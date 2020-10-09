package store

//go:generate mockgen -source=store.go -destination=mocks/mock.go -package=mocks

type Agents interface {
	Ethereum() EthereumAgent
}

type Vault interface {
	Agents
}

// Interfaces data agents
type EthereumAgent interface{}
