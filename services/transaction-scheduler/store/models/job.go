package models

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

type Job struct {
	tableName struct{} `pg:"jobs"` // nolint:unused,structcheck // reason

	ID            int `pg:"alias:id"`
	UUID          string
	ChainUUID     string
	NextJobUUID   string `pg:"alias:next_job_uuid"`
	ScheduleID    *int   `pg:"alias:schedule_id,notnull"`
	Schedule      *Schedule
	Type          string
	TransactionID *int `pg:"alias:transaction_id,notnull"`
	Transaction   *Transaction
	Logs          []*Log
	Labels        map[string]string
	InternalData  *entities.InternalData
	CreatedAt     time.Time `pg:"default:now()"`
	UpdatedAt     time.Time `pg:"default:now()"`
}
