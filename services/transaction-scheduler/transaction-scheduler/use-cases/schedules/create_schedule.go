package schedules

import (
	"context"

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
	Execute(ctx context.Context, scheduleRequest *types.ScheduleRequest, tenantID string) (scheduleResponse *types.ScheduleResponse, scheduleID int, err error)
}

// createSchedule is a use case to create a new transaction schedule
type createSchedule struct {
	chainRegistryClient client.ChainRegistryClient
	scheduleDataAgent   store.ScheduleAgent
}

// NewCreateSchedule creates a new CreateScheduleUseCase
func NewCreateSchedule(chainRegistryClient client.ChainRegistryClient, scheduleDataAgent store.ScheduleAgent) CreateScheduleUseCase {
	return &createSchedule{
		chainRegistryClient: chainRegistryClient,
		scheduleDataAgent:   scheduleDataAgent,
	}
}

// Execute validates and creates a new transaction schedule
func (usecase *createSchedule) Execute(
	ctx context.Context,
	scheduleRequest *types.ScheduleRequest,
	tenantID string,
) (scheduleResponse *types.ScheduleResponse, scheduleID int, err error) {
	log.WithContext(ctx).Debug("creating new schedule")

	// TODO: Add validation when use case becomes externally available through API

	// Validate that the chainUUID exists
	_, err = usecase.chainRegistryClient.GetChainByUUID(ctx, scheduleRequest.ChainID)
	if err != nil {
		return nil, 0, errors.FromError(err).ExtendComponent(createScheduleComponent)
	}

	schedule := &models.Schedule{
		ChainID:  scheduleRequest.ChainID,
		TenantID: tenantID,
	}
	err = usecase.scheduleDataAgent.Insert(ctx, schedule)
	if err != nil {
		return nil, 0, errors.FromError(err).ExtendComponent(createScheduleComponent)
	}

	log.WithContext(ctx).WithField("schedule_uuid", schedule.UUID).Info("schedule created successfully")
	return &types.ScheduleResponse{
		UUID:      schedule.UUID,
		ChainID:   schedule.ChainID,
		CreatedAt: schedule.CreatedAt,
	}, schedule.ID, nil
}
