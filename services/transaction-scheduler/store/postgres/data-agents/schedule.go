package dataagents

import (
	"context"

	"github.com/go-pg/pg/v9/orm"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"

	"github.com/go-pg/pg/v9"
)

const scheduleDAComponent = "data-agents.schedule"

// PGSchedule is a schedule data agent for PostgreSQL
type PGSchedule struct {
	db orm.DB
}

// NewPGSchedule creates a new PGSchedule
func NewPGSchedule(db orm.DB) *PGSchedule {
	return &PGSchedule{db: db}
}

// Insert Inserts a new schedule in DB
func (agent *PGSchedule) Insert(ctx context.Context, schedule *models.Schedule) error {
	schedule.UUID = uuid.NewV4().String()
	err := postgres.Insert(ctx, agent.db, schedule)
	if err != nil {
		return errors.FromError(err).ExtendComponent(scheduleDAComponent)
	}

	return nil
}

// FindOneByUUID Finds a schedule in DB
func (agent *PGSchedule) FindOneByUUID(ctx context.Context, scheduleUUID, tenantID string) (*models.Schedule, error) {
	schedule := &models.Schedule{}
	logger := log.WithField("schedule_uuid", scheduleUUID).WithField("tenant_id", tenantID)

	query := agent.db.ModelContext(ctx, schedule).Where("schedule.uuid = ?", scheduleUUID)
	if tenantID != multitenancy.DefaultTenantIDName {
		query.Where("schedule.tenant_id = ?", tenantID)
	}

	err := query.Relation("Jobs").Select()
	if err != nil && err == pg.ErrNoRows {
		errMessage := "schedule does not exist"
		logger.WithError(err).Error(errMessage)
		return nil, errors.NotFoundError(errMessage).ExtendComponent(jobDAComponent)
	} else if err != nil {
		errMessage := "could not load schedule"
		logger.WithError(err).Errorf(errMessage)
		return nil, errors.PostgresConnectionError(errMessage).ExtendComponent(jobDAComponent)
	}

	return schedule, nil
}
