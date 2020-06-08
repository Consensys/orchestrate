package models

import (
	"time"
)

type Job struct {
	tableName struct{} `pg:"jobs"` // nolint:unused,structcheck // reason

	ID            int `pg:"alias:id"`
	UUID          string
	ChainUUID     string
	ScheduleID    *int `pg:"alias:schedule_id,notnull"`
	Schedule      *Schedule
	Type          string // @TODO Replace by enum
	TransactionID *int   `pg:"alias:transaction_id,notnull"`
	Transaction   *Transaction
	Logs          []*Log
	Labels        map[string]string
	CreatedAt     time.Time `pg:"default:now()"`
}
