package schedules

import (
	"context"

	"github.com/consensys/orchestrate/pkg/types/entities"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	"github.com/consensys/orchestrate/services/api/store"
)

const createScheduleComponent = "use-cases.create-schedule"

// createScheduleUseCase is a use case to create a new transaction schedule
type createScheduleUseCase struct {
	db     store.DB
	logger *log.Logger
}

// NewCreateScheduleUseCase creates a new CreateScheduleUseCase
func NewCreateScheduleUseCase(db store.DB) usecases.CreateScheduleUseCase {
	return &createScheduleUseCase{
		db:     db,
		logger: log.NewLogger().SetComponent(createScheduleComponent),
	}
}

func (uc createScheduleUseCase) WithDBTransaction(dbtx store.Tx) usecases.CreateScheduleUseCase {
	uc.db = dbtx
	return &uc
}

// Execute validates and creates a new transaction schedule
func (uc *createScheduleUseCase) Execute(ctx context.Context, schedule *entities.Schedule) (*entities.Schedule, error) {
	logger := uc.logger.WithContext(ctx)
	logger.Debug("creating new schedule")

	scheduleModel := parsers.NewScheduleModelFromEntities(schedule)

	if err := uc.db.Schedule().Insert(ctx, scheduleModel); err != nil {
		return nil, errors.FromError(err).ExtendComponent(createScheduleComponent)
	}

	logger.WithField("schedule", scheduleModel.UUID).Info("schedule created successfully")

	return parsers.NewScheduleEntityFromModels(scheduleModel), nil
}
