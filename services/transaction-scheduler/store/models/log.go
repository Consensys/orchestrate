package models

import (
	"time"
)

type Log struct {
	tableName struct{} `pg:"logs"` // nolint:unused,structcheck // reason

	ID        int `pg:"alias:id"`
	UUID      string
	JobID     *int `pg:"alias:job_id,notnull"`
	Job       *Job
	Status    string // @TODO Replace "status" by enum
	Message   string
	CreatedAt time.Time `pg:"default:now()"`
}
