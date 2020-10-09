package schedules

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

const getScheduleComponent = "use-cases.get-schedule"

// getScheduleUseCase is a use case to get a schedule
type getScheduleUseCase struct {
	db store.DB
}

// NewGetScheduleUseCase creates a new GetScheduleUseCase
func NewGetScheduleUseCase(db store.DB) usecases.GetScheduleUseCase {
	return &getScheduleUseCase{
		db: db,
	}
}

// Execute gets a schedule
func (uc *getScheduleUseCase) Execute(ctx context.Context, scheduleUUID string, tenants []string) (*entities.Schedule, error) {
	log.WithContext(ctx).
		WithField("schedule_uuid", scheduleUUID).
		Debug("getting schedule")

	scheduleModel, err := fetchScheduleByUUID(ctx, uc.db, scheduleUUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getScheduleComponent)
	}

	log.WithContext(ctx).
		WithField("schedule_uuid", scheduleModel.UUID).
		Info("schedule found successfully")

	return parsers.NewScheduleEntityFromModels(scheduleModel), nil
}

func fetchScheduleByUUID(ctx context.Context, db store.DB, scheduleUUID string, tenants []string) (*models.Schedule, error) {
	schedule, err := db.Schedule().FindOneByUUID(ctx, scheduleUUID, tenants)
	if err != nil {
		return nil, err
	}

	for idx, job := range schedule.Jobs {
		schedule.Jobs[idx], err = db.Job().FindOneByUUID(ctx, job.UUID, tenants)
		if err != nil {
			return nil, err
		}
	}

	return schedule, nil
}
