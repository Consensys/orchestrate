package store

import (
	"context"

	healthz "github.com/heptiolabs/healthcheck"
)

//go:generate mockgen -source=store.go -destination=mocks/mock.go -package=mocks

type Vault interface {
	Agents
	HealthCheck() healthz.Check
}

type Agents interface {
	Ethereum() EthereumAgent
}

// Interfaces data agents
type EthereumAgent interface {
	Insert(ctx context.Context, address, privKey, namespace string) error
	FindOne(ctx context.Context, address, namespace string) (string, error)
}
