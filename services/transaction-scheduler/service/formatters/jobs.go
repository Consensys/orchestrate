package formatters

import (
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/txscheduler"

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
		TenantID:     job.TenantID,
		Annotations: types.Annotations{
			OneTimeKey:     job.InternalData.OneTimeKey,
			GasPricePolicy: FormatJobAnnotationsGasPricePolicy(job.InternalData),
		},
		Type:          job.Type,
		Status:        job.GetStatus(),
		ParentJobUUID: job.InternalData.ParentJobUUID,
		CreatedAt:     job.CreatedAt,
		UpdatedAt:     job.UpdatedAt,
	}
}

func FormatJobAnnotationsGasPricePolicy(data *entities.InternalData) types.GasPriceParams {
	gasPricePolicy := types.GasPriceParams{
		Priority: data.Priority,
		RetryPolicy: types.RetryParams{
			Increment: data.GasPriceIncrement,
			Limit:     data.GasPriceLimit,
		},
	}

	if data.RetryInterval != 0 {
		gasPricePolicy.RetryPolicy.Interval = data.RetryInterval.String()
	}

	return gasPricePolicy
}

func FormatJobCreateRequest(request *types.CreateJobRequest) *entities.Job {
	return &entities.Job{
		ChainUUID:    request.ChainUUID,
		ScheduleUUID: request.ScheduleUUID,
		NextJobUUID:  request.NextJobUUID,
		Type:         request.Type,
		Labels:       request.Labels,
		InternalData: formatInternalData(request.Annotations.OneTimeKey,
			&request.Annotations.GasPricePolicy,
			request.ParentJobUUID,
		),
		Transaction: &request.Transaction,
	}
}

func FormatJobUpdateRequest(request *types.UpdateJobRequest) *entities.Job {
	job := &entities.Job{
		Labels:      request.Labels,
		Transaction: request.Transaction,
	}

	if request.Annotations != nil {
		job.InternalData = &entities.InternalData{
			OneTimeKey:        request.Annotations.OneTimeKey,
			Priority:          request.Annotations.GasPricePolicy.Priority,
			GasPriceIncrement: request.Annotations.GasPricePolicy.RetryPolicy.Increment,
			GasPriceLimit:     request.Annotations.GasPricePolicy.RetryPolicy.Limit,
		}

		if request.Annotations.GasPricePolicy.RetryPolicy.Interval != "" {
			// we can skip the error check as at this point we know the interval is a duration as it already passed validation
			job.InternalData.RetryInterval, _ = time.ParseDuration(request.Annotations.GasPricePolicy.RetryPolicy.Interval)
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

	qParentJobUUID := req.URL.Query().Get("parent_job_uuid")
	if qParentJobUUID != "" {
		filters.ParentJobUUID = qParentJobUUID
	}

	qUpdatedAfter := req.URL.Query().Get("updated_after")
	if qUpdatedAfter != "" {
		updatedAfter, err := time.Parse(time.RFC3339, qUpdatedAfter)
		if err != nil {
			errMessage := "failed to parse updated_after as time"
			log.WithError(err).WithField("updated_after", qUpdatedAfter).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		filters.UpdatedAfter = updatedAfter
	}

	qOnlyParents := req.URL.Query().Get("only_parents")
	if qOnlyParents == "true" {
		filters.OnlyParents = true
	}

	if err := utils.GetValidator().Struct(filters); err != nil {
		return nil, err
	}

	return filters, nil
}

func formatInternalData(oneTimeKey bool, gasPricePolicy *types.GasPriceParams, parentJobUUID string) *entities.InternalData {
	internalData := &entities.InternalData{
		OneTimeKey:        oneTimeKey,
		Priority:          gasPricePolicy.Priority,
		GasPriceIncrement: gasPricePolicy.RetryPolicy.Increment,
		GasPriceLimit:     gasPricePolicy.RetryPolicy.Limit,
		ParentJobUUID:     parentJobUUID,
	}

	if gasPricePolicy.RetryPolicy.Interval != "" {
		// we can skip the error check as at this point we know the interval is a duration as it already passed validation
		internalData.RetryInterval, _ = time.ParseDuration(gasPricePolicy.RetryPolicy.Interval)
	}

	if internalData.Priority == "" {
		internalData.Priority = utils.PriorityMedium
	}

	return internalData
}
