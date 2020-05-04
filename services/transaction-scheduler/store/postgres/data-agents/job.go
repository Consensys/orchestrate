package dataagents

import (
	"context"

	"github.com/go-pg/pg/v9"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

const jobDAComponent = "data-agents.job"

// PGJob is a job data agent for PostgreSQL
type PGJob struct {
	db *pg.DB
}

// NewPGJob creates a new pgJob
func NewPGJob(db *pg.DB) *PGJob {
	return &PGJob{db: db}
}

// Insert Inserts a new job in DB
func (agent *PGJob) Insert(ctx context.Context, job *models.Job) error {
	tx, err := agent.db.Begin()
	if err != nil {
		return errors.PostgresConnectionError("Failed to create DB transaction").ExtendComponent(jobDAComponent)
	}
	pgctx := postgres.WithTx(ctx, tx)

	transaction := job.Transaction
	transaction.UUID = uuid.NewV4().String()
	err = insert(pgctx, tx.ModelContext(pgctx, transaction), jobDAComponent)
	if err != nil {
		return errors.FromError(err).ExtendComponent(jobDAComponent)
	}

	job.TransactionID = transaction.ID
	job.UUID = uuid.NewV4().String()
	err = insert(pgctx, tx.ModelContext(pgctx, job), jobDAComponent)
	if err != nil {
		return errors.FromError(err).ExtendComponent(jobDAComponent)
	}

	for _, log := range job.Logs {
		log.UUID = uuid.NewV4().String()
		log.JobID = job.ID
		err = insert(pgctx, tx.ModelContext(pgctx, log), jobDAComponent)
		if err != nil {
			return errors.FromError(err).ExtendComponent(jobDAComponent)
		}
	}

	return tx.Commit()
}
