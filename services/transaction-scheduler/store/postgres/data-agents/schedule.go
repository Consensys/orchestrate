package dataagents

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"

	"github.com/go-pg/pg/v9"
)

const scheduleDAComponent = "data-agents.schedule"

// PGSchedule is a schedule data agent for PostgreSQL
type PGSchedule struct {
	db *pg.DB
}

// NewPGSchedule creates a new PGSchedule
func NewPGSchedule(db *pg.DB) *PGSchedule {
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
