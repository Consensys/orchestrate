package parsers

import (
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
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
		InternalData: types.FormatAnnotationsToInternalData(jobResponse.Annotations),
		Transaction:  &jobResponse.Transaction,
		Logs:         jobResponse.Logs,
		CreatedAt:    jobResponse.CreatedAt,
		UpdatedAt:    jobResponse.UpdatedAt,
	}
}
