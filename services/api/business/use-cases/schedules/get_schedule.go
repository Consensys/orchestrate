package schedules

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/services/api/business/parsers"
	"github.com/ConsenSys/orchestrate/services/api/store"
	"github.com/ConsenSys/orchestrate/services/api/store/models"
)

const getScheduleComponent = "use-cases.get-schedule"

// getScheduleUseCase is a use case to get a schedule
type getScheduleUseCase struct {
	db     store.DB
	logger *log.Logger
}

// NewGetScheduleUseCase creates a new GetScheduleUseCase
func NewGetScheduleUseCase(db store.DB) usecases.GetScheduleUseCase {
	return &getScheduleUseCase{
		db:     db,
		logger: log.NewLogger().SetComponent(getScheduleComponent),
	}
}

// Execute gets a schedule
func (uc *getScheduleUseCase) Execute(ctx context.Context, scheduleUUID string, tenants []string) (*entities.Schedule, error) {
	ctx = log.WithFields(ctx, log.Field("schedule", scheduleUUID))

	scheduleModel, err := fetchScheduleByUUID(ctx, uc.db, scheduleUUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getScheduleComponent)
	}

	uc.logger.Debug("schedule found successfully")
	return parsers.NewScheduleEntityFromModels(scheduleModel), nil
}

func fetchScheduleByUUID(ctx context.Context, db store.DB, scheduleUUID string, tenants []string) (*models.Schedule, error) {
	schedule, err := db.Schedule().FindOneByUUID(ctx, scheduleUUID, tenants)
	if err != nil {
		return nil, err
	}

	for idx, job := range schedule.Jobs {
		schedule.Jobs[idx], err = db.Job().FindOneByUUID(ctx, job.UUID, tenants, false)
		if err != nil {
			return nil, err
		}
	}

	return schedule, nil
}
