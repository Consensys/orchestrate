package sender

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/tx-sender/nonce"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/tx-sender/use-cases"
	utils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/tx-sender/utils"
)

const sendTesseraMarkingTxComponent = "use-cases.send-tessera-marking-tx"

type sendTesseraMarkingTxUseCase struct {
	signTx           usecases.SignETHTransactionUseCase
	crafter          usecases.CraftTransactionUseCase
	nonceChecker     nonce.Manager
	jobClient        client.JobClient
	ec               ethclient.QuorumTransactionSender
	chainRegistryURL string
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
	}
}

func (uc *sendTesseraMarkingTxUseCase) Execute(ctx context.Context, job *entities.Job) error {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("processing tessera marking transaction job")

	err := uc.crafter.Execute(ctx, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendTesseraMarkingTxComponent)
	}

	// In case of job resending we don't need to sign again
	if job.InternalData.ParentJobUUID == job.UUID {
		err = utils2.UpdateJobStatus(ctx, uc.jobClient, job.UUID, utils.StatusResending, "", job.Transaction)
		if err != nil {
			return errors.FromError(err).ExtendComponent(sendETHTxComponent)
		}
	} else {
		job.Transaction.Raw, job.Transaction.Hash, err = uc.signTx.Execute(ctx, job)
		if err != nil {
			return errors.FromError(err).ExtendComponent(sendTesseraMarkingTxComponent)
		}

		err = utils2.UpdateJobStatus(ctx, uc.jobClient, job.UUID, utils.StatusPending, "", job.Transaction)
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
		err = utils2.UpdateJobStatus(ctx, uc.jobClient, job.UUID, utils.StatusWarning, warnMessage, job.Transaction)
		if err != nil {
			return errors.FromError(err).ExtendComponent(sendTesseraMarkingTxComponent)
		}
	}

	logger.Info("tessera marking transaction job was sent successfully")
	return nil
}

func (uc *sendTesseraMarkingTxUseCase) sendTx(ctx context.Context, job *entities.Job) (string, error) {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("sending Tessera marking transaction")

	proxyURL := utils.GetProxyURL(uc.chainRegistryURL, job.ChainUUID)
	txHash, err := uc.ec.SendQuorumRawPrivateTransaction(ctx, proxyURL, job.Transaction.Raw, job.Transaction.PrivateFor)
	if err != nil {
		errMsg := "cannot send tessera marking transaction"
		logger.WithError(err).Errorf(errMsg)
		return "", err
	}

	return txHash.String(), nil
}
