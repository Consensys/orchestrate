package sender

import (
	"context"
	"fmt"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	usecases "github.com/consensys/orchestrate/services/tx-sender/tx-sender/use-cases"
	utils2 "github.com/consensys/orchestrate/services/tx-sender/tx-sender/utils"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const sendETHRawTxComponent = "use-cases.send-eth-raw-tx"

type sendETHRawTxUseCase struct {
	ec               ethclient.TransactionSender
	chainRegistryURL string
	jobClient        client.JobClient
	logger           *log.Logger
}

func NewSendETHRawTxUseCase(ec ethclient.TransactionSender, jobClient client.JobClient,
	chainRegistryURL string) usecases.SendETHRawTxUseCase {
	return &sendETHRawTxUseCase{
		jobClient:        jobClient,
		chainRegistryURL: chainRegistryURL,
		ec:               ec,
		logger:           log.NewLogger().SetComponent(sendETHRawTxComponent),
	}
}

func (uc *sendETHRawTxUseCase) Execute(ctx context.Context, job *entities.Job) error {
	ctx = log.With(log.WithFields(
		ctx,
		log.Field("job", job.UUID),
		log.Field("tenant_id", job.TenantID),
		log.Field("owner_id", job.OwnerID),
		log.Field("schedule_uuid", job.ScheduleUUID),
	), uc.logger)
	logger := uc.logger.WithContext(ctx)
	logger.Debug("processing ethereum raw transaction job")

	var err error
	if job.InternalData.ParentJobUUID == job.UUID || job.Status == entities.StatusPending || job.Status == entities.StatusResending {
		err = utils2.UpdateJobStatus(ctx, uc.jobClient, job, entities.StatusResending, "", nil)
	} else {
		err = utils2.UpdateJobStatus(ctx, uc.jobClient, job, entities.StatusPending, "", nil)
	}

	if err != nil {
		return errors.FromError(err).ExtendComponent(sendETHRawTxComponent)
	}

	txHash, err := uc.sendTx(ctx, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendETHRawTxComponent)
	}

	if txHash.String() != job.Transaction.Hash.String() {
		warnMessage := fmt.Sprintf("expected transaction hash %s, but got %s. Overriding", job.Transaction.Hash.String(), txHash.String())
		job.Transaction.Hash = txHash
		err = utils2.UpdateJobStatus(ctx, uc.jobClient, job, entities.StatusWarning, warnMessage, job.Transaction)
		if err != nil {
			return errors.FromError(err).ExtendComponent(sendETHRawTxComponent)
		}
	}

	logger.Info("ethereum raw transaction job was processed successfully")
	return nil
}

func (uc *sendETHRawTxUseCase) sendTx(ctx context.Context, job *entities.Job) (*ethcommon.Hash, error) {
	proxyURL := utils.GetProxyURL(uc.chainRegistryURL, job.ChainUUID)
	txHash, err := uc.ec.SendRawTransaction(ctx, proxyURL, job.Transaction.Raw)
	if err != nil {
		errMsg := "cannot send ethereum raw transaction"
		uc.logger.WithContext(ctx).WithError(err).Error(errMsg)
		return nil, err
	}

	return &txHash, nil
}
