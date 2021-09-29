package models

import (
	"time"

	"github.com/consensys/orchestrate/pkg/types/entities"
)

type TransactionRequest struct {
	tableName struct{} `pg:"transaction_requests"` // nolint:unused,structcheck // reason

	ID             int
	IdempotencyKey string
	ChainName      string
	ScheduleID     *int
	Schedule       *Schedule
	RequestHash    string
	Params         *entities.ETHTransactionParams // This will be automatically transformed in JSON by go-pg (and vice-versa)
	CreatedAt      time.Time                      `pg:"default:now()"`
}
