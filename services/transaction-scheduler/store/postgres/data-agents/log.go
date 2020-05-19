package dataagents

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	uuid "github.com/satori/go.uuid"
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

const logDAComponent = "data-agents.log"

// PGLog is a log data agent for PostgreSQL
type PGLog struct {
	db pg.DB
}

// NewPGLog creates a new PGLog
func NewPGLog(db pg.DB) *PGLog {
	return &PGLog{db: db}
}

// Insert Inserts a new log in DB
func (agent *PGLog) Insert(ctx context.Context, logModel *models.Log) error {
	if logModel.UUID == "" {
		logModel.UUID = uuid.NewV4().String()
	}

	if logModel.JobID == nil && logModel.Job != nil {
		logModel.JobID = &logModel.Job.ID
	}

	err := pg.Insert(ctx, agent.db, logModel)
	if err != nil {
		return errors.FromError(err).ExtendComponent(logDAComponent)
	}

	return nil
}
