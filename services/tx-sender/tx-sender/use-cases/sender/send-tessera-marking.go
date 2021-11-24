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
)

const sendTesseraMarkingTxComponent = "use-cases.send-tessera-marking-tx"

type sendTesseraMarkingTxUseCase struct {
	signTx           usecases.SignETHTransactionUseCase
	crafter          usecases.CraftTransactionUseCase
	nonceChecker     nonce.Manager
	jobClient        client.JobClient
	ec               ethclient.QuorumTransactionSender
	chainRegistryURL string
	logger           *log.Logger
}

func NewSendTesseraMarkingTxUseCase(signTx usecases.SignQuorumPrivateTransactionUseCase, crafter usecases.CraftTransactionUseCase,
	ec ethclient.QuorumTransactionSender, jobClient client.JobClient, chainRegistryURL string,
	nonceChecker nonce.Manager) usecases.SendTesseraMarkingTxUseCase {
	return &sendTesseraMarkingTxUseCase{
		signTx:           signTx,
		nonceChecker:     nonceChecker,
		ec:               ec,
		jobClient:        jobClient,
		chainRegistryURL: chainRegistryURL,
		crafter:          crafter,
		logger:           log.NewLogger().SetComponent(sendTesseraMarkingTxComponent),
	}
}

func (uc *sendTesseraMarkingTxUseCase) Execute(ctx context.Context, job *entities.Job) error {
	ctx = log.With(log.WithFields(
		ctx,
		log.Field("job", job.UUID),
		log.Field("tenant_id", job.TenantID),
		log.Field("owner_id", job.OwnerID),
		log.Field("schedule_uuid", job.ScheduleUUID),
	), uc.logger)
	logger := uc.logger.WithContext(ctx)
	logger.Debug("processing tessera marking transaction job")

	err := uc.crafter.Execute(ctx, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendTesseraMarkingTxComponent)
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
			return errors.FromError(err).ExtendComponent(sendTesseraMarkingTxComponent)
		}

		err = utils2.UpdateJobStatus(ctx, uc.jobClient, job, entities.StatusPending, "", job.Transaction)
		if err != nil {
			return errors.FromError(err).ExtendComponent(sendTesseraMarkingTxComponent)
		}
	}

	txHash, err := uc.sendTx(ctx, job)
	if err != nil {
		if err2 := uc.nonceChecker.CleanNonce(ctx, job, err); err2 != nil {
			return errors.FromError(err2).ExtendComponent(sendTesseraMarkingTxComponent)
		}
		return errors.FromError(err).ExtendComponent(sendTesseraMarkingTxComponent)
	}

	err = uc.nonceChecker.IncrementNonce(ctx, job)
	if err != nil {
		return err
	}

	if txHash != job.Transaction.Hash {
		warnMessage := fmt.Sprintf("expected transaction hash %s, but got %s. Overriding", job.Transaction.Hash, txHash)
		job.Transaction.Hash = txHash
		err = utils2.UpdateJobStatus(ctx, uc.jobClient, job, entities.StatusWarning, warnMessage, job.Transaction)
		if err != nil {
			return errors.FromError(err).ExtendComponent(sendTesseraMarkingTxComponent)
		}
	}

	logger.Info("tessera marking transaction job was sent successfully")
	return nil
}

func (uc *sendTesseraMarkingTxUseCase) sendTx(ctx context.Context, job *entities.Job) (string, error) {
	proxyURL := utils.GetProxyURL(uc.chainRegistryURL, job.ChainUUID)
	txHash, err := uc.ec.SendQuorumRawPrivateTransaction(ctx, proxyURL, job.Transaction.Raw, job.Transaction.PrivateFor,
		job.Transaction.MandatoryFor, int(job.Transaction.PrivacyFlag))
	if err != nil {
		uc.logger.WithContext(ctx).WithError(err).Error("cannot send tessera marking transaction")
		return "", err
	}

	return txHash.String(), nil
}
