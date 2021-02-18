package utils

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

func UpdateJobStatus(ctx context.Context, apiClient client.JobClient, jobUUID string, status entities.JobStatus,
	msg string, transaction *entities.ETHTransaction) error {
	logger := log.FromContext(ctx).WithField("status", status)

	txUpdateReq := &api.UpdateJobRequest{
		Status:      status,
		Message:     msg,
		Transaction: transaction,
	}

	_, err := apiClient.UpdateJob(ctx, jobUUID, txUpdateReq)
	if err != nil {
		logger.WithError(err).Error("failed to update job status")
		return err
	}

	logger.Debug("job status was updated successfully")
	return nil
}
