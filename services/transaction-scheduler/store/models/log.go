package models

import "time"

type Log struct {
	tableName struct{} `pg:"logs"` // nolint:unused,structcheck // reason

	ID        int
	UUID      string
	JobID     int
	Status    string
	Message   string
	CreatedAt time.Time `pg:"default:now()"`
}
