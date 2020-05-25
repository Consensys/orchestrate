package models

import (
	"time"
)

type Job struct {
	tableName struct{} `pg:"jobs"` // nolint:unused,structcheck // reason

	ID            int `pg:"alias:id"`
	UUID          string
	ScheduleID    *int `pg:"alias:schedule_id,notnull"`
	Schedule      *Schedule
	Type          string // @TODO Replace by enum
	TransactionID *int   `pg:"alias:transaction_id,notnull"`
	Transaction   *Transaction
	Logs          []*Log
	Labels        map[string]string
	CreatedAt     time.Time `pg:"default:now()"`
}

// GetStatus Computes the status of a Job by checking its logs
func (job *Job) GetStatus() string {
	var status string
	var logCreatedAt *time.Time
	for idx := range job.Logs {
		if logCreatedAt == nil || job.Logs[idx].CreatedAt.After(*logCreatedAt) {
			status = job.Logs[idx].Status
			logCreatedAt = &job.Logs[idx].CreatedAt
		}
	}

	return status
}
