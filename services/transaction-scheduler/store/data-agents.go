package store

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

//go:generate mockgen -source=data-agents.go -destination=mocks/data-agents.go -package=mocks

type DataAgents struct {
	TransactionRequest TransactionRequestAgent
}

// Interfaces data agents

type TransactionRequestAgent interface {
	Insert(ctx context.Context, txRequest *models.TransactionRequest) error
}
