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

const sendEEAPrivateTxComponent = "use-cases.send-eea-private-tx"

type sendEEAPrivateTxUseCase struct {
	signTx           usecases.SignETHTransactionUseCase
	nonceChecker     nonce.Checker
	client           orchestrateclient.OrchestrateClient
	ec               ethclient.EEATransactionSender
	chainRegistryURL string
}

func NewSendEEAPrivateTxUseCase(
	signTx usecases.SignEEATransactionUseCase,
	ec ethclient.EEATransactionSender,
	client orchestrateclient.OrchestrateClient,
	chainRegistryURL string,
	nonceChecker nonce.Checker,
) usecases.SendEEAPrivateTxUseCase {
	return &sendEEAPrivateTxUseCase{
		client:           client,
		chainRegistryURL: chainRegistryURL,
		signTx:           signTx,
		ec:               ec,
		nonceChecker:     nonceChecker,
	}
}

// Execute signs a public Ethereum transaction
func (uc *sendEEAPrivateTxUseCase) Execute(ctx context.Context, job *entities.Job) error {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("processing EEA private transaction job")

	err := uc.nonceChecker.Check(ctx, job)
	if err != nil {
		return err
	}

	job.Transaction.Raw, _, err = uc.signTx.Execute(ctx, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendEEAPrivateTxComponent)
	}

	job.Transaction.Hash, err = uc.sendTx(ctx, job)
	if err != nil {
		if err2 := uc.nonceChecker.OnFailure(ctx, job, err); err2 != nil {
			return errors.FromError(err2).ExtendComponent(sendEEAPrivateTxComponent)
		}
		return err
	}

	err = uc.nonceChecker.OnSuccess(ctx, job)
	if err != nil {
		return err
	}

	_, err = uc.client.UpdateJob(ctx, job.UUID, &txschedulertypes.UpdateJobRequest{
		Transaction: job.Transaction,
		Status:      utils.StatusStored,
	})
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendEEAPrivateTxComponent)
	}

	logger.Info("EEA private transaction job was sent successfully")
	return nil
}

// Execute signs a public Ethereum transaction
func (uc *sendEEAPrivateTxUseCase) sendTx(ctx context.Context, job *entities.Job) (string, error) {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)

	proxyURL := fmt.Sprintf("%s/%s", uc.chainRegistryURL, job.ChainUUID)
	txHash, err := uc.ec.PrivDistributeRawTransaction(ctx, proxyURL, job.Transaction.Raw)
	if err != nil {
		errMsg := "cannot send EEA private transaction"
		logger.WithError(err).Errorf(errMsg)
		return "", err
	}

	return txHash.String(), nil
}
