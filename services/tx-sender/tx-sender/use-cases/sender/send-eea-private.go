package sender

import (
	"context"

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

const sendEEAPrivateTxComponent = "use-cases.send-eea-private-tx"

type sendEEAPrivateTxUseCase struct {
	crafter          usecases.CraftTransactionUseCase
	signTx           usecases.SignETHTransactionUseCase
	nonceManager     nonce.Manager
	jobClient        client.JobClient
	ec               ethclient.EEATransactionSender
	chainRegistryURL string
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
	}
}

func (uc *sendEEAPrivateTxUseCase) Execute(ctx context.Context, job *entities.Job) error {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
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
		if err2 := uc.nonceManager.DecreaseNonce(ctx, job, err); err2 != nil {
			return errors.FromError(err2).ExtendComponent(sendEEAPrivateTxComponent)
		}
		return err
	}

	err = uc.nonceManager.IncrementNonce(ctx, job)
	if err != nil {
		return err
	}

	err = utils2.UpdateJobStatus(ctx, uc.jobClient, job.UUID, utils.StatusStored, "", job.Transaction)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendEEAPrivateTxComponent)
	}

	logger.Info("EEA private transaction job was sent successfully")
	return nil
}

func (uc *sendEEAPrivateTxUseCase) sendTx(ctx context.Context, job *entities.Job) (string, error) {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)

	proxyURL := utils.GetProxyURL(uc.chainRegistryURL, job.ChainUUID)
	txHash, err := uc.ec.PrivDistributeRawTransaction(ctx, proxyURL, job.Transaction.Raw)
	if err != nil {
		errMsg := "cannot send EEA private transaction"
		logger.WithError(err).Errorf(errMsg)
		return "", err
	}

	return txHash.String(), nil
}
