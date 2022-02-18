package dataagents

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	pg "github.com/consensys/orchestrate/pkg/toolkit/database/postgres"
	"github.com/consensys/orchestrate/services/api/store"
	"github.com/go-pg/pg/v9/orm"
	"github.com/gofrs/uuid"

	"github.com/consensys/orchestrate/services/api/store/models"
)

const scheduleDAComponent = "data-agents.schedule"

// PGSchedule is a schedule data agent for PostgreSQL
type PGSchedule struct {
	db     pg.DB
	logger *log.Logger
}

// NewPGSchedule creates a new PGSchedule
func NewPGSchedule(db pg.DB) store.ScheduleAgent {
	return &PGSchedule{db: db, logger: log.NewLogger().SetComponent(scheduleDAComponent)}
}

// Insert Inserts a new schedule in DB
func (agent *PGSchedule) Insert(ctx context.Context, schedule *models.Schedule) error {
	if schedule.UUID == "" {
		schedule.UUID = uuid.Must(uuid.NewV4()).String()
	}

	err := pg.Insert(ctx, agent.db, schedule)
	if err != nil {
		agent.logger.WithContext(ctx).WithError(err).Error("failed to insert schedule")
		return errors.FromError(err).ExtendComponent(scheduleDAComponent)
	}

	return nil
}

// FindOneByUUID Finds a schedule in DB
func (agent *PGSchedule) FindOneByUUID(ctx context.Context, scheduleUUID string, tenants []string, ownerID string) (*models.Schedule, error) {
	schedule := &models.Schedule{}

	query := agent.db.ModelContext(ctx, schedule).
		Relation("Jobs", func(q *orm.Query) (*orm.Query, error) {
			return q.Order("id ASC"), nil
		}).
		Where("schedule.uuid = ?", scheduleUUID)

	query = pg.WhereAllowedTenants(query, "schedule.tenant_id", tenants)
	query = pg.WhereAllowedOwner(query, "owner_id", ownerID)

	if err := pg.SelectOne(ctx, query); err != nil {
		if !errors.IsNotFoundError(err) {
			agent.logger.WithContext(ctx).WithError(err).Error("failed to find schedule")
		}
		return nil, errors.FromError(err).ExtendComponent(scheduleDAComponent)
	}

	return schedule, nil
}

// Search Finds schedules in DB
func (agent *PGSchedule) FindAll(ctx context.Context, tenants []string, ownerID string) ([]*models.Schedule, error) {
	var schedules []*models.Schedule

	query := agent.db.ModelContext(ctx, &schedules).
		Relation("Jobs", func(q *orm.Query) (*orm.Query, error) {
			return q.Order("id ASC"), nil
		})

	query = pg.WhereAllowedTenants(query, "schedule.tenant_id", tenants)
	query = pg.WhereAllowedOwner(query, "owner_id", ownerID)

	err := pg.Select(ctx, query)
	if err != nil {
		if !errors.IsNotFoundError(err) {
			agent.logger.WithContext(ctx).WithError(err).Error("failed to insert all schedules")
		}
		return nil, errors.FromError(err).ExtendComponent(jobDAComponent)
	}

	return schedules, nil
}
