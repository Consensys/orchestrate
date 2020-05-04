package dataagents

import (
	"context"

	"github.com/go-pg/pg/v9"
	uuid "github.com/satori/go.uuid"
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
	return insert(ctx, agent.db.ModelContext(ctx, logModel), logDAComponent)
}
