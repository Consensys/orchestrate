package models

import (
	"time"
)

type TransactionRequest struct {
	tableName struct{} `pg:"requests"` // nolint:unused,structcheck // reason

	ID             int
	IdempotencyKey string
	ScheduleID     int
	RequestHash    string
	Params         string
	CreatedAt      time.Time `pg:"default:now()"`
}
