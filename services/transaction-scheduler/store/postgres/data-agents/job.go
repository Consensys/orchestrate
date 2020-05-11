package dataagents

import (
	"context"

	"github.com/go-pg/pg/v9/orm"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"

	log "github.com/sirupsen/logrus"

	"github.com/go-pg/pg/v9"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

const jobDAComponent = "data-agents.job"

// PGJob is a job data agent for PostgreSQL
type PGJob struct {
	db orm.DB
}

// NewPGJob creates a new pgJob
func NewPGJob(db orm.DB) *PGJob {
	return &PGJob{db: db}
}

// Insert Inserts a new job in DB
func (agent *PGJob) Insert(ctx context.Context, job *models.Job) error {
	transaction := job.Transaction
	job.Transaction.UUID = uuid.NewV4().String()
	err := postgres.Insert(ctx, agent.db, transaction)
	if err != nil {
		return errors.FromError(err).ExtendComponent(jobDAComponent)
	}

	job.TransactionID = transaction.ID
	job.UUID = uuid.NewV4().String()
	err = postgres.Insert(ctx, agent.db, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(jobDAComponent)
	}

	return nil
}

// FindOneByUUID gets a job by UUID
func (agent *PGJob) FindOneByUUID(ctx context.Context, jobUUID, tenantID string) (*models.Job, error) {
	job := &models.Job{}
	logger := log.WithField("job_uuid", jobUUID)

	query := agent.db.ModelContext(ctx, job).Where("job.uuid = ?", jobUUID)
	if tenantID != multitenancy.DefaultTenantIDName {
		query.Where("tenant_id = ?", tenantID)
	}

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
