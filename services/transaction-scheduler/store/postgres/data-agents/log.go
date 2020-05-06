package dataagents

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	"github.com/go-pg/pg/v9"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

const logDAComponent = "data-agents.log"

// PGLog is a log data agent for PostgreSQL
type PGLog struct {
	db *pg.DB
}

// NewPGLog creates a new PGLog
func NewPGLog(db *pg.DB) *PGLog {
	return &PGLog{db: db}
}

// Insert Inserts a new log in DB
func (agent *PGLog) Insert(ctx context.Context, logModel *models.Log) error {
	logModel.UUID = uuid.NewV4().String()
	err := postgres.Insert(ctx, agent.db, logModel)
	if err != nil {
		return errors.FromError(err).ExtendComponent(logDAComponent)
	}

	return nil
}
