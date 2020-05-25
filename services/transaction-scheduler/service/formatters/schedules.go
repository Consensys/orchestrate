package formatters

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

func FormatScheduleResponse(scheduleModel *entities.Schedule) *types.ScheduleResponse {
	scheduleResponse := &types.ScheduleResponse{
		UUID:      scheduleModel.UUID,
		ChainUUID: scheduleModel.ChainUUID,
		CreatedAt: scheduleModel.CreatedAt,
		Jobs:      []*types.JobResponse{},
	}

	for idx := range scheduleModel.Jobs {
		scheduleResponse.Jobs = append(scheduleResponse.Jobs, FormatJobResponse(scheduleModel.Jobs[idx]))
	}

	return scheduleResponse
}

func FormatScheduleCreateRequest(request *types.CreateScheduleRequest) *entities.Schedule {
	scheduleResponse := &entities.Schedule{
		ChainUUID: request.ChainUUID,
	}

	return scheduleResponse
}
