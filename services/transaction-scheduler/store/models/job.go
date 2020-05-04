package models

import "time"

type Job struct {
	tableName struct{} `pg:"jobs"` // nolint:unused,structcheck // reason

	ID            int
	UUID          string
	ScheduleID    int
	Schedule      *Schedule `pg:"-"`
	Type          string
	TransactionID int
	Transaction   *Transaction `pg:"-"`
	Logs          []*Log       `pg:"-"`
	Labels        *map[string]string
	CreatedAt     time.Time `pg:"default:now()"`
}
