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
	"github.com/consensys/orchestrate/services/tx-sender/tx-sender/nonce"
	usecases "github.com/consensys/orchestrate/services/tx-sender/tx-sender/use-cases"
	utils2 "github.com/consensys/orchestrate/services/tx-sender/tx-sender/utils"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const sendETHTxComponent = "use-cases.send-eth-tx"

type sendETHTxUseCase struct {
	signTx           usecases.SignETHTransactionUseCase
	crafter          usecases.CraftTransactionUseCase
	nonceChecker     nonce.Manager
	ec               ethclient.TransactionSender
	chainRegistryURL string
	jobClient        client.JobClient
	logger           *log.Logger
}

func NewSendEthTxUseCase(signTx usecases.SignETHTransactionUseCase, crafter usecases.CraftTransactionUseCase,
	ec ethclient.TransactionSender, jobClient client.JobClient, chainRegistryURL string,
	nonceChecker nonce.Manager) usecases.SendETHTxUseCase {
	return &sendETHTxUseCase{
		jobClient:        jobClient,
		ec:               ec,
		chainRegistryURL: chainRegistryURL,
		signTx:           signTx,
		nonceChecker:     nonceChecker,
		crafter:          crafter,
		logger:           log.NewLogger().SetComponent(sendETHTxComponent),
	}
}

func (uc *sendETHTxUseCase) Execute(ctx context.Context, job *entities.Job) error {
	ctx = log.With(log.WithFields(
		ctx,
		log.Field("job", job.UUID),
		log.Field("tenant_id", job.TenantID),
		log.Field("owner_id", job.OwnerID),
		log.Field("schedule_uuid", job.ScheduleUUID),
	), uc.logger)

	logger := uc.logger.WithContext(ctx)
	logger.Debug("processing ethereum transaction job")

	err := uc.crafter.Execute(ctx, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendETHTxComponent)
	}

	// In case of job resending we don't need to sign again
	if job.InternalData.ParentJobUUID == job.UUID || job.Status == entities.StatusPending || job.Status == entities.StatusResending {
		err = utils2.UpdateJobStatus(ctx, uc.jobClient, job, entities.StatusResending, "", job.Transaction)
		if err != nil {
			return errors.FromError(err).ExtendComponent(sendETHTxComponent)
		}
	} else {
		job.Transaction.Raw, job.Transaction.Hash, err = uc.signTx.Execute(ctx, job)
		if err != nil {
			return errors.FromError(err).ExtendComponent(sendETHTxComponent)
		}

		err = utils2.UpdateJobStatus(ctx, uc.jobClient, job, entities.StatusPending, "", job.Transaction)
		if err != nil {
			return errors.FromError(err).ExtendComponent(sendETHTxComponent)
		}
	}

	txHash, err := uc.sendTx(ctx, job)
	if err != nil {
		if err2 := uc.nonceChecker.CleanNonce(ctx, job, err); err2 != nil {
			return errors.FromError(err2).ExtendComponent(sendETHTxComponent)
		}
		return errors.FromError(err).ExtendComponent(sendETHTxComponent)
	}

	if err = uc.nonceChecker.IncrementNonce(ctx, job); err != nil {
		return errors.FromError(err).ExtendComponent(sendETHTxComponent)
	}

	if txHash.String() != job.Transaction.Hash.String() {
		warnMessage := fmt.Sprintf("expected transaction hash %s, but got %s. overriding", job.Transaction.Hash, txHash)
		job.Transaction.Hash = txHash
		err = utils2.UpdateJobStatus(ctx, uc.jobClient, job, entities.StatusWarning, warnMessage, job.Transaction)
		if err != nil {
			return errors.FromError(err).ExtendComponent(sendETHTxComponent)
		}
	}

	logger.Info("ethereum transaction job was sent successfully")
	return nil
}

func (uc *sendETHTxUseCase) sendTx(ctx context.Context, job *entities.Job) (*ethcommon.Hash, error) {
	proxyURL := utils.GetProxyURL(uc.chainRegistryURL, job.ChainUUID)
	txHash, err := uc.ec.SendRawTransaction(ctx, proxyURL, job.Transaction.Raw)
	if err != nil {
		uc.logger.WithContext(ctx).WithError(err).Error("cannot send raw ethereum transaction")
		return nil, err
	}

	return &txHash, nil
}
