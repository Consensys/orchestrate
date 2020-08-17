package dataagents

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	"github.com/go-pg/pg/v9/orm"

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
func (agent *PGJob) FindOneByUUID(ctx context.Context, jobUUID string, tenants []string) (*models.Job, error) {
	job := &models.Job{}

	query := agent.db.ModelContext(ctx, job).
		Where("job.uuid = ?", jobUUID).
		Relation("Transaction").
		Relation("Schedule").
		Relation("Logs", func(q *orm.Query) (*orm.Query, error) {
			return q.Order("id ASC"), nil
		})

	query = pg.WhereAllowedTenants(query, "schedule.tenant_id", tenants).Order("id ASC")

	err := pg.Select(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(jobDAComponent)
	}

	return job, nil
}

func (agent *PGJob) Search(ctx context.Context, filters *entities.JobFilters, tenants []string) ([]*models.Job, error) {
	var jobs []*models.Job

	query := agent.db.ModelContext(ctx, &jobs).
		Relation("Transaction").
		Relation("Schedule").
		Relation("Logs", func(q *orm.Query) (*orm.Query, error) {
			return q.Order("id ASC"), nil
		})

	if len(filters.TxHashes) > 0 {
		query = query.Where("transaction.hash in (?)", gopg.In(filters.TxHashes))
	}

	if filters.ChainUUID != "" {
		query = query.Where("job.chain_uuid = ?", filters.ChainUUID)
	}

	if filters.Status != "" {
		query = query.
			Join("LEFT JOIN logs as log").
			JoinOn("log.job_id = job.id").
			Join("LEFT JOIN logs as tmpl").
			JoinOn("tmpl.job_id = job.id AND log.created_at < tmpl.created_at").
			Where("tmpl.id is null AND log.status = ?", filters.Status)
	}

	query = pg.WhereAllowedTenants(query, "schedule.tenant_id", tenants).
		Where("job.updated_at > ?", filters.UpdatedAfter). // No need to check this filter as the zero value is 1/1/0001 0h:0m:0s
		Order("id ASC")

	err := pg.Select(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(jobDAComponent)
	}

	return jobs, nil
}
