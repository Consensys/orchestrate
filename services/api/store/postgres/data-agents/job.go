package dataagents

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"

	"github.com/go-pg/pg/v9/orm"

	gopg "github.com/go-pg/pg/v9"
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"

	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
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
		errMsg := "cannot update job with missing ID"
		log.WithContext(ctx).Error(errMsg)
		return errors.InvalidArgError(errMsg)
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

// LockOneByUUID gets a job by UUID
func (agent *PGJob) LockOneByUUID(ctx context.Context, jobUUID string) error {
	query := agent.db.ModelContext(ctx, &models.Job{}).Where("job.uuid = ?", jobUUID).For("UPDATE")
	err := pg.Select(ctx, query)
	if err != nil {
		return errors.FromError(err).ExtendComponent(jobDAComponent)
	}

	return nil
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

	if filters.ParentJobUUID != "" {
		query = query.
			Where("job.uuid = ?", filters.ParentJobUUID).
			WhereOr(fmt.Sprintf("job.internal_data @> '{\"parentJobUUID\": \"%s\"}'", filters.ParentJobUUID))
	}

	if filters.OnlyParents {
		query = query.Where("job.internal_data->'parentJobUUID' is null")
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
