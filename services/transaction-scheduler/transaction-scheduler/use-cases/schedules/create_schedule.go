package schedules

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

const createScheduleComponent = "use-cases.create-schedule"

// createScheduleUseCase is a use case to create a new transaction schedule
type createScheduleUseCase struct {
	db store.DB
}

// NewCreateScheduleUseCase creates a new CreateScheduleUseCase
func NewCreateScheduleUseCase(db store.DB) usecases.CreateScheduleUseCase {
	return &createScheduleUseCase{
		db: db,
	}
}

func (uc createScheduleUseCase) WithDBTransaction(dbtx store.Tx) usecases.CreateScheduleUseCase {
	uc.db = dbtx
	return &uc
}

// Execute validates and creates a new transaction schedule
func (uc *createScheduleUseCase) Execute(ctx context.Context, schedule *entities.Schedule) (*entities.Schedule, error) {
	log.WithContext(ctx).Debug("creating new schedule")

	scheduleModel := parsers.NewScheduleModelFromEntities(schedule)

	if err := uc.db.Schedule().Insert(ctx, scheduleModel); err != nil {
		return nil, errors.FromError(err).ExtendComponent(createScheduleComponent)
	}

	log.WithContext(ctx).WithField("schedule_uuid", scheduleModel.UUID).Info("schedule created successfully")

	return parsers.NewScheduleEntityFromModels(scheduleModel), nil
}
