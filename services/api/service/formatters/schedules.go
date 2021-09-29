package formatters

import (
	types "github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
)

func FormatScheduleResponse(schedule *entities.Schedule) *types.ScheduleResponse {
	scheduleResponse := &types.ScheduleResponse{
		UUID:      schedule.UUID,
		TenantID:  schedule.TenantID,
		CreatedAt: schedule.CreatedAt,
		Jobs:      []*types.JobResponse{},
	}

	for idx := range schedule.Jobs {
		scheduleResponse.Jobs = append(scheduleResponse.Jobs, FormatJobResponse(schedule.Jobs[idx]))
	}

	return scheduleResponse
}
