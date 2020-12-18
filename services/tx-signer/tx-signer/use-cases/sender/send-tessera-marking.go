package sender

import (
	"context"
	"fmt"

	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/nonce"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/use-cases"
)

const sendTesseraMarkingTxComponent = "use-cases.send-tessera-marking-tx"

type sendTesseraMarkingTxUseCase struct {
	signTx           usecases.SignETHTransactionUseCase
	nonceChecker     nonce.Checker
	client           orchestrateclient.OrchestrateClient
	ec               ethclient.QuorumTransactionSender
	chainRegistryURL string
}

func NewSendTesseraMarkingTxUseCase(
	signTx usecases.SignQuorumPrivateTransactionUseCase,
	ec ethclient.QuorumTransactionSender,
	client orchestrateclient.OrchestrateClient,
	chainRegistryURL string,
	nonceChecker nonce.Checker,
) usecases.SendTesseraMarkingTxUseCase {
	return &sendTesseraMarkingTxUseCase{
		signTx:           signTx,
		nonceChecker:     nonceChecker,
		ec:               ec,
		client:           client,
		chainRegistryURL: chainRegistryURL,
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

	_, err = uc.client.UpdateJob(ctx, job.UUID, &types.UpdateJobRequest{
		Transaction: job.Transaction,
		Status:      utils.StatusPending,
	})
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
		job.Transaction.Hash = txHash
		_, err = uc.client.UpdateJob(ctx, job.UUID, &types.UpdateJobRequest{
			Message:     fmt.Sprintf("expected transaction hash %s, but got %s. Overriding", job.Transaction.Hash, txHash),
			Transaction: job.Transaction,
			Status:      utils.StatusWarning,
		})
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
