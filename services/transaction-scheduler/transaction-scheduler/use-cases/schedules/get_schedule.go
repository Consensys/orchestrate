package schedules

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

//go:generate mockgen -source=get_schedule.go -destination=mocks/get_schedule.go -package=mocks

const getScheduleComponent = "use-cases.get-schedule"

type GetScheduleUseCase interface {
	Execute(ctx context.Context, scheduleUUID, tenantID string) (*entities.Schedule, error)
}

// getScheduleUseCase is a use case to get a schedule
type getScheduleUseCase struct {
	db store.DB
}

// NewGetScheduleUseCase creates a new GetScheduleUseCase
func NewGetScheduleUseCase(db store.DB) GetScheduleUseCase {
	return &getScheduleUseCase{
		db: db,
	}
}

// Execute gets a schedule
func (uc *getScheduleUseCase) Execute(ctx context.Context, scheduleUUID, tenantID string) (*entities.Schedule, error) {
	log.WithContext(ctx).
		WithField("schedule_uuid", scheduleUUID).
		Debug("getting schedule")

	scheduleModel, err := fetchScheduleByUUID(ctx, uc.db, scheduleUUID, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getScheduleComponent)
	}

	log.WithContext(ctx).
		WithField("schedule_uuid", scheduleModel.UUID).
		Info("schedule found successfully")

	return parsers.NewScheduleEntityFromModels(scheduleModel), nil
}

func fetchScheduleByUUID(ctx context.Context, db store.DB, scheduleUUID, tenantID string) (*models.Schedule, error) {
	schedule, err := db.Schedule().FindOneByUUID(ctx, scheduleUUID, tenantID)
	if err != nil {
		return nil, err
	}

	for idx, job := range schedule.Jobs {
		schedule.Jobs[idx], err = db.Job().FindOneByUUID(ctx, job.UUID, tenantID)
		if err != nil {
			return nil, err
		}
	}

	return schedule, nil
}
