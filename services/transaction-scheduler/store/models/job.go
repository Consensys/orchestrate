package models

import "time"

type Job struct {
	tableName struct{} `pg:"jobs"` // nolint:unused,structcheck // reason

	ID            int
	UUID          string
	ScheduleID    int
	Schedule      *Schedule
	Type          string
	TransactionID int
	Transaction   *Transaction
	Logs          []*Log
	Labels        map[string]string
	CreatedAt     time.Time `pg:"default:now()"`
}

// GetStatus Computes the status of a Job by checking its logs
func (job *Job) GetStatus() string {
	// TODO: Order logs by createdAt when getting them from DB
	return job.Logs[len(job.Logs)-1].Status
}
