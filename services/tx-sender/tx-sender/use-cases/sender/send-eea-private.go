package sender

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/orchestrate/services/tx-sender/tx-sender/nonce"
	usecases "github.com/consensys/orchestrate/services/tx-sender/tx-sender/use-cases"
	utils2 "github.com/consensys/orchestrate/services/tx-sender/tx-sender/utils"
)

const sendEEAPrivateTxComponent = "use-cases.send-eea-private-tx"

type sendEEAPrivateTxUseCase struct {
	crafter          usecases.CraftTransactionUseCase
	signTx           usecases.SignETHTransactionUseCase
	nonceManager     nonce.Manager
	jobClient        client.JobClient
	ec               ethclient.EEATransactionSender
	chainRegistryURL string
	logger           *log.Logger
}

func NewSendEEAPrivateTxUseCase(signTx usecases.SignEEATransactionUseCase, crafter usecases.CraftTransactionUseCase,
	ec ethclient.EEATransactionSender, jobClient client.JobClient, chainRegistryURL string,
	nonceManager nonce.Manager) usecases.SendEEAPrivateTxUseCase {
	return &sendEEAPrivateTxUseCase{
		jobClient:        jobClient,
		chainRegistryURL: chainRegistryURL,
		signTx:           signTx,
		ec:               ec,
		nonceManager:     nonceManager,
		crafter:          crafter,
		logger:           log.NewLogger().SetComponent(sendEEAPrivateTxComponent),
	}
}

func (uc *sendEEAPrivateTxUseCase) Execute(ctx context.Context, job *entities.Job) error {
	ctx = log.With(log.WithFields(
		ctx,
		log.Field("job", job.UUID),
		log.Field("tenant_id", job.TenantID),
		log.Field("schedule_uuid", job.ScheduleUUID),
	), uc.logger)
	logger := uc.logger.WithContext(ctx)

	logger.Debug("processing EEA private transaction job")

	err := uc.crafter.Execute(ctx, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendEEAPrivateTxComponent)
	}

	job.Transaction.Raw, _, err = uc.signTx.Execute(ctx, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendEEAPrivateTxComponent)
	}

	job.Transaction.Hash, err = uc.sendTx(ctx, job)
	if err != nil {
		if err2 := uc.nonceManager.CleanNonce(ctx, job, err); err2 != nil {
			return errors.FromError(err2).ExtendComponent(sendEEAPrivateTxComponent)
		}
		return err
	}

	err = uc.nonceManager.IncrementNonce(ctx, job)
	if err != nil {
		return err
	}

	err = utils2.UpdateJobStatus(ctx, uc.jobClient, job, entities.StatusStored, "", job.Transaction)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendEEAPrivateTxComponent)
	}

	logger.Info("EEA private transaction job was sent successfully")
	return nil
}

func (uc *sendEEAPrivateTxUseCase) sendTx(ctx context.Context, job *entities.Job) (string, error) {
	proxyURL := utils.GetProxyURL(uc.chainRegistryURL, job.ChainUUID)
	txHash, err := uc.ec.PrivDistributeRawTransaction(ctx, proxyURL, job.Transaction.Raw)
	if err != nil {
		errMsg := "cannot send EEA private transaction"
		uc.logger.WithContext(ctx).WithError(err).Error(errMsg)
		return "", err
	}

	return txHash.String(), nil
}
