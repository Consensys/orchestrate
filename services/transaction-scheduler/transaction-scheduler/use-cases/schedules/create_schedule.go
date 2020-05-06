package schedules

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
)

//go:generate mockgen -source=create_schedule.go -destination=mocks/create_schedule.go -package=mocks

const createScheduleComponent = "use-cases.create-schedule"

type CreateScheduleUseCase interface {
	Execute(ctx context.Context, scheduleRequest *types.ScheduleRequest, tenantID string) (scheduleResponse *types.ScheduleResponse, err error)
}

// createScheduleUseCase is a use case to create a new transaction schedule
type createScheduleUseCase struct {
	chainRegistryClient client.ChainRegistryClient
	scheduleDataAgent   store.ScheduleAgent
}

// NewCreateScheduleUseCase creates a new CreateScheduleUseCase
func NewCreateScheduleUseCase(chainRegistryClient client.ChainRegistryClient, scheduleDataAgent store.ScheduleAgent) CreateScheduleUseCase {
	return &createScheduleUseCase{
		chainRegistryClient: chainRegistryClient,
		scheduleDataAgent:   scheduleDataAgent,
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

	// Validate that the chainUUID exists
	_, err = uc.chainRegistryClient.GetChainByUUID(ctx, scheduleRequest.ChainID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createScheduleComponent)
	}

	schedule := &models.Schedule{
		ChainID:  scheduleRequest.ChainID,
		TenantID: tenantID,
	}
	err = uc.scheduleDataAgent.Insert(ctx, schedule)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createScheduleComponent)
	}

	log.WithContext(ctx).WithField("schedule_uuid", schedule.UUID).Info("schedule created successfully")
	return &types.ScheduleResponse{
		UUID:      schedule.UUID,
		ChainID:   schedule.ChainID,
		CreatedAt: schedule.CreatedAt,
	}, nil
}
