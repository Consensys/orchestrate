package utils

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/txscheduler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
)

func UpdateJobStatus(ctx context.Context, txClient client.TransactionSchedulerClient, jobUUID, status, msg string,
	transaction *entities.ETHTransaction) error {
	logger := log.WithContext(ctx).WithField("job_uuid", jobUUID).WithField("status", status)
	logger.Debug("updating job status")

	txUpdateReq := &txschedulertypes.UpdateJobRequest{
		Status:      status,
		Message:     msg,
		Transaction: transaction,
	}

	_, err := txClient.UpdateJob(ctx, jobUUID, txUpdateReq)
	if err != nil {
		errMsg := "failed to update job status"
		logger.WithError(err).Errorf(errMsg)
		return err
	}

	logger.Info("job status updated successfully")
	return nil
}
