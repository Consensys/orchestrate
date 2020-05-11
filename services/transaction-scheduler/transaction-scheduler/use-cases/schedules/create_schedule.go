package schedules

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/interfaces"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
)

//go:generate mockgen -source=create_schedule.go -destination=mocks/create_schedule.go -package=mocks

const createScheduleComponent = "use-cases.create-schedule"

type CreateScheduleUseCase interface {
	Execute(ctx context.Context, scheduleRequest *types.ScheduleRequest, tenantID string) (scheduleResponse *types.ScheduleResponse, err error)
}

// createScheduleUseCase is a use case to create a new transaction schedule
type createScheduleUseCase struct {
	validator validators.TransactionValidator
	db        interfaces.DB
}

// NewCreateScheduleUseCase creates a new CreateScheduleUseCase
func NewCreateScheduleUseCase(validator validators.TransactionValidator, db interfaces.DB) CreateScheduleUseCase {
	return &createScheduleUseCase{
		validator: validator,
		db:        db,
	}
}

// Execute validates and creates a new transaction schedule
func (uc *createScheduleUseCase) Execute(
	ctx context.Context,
	scheduleRequest *types.ScheduleRequest,
	tenantID string,
) (scheduleResponse *types.ScheduleResponse, err error) {
	log.WithContext(ctx).Debug("creating new schedule")

	err = utils.GetValidator().Struct(scheduleRequest)
	if err != nil {
		errMessage := "failed to validate schedule request"
		log.WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage).ExtendComponent(createScheduleComponent)
	}

	err = uc.validator.ValidateChainExists(ctx, scheduleRequest.ChainUUID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createScheduleComponent)
	}

	schedule := &models.Schedule{
		ChainUUID: scheduleRequest.ChainUUID,
		TenantID:  tenantID,
	}
	err = uc.db.Schedule().Insert(ctx, schedule)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createScheduleComponent)
	}

	log.WithContext(ctx).WithField("schedule_uuid", schedule.UUID).Info("schedule created successfully")
	return &types.ScheduleResponse{
		UUID:      schedule.UUID,
		ChainUUID: schedule.ChainUUID,
		CreatedAt: schedule.CreatedAt,
	}, nil
}
