package parsers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/txscheduler"
)

func JobResponseToEntity(jobResponse *txschedulertypes.JobResponse) *entities.Job {
	// Cannot fail as the duration coming from a response is expected to be valid
	return &entities.Job{
		UUID:         jobResponse.UUID,
		ChainUUID:    jobResponse.ChainUUID,
		ScheduleUUID: jobResponse.ScheduleUUID,
		Type:         jobResponse.Type,
		Labels:       jobResponse.Labels,
		TenantID:     jobResponse.TenantID,
		InternalData: txschedulertypes.FormatAnnotationsToInternalData(jobResponse.Annotations),
		Transaction:  &jobResponse.Transaction,
		Logs:         jobResponse.Logs,
		CreatedAt:    jobResponse.CreatedAt,
		UpdatedAt:    jobResponse.UpdatedAt,
	}
}
