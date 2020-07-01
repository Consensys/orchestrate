package dataagents

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	"github.com/gofrs/uuid"
	pg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

const scheduleDAComponent = "data-agents.schedule"

// PGSchedule is a schedule data agent for PostgreSQL
type PGSchedule struct {
	db pg.DB
}

// NewPGSchedule creates a new PGSchedule
func NewPGSchedule(db pg.DB) *PGSchedule {
	return &PGSchedule{db: db}
}

// Insert Inserts a new schedule in DB
func (agent *PGSchedule) Insert(ctx context.Context, schedule *models.Schedule) error {
	if schedule.UUID == "" {
		schedule.UUID = uuid.Must(uuid.NewV4()).String()
	}

	if schedule.TransactionRequest != nil && schedule.TransactionRequestID == nil {
		schedule.TransactionRequestID = &schedule.TransactionRequest.ID
	}

	err := pg.Insert(ctx, agent.db, schedule)
	if err != nil {
		return errors.FromError(err).ExtendComponent(scheduleDAComponent)
	}

	return nil
}

// FindOneByUUID Finds a schedule in DB
func (agent *PGSchedule) FindOneByUUID(ctx context.Context, scheduleUUID string, tenants []string) (*models.Schedule, error) {
	schedule := &models.Schedule{}

	query := agent.db.ModelContext(ctx, schedule).
		Relation("Jobs").
		Where("schedule.uuid = ?", scheduleUUID)

	query = pg.WhereAllowedTenants(query, "schedule.tenant_id", tenants)

	if err := pg.SelectOne(ctx, query); err != nil {
		return nil, errors.FromError(err).ExtendComponent(scheduleDAComponent)
	}

	return schedule, nil
}

// FindOneByUUID Finds a schedule in DB
func (agent *PGSchedule) FindAll(ctx context.Context, tenants []string) ([]*models.Schedule, error) {
	schedules := []*models.Schedule{}

	query := agent.db.ModelContext(ctx, &schedules).
		Relation("Jobs")

	query = pg.WhereAllowedTenants(query, "schedule.tenant_id", tenants)

	err := pg.Select(ctx, query)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(jobDAComponent)
	}

	return schedules, nil
}
