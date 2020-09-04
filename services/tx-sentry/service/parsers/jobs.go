package parsers

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx-scheduler"
)

func JobResponseToEntity(jobResponse *txschedulertypes.JobResponse) *entities.Job {
	// Cannot fail as the duration coming from a response is expected to be valid
	retryInterval, _ := time.ParseDuration(jobResponse.Annotations.GasPricePolicy.RetryPolicy.Interval)

	return &entities.Job{
		UUID:         jobResponse.UUID,
		ChainUUID:    jobResponse.ChainUUID,
		ScheduleUUID: jobResponse.ScheduleUUID,
		Type:         jobResponse.Type,
		Labels:       jobResponse.Labels,
		InternalData: &entities.InternalData{
			OneTimeKey:        jobResponse.Annotations.OneTimeKey,
			Priority:          jobResponse.Annotations.GasPricePolicy.Priority,
			RetryInterval:     retryInterval,
			GasPriceIncrement: jobResponse.Annotations.GasPricePolicy.RetryPolicy.Increment,
			GasPriceLimit:     jobResponse.Annotations.GasPricePolicy.RetryPolicy.Limit,
			ParentJobUUID:     jobResponse.ParentJobUUID,
		},
		Transaction: &jobResponse.Transaction,
		Logs:        jobResponse.Logs,
		CreatedAt:   jobResponse.CreatedAt,
		UpdatedAt:   jobResponse.UpdatedAt,
	}
}
