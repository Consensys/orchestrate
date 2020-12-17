package sender

import (
	"context"
	"fmt"

	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/txscheduler"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/nonce"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/use-cases"
)

const sendETHTxComponent = "use-cases.send-eth-tx"

type sendETHTxUseCase struct {
	signTx           usecases.SignETHTransactionUseCase
	nonceChecker     nonce.Checker
	ec               ethclient.TransactionSender
	chainRegistryURL string
	client           orchestrateclient.OrchestrateClient
}

func NewSendEthTxUseCase(
	signTx usecases.SignETHTransactionUseCase,
	ec ethclient.TransactionSender,
	client orchestrateclient.OrchestrateClient,
	chainRegistryURL string,
	nonceChecker nonce.Checker,
) usecases.SendETHTxUseCase {
	return &sendETHTxUseCase{
		client:           client,
		ec:               ec,
		chainRegistryURL: chainRegistryURL,
		signTx:           signTx,
		nonceChecker:     nonceChecker,
	}
}

// Execute signs a public Ethereum transaction
func (uc *sendETHTxUseCase) Execute(ctx context.Context, job *entities.Job) error {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("processing ethereum transaction job")

	err := uc.nonceChecker.Check(ctx, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendETHTxComponent)
	}

	job.Transaction.Raw, job.Transaction.Hash, err = uc.signTx.Execute(ctx, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendETHTxComponent)
	}

	txUpdateReq := &txschedulertypes.UpdateJobRequest{
		Transaction: job.Transaction,
	}
	if job.InternalData.ParentJobUUID == job.UUID {
		txUpdateReq.Status = utils.StatusResending
	} else {
		txUpdateReq.Status = utils.StatusPending
	}

	_, err = uc.client.UpdateJob(ctx, job.UUID, txUpdateReq)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendETHTxComponent)
	}

	txHash, err := uc.sendTx(ctx, job)
	if err != nil {
		if err2 := uc.nonceChecker.OnFailure(ctx, job, err); err2 != nil {
			return errors.FromError(err2).ExtendComponent(sendETHTxComponent)
		}
		return errors.FromError(err).ExtendComponent(sendETHTxComponent)
	}

	if err = uc.nonceChecker.OnSuccess(ctx, job); err != nil {
		return errors.FromError(err).ExtendComponent(sendETHTxComponent)
	}

	if txHash != job.Transaction.Hash {
		job.Transaction.Hash = txHash
		_, err = uc.client.UpdateJob(ctx, job.UUID, &txschedulertypes.UpdateJobRequest{
			Message:     fmt.Sprintf("expected transaction hash %s, but got %s. Overriding", job.Transaction.Hash, txHash),
			Transaction: job.Transaction,
			Status:      utils.StatusWarning,
		})
		if err != nil {
			return errors.FromError(err).ExtendComponent(sendETHTxComponent)
		}
	}

	logger.Info("ethereum transaction job was sent successfully")
	return nil
}

func (uc *sendETHTxUseCase) sendTx(ctx context.Context, job *entities.Job) (string, error) {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("sending ethereum transaction")

	proxyURL := fmt.Sprintf("%s/%s", uc.chainRegistryURL, job.ChainUUID)
	txHash, err := uc.ec.SendRawTransaction(ctx, proxyURL, job.Transaction.Raw)
	if err != nil {
		errMsg := "cannot send ethereum transaction"
		logger.WithError(err).Errorf(errMsg)
		return "", err
	}

	return txHash.String(), nil
}
