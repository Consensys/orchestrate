package formatters

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx-scheduler"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

func FormatJobResponse(job *entities.Job) *types.JobResponse {
	return &types.JobResponse{
		UUID:         job.UUID,
		ChainUUID:    job.ChainUUID,
		ScheduleUUID: job.ScheduleUUID,
		NextJobUUID:  job.NextJobUUID,
		Transaction:  *job.Transaction,
		Logs:         job.Logs,
		Labels:       job.Labels,
		Annotations: types.Annotations{
			OneTimeKey: job.InternalData.OneTimeKey,
			Priority:   job.InternalData.Priority,
			RetryPolicy: types.GasPriceRetryParams{
				Interval:       fmt.Sprintf("%vs", job.InternalData.RetryInterval.Seconds()),
				IncrementLevel: job.InternalData.GasPriceIncrementLevel,
				Increment:      job.InternalData.GasPriceIncrement,
				Limit:          job.InternalData.GasPriceLimit,
			},
		},
		Type:      job.Type,
		Status:    job.GetStatus(),
		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
	}
}

func FormatJobCreateRequest(request *types.CreateJobRequest, defaultRetryInterval time.Duration) *entities.Job {
	return &entities.Job{
		ChainUUID:    request.ChainUUID,
		ScheduleUUID: request.ScheduleUUID,
		NextJobUUID:  request.NextJobUUID,
		Type:         request.Type,
		Labels:       request.Labels,
		InternalData: formatAnnotations(&request.Annotations, defaultRetryInterval),
		Transaction:  &request.Transaction,
	}
}

func FormatJobUpdateRequest(request *types.UpdateJobRequest) *entities.Job {
	job := &entities.Job{
		Labels:      request.Labels,
		Transaction: request.Transaction,
	}

	if request.Annotations != nil {
		job.InternalData = &entities.InternalData{
			OneTimeKey:             request.Annotations.OneTimeKey,
			Priority:               request.Annotations.Priority,
			GasPriceIncrementLevel: request.Annotations.RetryPolicy.IncrementLevel,
			GasPriceIncrement:      request.Annotations.RetryPolicy.Increment,
			GasPriceLimit:          request.Annotations.RetryPolicy.Limit,
		}

		if request.Annotations.RetryPolicy.Interval != "" {
			// we can skip the error check as at this point we know the interval is a duration as it already passed validation
			job.InternalData.RetryInterval, _ = time.ParseDuration(request.Annotations.RetryPolicy.Interval)
		}
	}

	return job
}

func FormatJobFilterRequest(req *http.Request) (*entities.JobFilters, error) {
	filters := &entities.JobFilters{}

	qTxHashes := req.URL.Query().Get("tx_hashes")
	if qTxHashes != "" {
		filters.TxHashes = strings.Split(qTxHashes, ",")
	}

	qChainUUID := req.URL.Query().Get("chain_uuid")
	if qChainUUID != "" {
		filters.ChainUUID = qChainUUID
	}

	qStatus := req.URL.Query().Get("status")
	if qStatus != "" {
		filters.Status = qStatus
	}

	if err := utils.GetValidator().Struct(filters); err != nil {
		return nil, err
	}

	return filters, nil
}

func formatAnnotations(annotations *types.Annotations, defaultRetryInterval time.Duration) *entities.InternalData {
	internalData := &entities.InternalData{
		OneTimeKey:             annotations.OneTimeKey,
		Priority:               annotations.Priority,
		GasPriceIncrementLevel: annotations.RetryPolicy.IncrementLevel,
		GasPriceIncrement:      annotations.RetryPolicy.Increment,
		GasPriceLimit:          annotations.RetryPolicy.Limit,
	}

	if annotations.RetryPolicy.Interval == "" {
		internalData.RetryInterval = defaultRetryInterval
	} else {
		// we can skip the error check as at this point we know the interval is a duration as it already passed validation
		internalData.RetryInterval, _ = time.ParseDuration(annotations.RetryPolicy.Interval)
	}

	if internalData.Priority == "" {
		internalData.Priority = utils.PriorityMedium
	}

	return internalData
}
