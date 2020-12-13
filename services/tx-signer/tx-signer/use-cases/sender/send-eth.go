package sender

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/nonce"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/use-cases"
	utils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/utils"
)

const sendETHTxComponent = "use-cases.send-eth-tx"

type sendETHTxUseCase struct {
	signTx            usecases.SignETHTransactionUseCase
	nonceChecker      nonce.Checker
	ec                ethclient.TransactionSender
	chainRegistryURL  string
	txSchedulerClient client.TransactionSchedulerClient
}

func NewSendEthTxUseCase(signTx usecases.SignETHTransactionUseCase, ec ethclient.TransactionSender,
	txSchedulerClient client.TransactionSchedulerClient, chainRegistryURL string, nonceChecker nonce.Checker,
) usecases.SendETHTxUseCase {
	return &sendETHTxUseCase{
		txSchedulerClient: txSchedulerClient,
		ec:                ec,
		chainRegistryURL:  chainRegistryURL,
		signTx:            signTx,
		nonceChecker:      nonceChecker,
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

	if job.InternalData.ParentJobUUID == job.UUID {
		err = utils2.UpdateJobStatus(ctx, uc.txSchedulerClient, job.UUID, utils.StatusResending, "", job.Transaction)
	} else {
		err = utils2.UpdateJobStatus(ctx, uc.txSchedulerClient, job.UUID, utils.StatusPending, "", job.Transaction)
	}

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
		warnMessage := fmt.Sprintf("expected transaction hash %s, but got %s. Overriding", job.Transaction.Hash, txHash)
		job.Transaction.Hash = txHash
		err = utils2.UpdateJobStatus(ctx, uc.txSchedulerClient, job.UUID, utils.StatusWarning, warnMessage, job.Transaction)
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
