package models

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

type Log struct {
	tableName struct{} `pg:"logs"` // nolint:unused,structcheck // reason

	ID        int `pg:"alias:id"`
	UUID      string
	JobID     *int `pg:"alias:job_id,notnull"`
	Job       *Job
	Status    entities.JobStatus
	Message   string
	CreatedAt time.Time `pg:"default:now()"`
}
