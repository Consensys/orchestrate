package dataagents

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"

	log "github.com/sirupsen/logrus"

	"github.com/go-pg/pg/v9"
	uuid "github.com/satori/go.uuid"
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
		errMessage := "failed to create DB transaction for job"
		log.WithError(err).Error(errMessage)
		return errors.PostgresConnectionError(errMessage).ExtendComponent(jobDAComponent)
	}

	transaction := job.Transaction
	transaction.UUID = uuid.NewV4().String()
	err = postgres.Insert(ctx, tx, transaction)
	if err != nil {
		return errors.FromError(err).ExtendComponent(jobDAComponent)
	}

	job.TransactionID = transaction.ID
	job.UUID = uuid.NewV4().String()
	err = postgres.Insert(ctx, tx, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(jobDAComponent)
	}

	for _, logModel := range job.Logs {
		logModel.UUID = uuid.NewV4().String()
		logModel.JobID = job.ID
		err = postgres.Insert(ctx, tx, logModel)
		if err != nil {
			return errors.FromError(err).ExtendComponent(jobDAComponent)
		}
	}

	return tx.Commit()
}

// FindOneByUUID gets a job by UUID
func (agent *PGJob) FindOneByUUID(ctx context.Context, jobUUID string) (*models.Job, error) {
	job := &models.Job{}
	logger := log.WithField("job_uuid", jobUUID)

	err := agent.db.ModelContext(ctx, job).
		Where("job.uuid = ?", jobUUID).
		Relation("Transaction").
		Relation("Schedule").
		Relation("Logs").
		Select()
	if err != nil && err == pg.ErrNoRows {
		errMessage := "job does not exist"
		logger.WithError(err).Error(errMessage)
		return nil, errors.NotFoundError(errMessage).ExtendComponent(jobDAComponent)
	} else if err != nil {
		errMessage := "could not load job"
		logger.WithError(err).Errorf("could not load job")
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(jobDAComponent)
	}

	return job, nil
}
