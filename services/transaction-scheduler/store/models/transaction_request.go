package models

import (
	"time"
)

type TransactionRequest struct {
	tableName struct{} `pg:"transaction_requests"` // nolint:unused,structcheck // reason

	ID             int
	IdempotencyKey string
	Schedules      []*Schedule
	RequestHash    string
	Params         string
	CreatedAt      time.Time `pg:"default:now()"`
}
