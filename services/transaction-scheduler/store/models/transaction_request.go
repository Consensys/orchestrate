package models

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
)

type TransactionRequest struct {
	tableName struct{} `pg:"transaction_requests"` // nolint:unused,structcheck // reason

	ID             int
	UUID           string
	IdempotencyKey string
	Schedules      []*Schedule
	RequestHash    string
	Params         *types.ETHTransactionParams // This will be automatically transformed in JSON by go-pg (and vice-versa)
	CreatedAt      time.Time                   `pg:"default:now()"`
}
