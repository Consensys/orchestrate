package parsers

import (
	types "github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
)

func JobResponseToEntity(jobResponse *types.JobResponse) *entities.Job {
	// Cannot fail as the duration coming from a response is expected to be valid
	return &entities.Job{
		UUID:         jobResponse.UUID,
		ChainUUID:    jobResponse.ChainUUID,
		ScheduleUUID: jobResponse.ScheduleUUID,
		Type:         jobResponse.Type,
		Labels:       jobResponse.Labels,
		TenantID:     jobResponse.TenantID,
		OwnerID:      jobResponse.OwnerID,
		InternalData: types.FormatAnnotationsToInternalData(jobResponse.Annotations),
		Transaction:  &jobResponse.Transaction,
		Logs:         jobResponse.Logs,
		CreatedAt:    jobResponse.CreatedAt,
		UpdatedAt:    jobResponse.UpdatedAt,
	}
}
