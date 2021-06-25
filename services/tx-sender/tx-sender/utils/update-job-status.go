package utils

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/sdk/client"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/pkg/types/api"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
)

func UpdateJobStatus(ctx context.Context, apiClient client.JobClient, job *entities.Job, status entities.JobStatus,
	msg string, transaction *entities.ETHTransaction) error {
	logger := log.FromContext(ctx).WithField("status", status)

	txUpdateReq := &api.UpdateJobRequest{
		Status:      status,
		Message:     msg,
		Transaction: transaction,
	}

	_, err := apiClient.UpdateJob(ctx, job.UUID, txUpdateReq)
	if err != nil {
		logger.WithError(err).Error("failed to update job status")
		return err
	}

	job.Status = status
	logger.Debug("job status was updated successfully")
	return nil
}
