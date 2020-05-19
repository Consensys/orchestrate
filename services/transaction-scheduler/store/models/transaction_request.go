package models

import (
	"time"
)

type TransactionRequest struct {
	tableName struct{} `pg:"requests"` // nolint:unused,structcheck // reason

	ID             int
	IdempotencyKey string
	ScheduleID     *int `pg:"alias:schedule_id,notnull"`
	Schedule       *Schedule
	RequestHash    string
	Params         string
	CreatedAt      time.Time `pg:"default:now()"`
}
