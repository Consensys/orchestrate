package schedules

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
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
