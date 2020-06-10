package dataagents

import (
	"context"

	gopg "github.com/go-pg/pg/v9"
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"

	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

const jobDAComponent = "data-agents.job"

// PGJob is a job data agent for PostgreSQL
type PGJob struct {
	db pg.DB
}

// NewPGJob creates a new pgJob
func NewPGJob(db pg.DB) *PGJob {
	return &PGJob{db: db}
}

// Insert Inserts a new job in DB
func (agent *PGJob) Insert(ctx context.Context, job *models.Job) error {
	if job.UUID == "" {
		job.UUID = uuid.Must(uuid.NewV4()).String()
	}

	if job.Transaction != nil && job.TransactionID == nil {
		job.TransactionID = &job.Transaction.ID
	}
	if job.Schedule != nil && job.ScheduleID == nil {
		job.ScheduleID = &job.Transaction.ID
	}

	agent.db.ModelContext(ctx, job)
	err := pg.Insert(ctx, agent.db, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(jobDAComponent)
	}

	return nil
}

// Insert Inserts a new job in DB
func (agent *PGJob) Update(ctx context.Context, job *models.Job) error {
	if job.ID == 0 {
		return errors.InvalidArgError("cannot update job with missing ID")
	}

	if job.Transaction != nil && job.TransactionID == nil {
		job.TransactionID = &job.Transaction.ID
	}
	if job.Schedule != nil && job.ScheduleID == nil {
		job.ScheduleID = &job.Transaction.ID
	}

	agent.db.ModelContext(ctx, job)
	err := pg.Update(ctx, agent.db, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(jobDAComponent)
	}

	return nil
}

// FindOneByUUID gets a job by UUID
func (agent *PGJob) FindOneByUUID(ctx context.Context, jobUUID, tenantID string) (*models.Job, error) {
	job := &models.Job{}

	query := agent.db.ModelContext(ctx, job).
		Where("job.uuid = ?", jobUUID).
		Relation("Transaction").
		Relation("Schedule").
		Relation("Logs")

	if tenantID != "" {
		query.Where("schedule.tenant_id = ?", tenantID)
	}

	err := pg.Select(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(jobDAComponent)
	}

	return job, nil
}

func (agent *PGJob) Search(ctx context.Context, tenantID string, txHashes []string, chainUUID string) ([]*models.Job, error) {
	jobs := []*models.Job{}

	query := agent.db.ModelContext(ctx, &jobs).
		Relation("Transaction").
		Relation("Schedule").
		Relation("Logs")

	if len(txHashes) > 0 {
		query = query.Where("transaction.hash in (?)", gopg.In(txHashes))
	}

	if chainUUID != "" {
		query = query.Where("job.chain_uuid = ?", chainUUID)
	}

	if tenantID != "" {
		query.Where("schedule.tenant_id = ?", tenantID)
	}

	err := pg.Select(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(jobDAComponent)
	}

	return jobs, nil
}
