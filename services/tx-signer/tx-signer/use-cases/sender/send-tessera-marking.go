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

const sendTesseraMarkingTxComponent = "use-cases.send-tessera-marking-tx"

type sendTesseraMarkingTxUseCase struct {
	signTx            usecases.SignETHTransactionUseCase
	nonceChecker      nonce.Checker
	txSchedulerClient client.TransactionSchedulerClient
	ec                ethclient.QuorumTransactionSender
	chainRegistryURL  string
}

func NewSendTesseraMarkingTxUseCase(signTx usecases.SignQuorumPrivateTransactionUseCase,
	ec ethclient.QuorumTransactionSender, txSchedulerClient client.TransactionSchedulerClient, chainRegistryURL string,
	nonceChecker nonce.Checker) usecases.SendTesseraMarkingTxUseCase {
	return &sendTesseraMarkingTxUseCase{
		signTx:            signTx,
		nonceChecker:      nonceChecker,
		ec:                ec,
		txSchedulerClient: txSchedulerClient,
		chainRegistryURL:  chainRegistryURL,
	}
}

// Execute signs a public Ethereum transaction
func (uc *sendTesseraMarkingTxUseCase) Execute(ctx context.Context, job *entities.Job) error {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("processing tessera marking transaction job")

	err := uc.nonceChecker.Check(ctx, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendTesseraMarkingTxComponent)
	}

	job.Transaction.Raw, job.Transaction.Hash, err = uc.signTx.Execute(ctx, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendTesseraMarkingTxComponent)
	}

	err = utils2.UpdateJobStatus(ctx, uc.txSchedulerClient, job.UUID, utils.StatusPending, "", job.Transaction)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendTesseraMarkingTxComponent)
	}

	txHash, err := uc.sendTx(ctx, job)
	if err != nil {
		if err2 := uc.nonceChecker.OnFailure(ctx, job, err); err2 != nil {
			return errors.FromError(err2).ExtendComponent(sendTesseraMarkingTxComponent)
		}
		return errors.FromError(err).ExtendComponent(sendTesseraMarkingTxComponent)
	}

	err = uc.nonceChecker.OnSuccess(ctx, job)
	if err != nil {
		return err
	}

	if txHash != job.Transaction.Hash {
		warnMessage := fmt.Sprintf("expected transaction hash %s, but got %s. Overriding", job.Transaction.Hash, txHash)
		job.Transaction.Hash = txHash
		err = utils2.UpdateJobStatus(ctx, uc.txSchedulerClient, job.UUID, utils.StatusWarning, warnMessage, job.Transaction)
		if err != nil {
			return errors.FromError(err).ExtendComponent(sendTesseraMarkingTxComponent)
		}
	}

	logger.Info("tessera marking job was sent successfully")
	return nil
}

func (uc *sendTesseraMarkingTxUseCase) sendTx(ctx context.Context, job *entities.Job) (string, error) {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("sending Tessera marking transaction")

	proxyTessera := fmt.Sprintf("%s/%s", uc.chainRegistryURL, job.ChainUUID)
	txHash, err := uc.ec.SendQuorumRawPrivateTransaction(ctx, proxyTessera, job.Transaction.Raw, job.Transaction.PrivateFor)
	if err != nil {
		errMsg := "cannot send tessera marking transaction"
		logger.WithError(err).Errorf(errMsg)
		return "", err
	}

	return txHash.String(), nil
}
