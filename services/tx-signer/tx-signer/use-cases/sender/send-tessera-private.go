package sender

import (
	"context"
	"fmt"

	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"

	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/use-cases"
)

const sendTesseraPrivateTxComponent = "use-cases.send-tessera-private-tx"

type sendTesseraPrivateTxUseCase struct {
	ec               ethclient.QuorumTransactionSender
	chainRegistryURL string
	client           orchestrateclient.OrchestrateClient
}

func NewSendTesseraPrivateTxUseCase(
	ec ethclient.QuorumTransactionSender,
	client orchestrateclient.OrchestrateClient,
	chainRegistryURL string,
) usecases.SendTesseraPrivateTxUseCase {
	return &sendTesseraPrivateTxUseCase{
		ec:               ec,
		chainRegistryURL: chainRegistryURL,
		client:           client,
	}
}

func (uc *sendTesseraPrivateTxUseCase) Execute(ctx context.Context, job *entities.Job) error {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("processing tessera private transaction job")

	var err error
	job.Transaction.EnclaveKey, err = uc.sendTx(ctx, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendTesseraPrivateTxComponent)
	}

	_, err = uc.client.UpdateJob(ctx, job.UUID, &types.UpdateJobRequest{
		Transaction: job.Transaction,
		Status:      utils.StatusStored,
	})
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendTesseraPrivateTxComponent)
	}

	logger.Info("Tessera private job was sent successfully")
	return nil
}

func (uc *sendTesseraPrivateTxUseCase) sendTx(ctx context.Context, job *entities.Job) (string, error) {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("sending Tessera private transaction")

	proxyTessera := fmt.Sprintf("%s/tessera/%s", uc.chainRegistryURL, job.ChainUUID)
	data, err := hexutil.Decode(job.Transaction.Data)
	if err != nil {
		errMsg := "cannot decode transaction data"
		logger.WithError(err).Errorf(errMsg)
		return "", err
	}

	enclaveKey, err := uc.ec.StoreRaw(ctx, proxyTessera, data, job.Transaction.PrivateFrom)
	if err != nil {
		errMsg := "cannot send tessera private transaction"
		logger.WithError(err).Errorf(errMsg)
		return "", err
	}

	return enclaveKey, nil
}
