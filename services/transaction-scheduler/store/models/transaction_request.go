package models

import (
	"time"
)

type TransactionRequest struct {
	tableName struct{} `pg:"requests"` // nolint:unused,structcheck // reason

	ID             int
	IdempotencyKey string
	Chain          string
	Method         string
	Params         string
	Labels         *map[string]string
	CreatedAt      time.Time `pg:"default:now()"`
}
