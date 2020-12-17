package sender

import (
	"context"
	"fmt"

	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/txscheduler"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/use-cases"
)

const sendETHRawTxComponent = "use-cases.send-eth-raw-tx"

type sendETHRawTxUseCase struct {
	ec               ethclient.TransactionSender
	chainRegistryURL string
	client           orchestrateclient.OrchestrateClient
}

func NewSendETHRawTxUseCase(ec ethclient.TransactionSender, client orchestrateclient.OrchestrateClient,
	chainRegistryURL string) usecases.SendETHRawTxUseCase {
	return &sendETHRawTxUseCase{
		client:           client,
		chainRegistryURL: chainRegistryURL,
		ec:               ec,
	}
}

// Execute signs a public Ethereum transaction
func (uc *sendETHRawTxUseCase) Execute(ctx context.Context, job *entities.Job) error {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("processing ethereum raw transaction job")

	var err error
	job.Transaction, err = uc.rawTxDecoder(job.Transaction.Raw)
	if err != nil {
		logger.WithError(err).Error("failed to decode transaction")
		return errors.FromError(err).ExtendComponent(sendETHRawTxComponent)
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
		return errors.FromError(err).ExtendComponent(sendETHRawTxComponent)
	}

	txHash, err := uc.sendTx(ctx, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendETHRawTxComponent)
	}

	if txHash != job.Transaction.Hash {
		job.Transaction.Hash = txHash
		_, err = uc.client.UpdateJob(ctx, job.UUID, &txschedulertypes.UpdateJobRequest{
			Message:     fmt.Sprintf("expected transaction hash %s, but got %s. Overriding", job.Transaction.Hash, txHash),
			Transaction: job.Transaction,
			Status:      utils.StatusWarning,
		})
		if err != nil {
			return errors.FromError(err).ExtendComponent(sendETHRawTxComponent)
		}
	}

	logger.Info("ethereum raw transaction job was processed successfully")
	return nil
}

func (uc *sendETHRawTxUseCase) sendTx(ctx context.Context, job *entities.Job) (string, error) {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("sending ethereum raw transaction")

	proxyURL := fmt.Sprintf("%s/%s", uc.chainRegistryURL, job.ChainUUID)
	txHash, err := uc.ec.SendRawTransaction(ctx, proxyURL, job.Transaction.Raw)
	if err != nil {
		errMsg := "cannot send ethereum raw transaction"
		logger.WithError(err).Errorf(errMsg)
		return "", err
	}

	return txHash.String(), nil
}

func (uc *sendETHRawTxUseCase) rawTxDecoder(raw string) (*entities.ETHTransaction, error) {
	var tx *types.Transaction

	rawb, err := hexutil.Decode(raw)
	if err != nil {
		return nil, err
	}

	err = rlp.DecodeBytes(rawb, &tx)
	if err != nil {
		return nil, err
	}

	msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
	if err != nil {
		return nil, err
	}

	jobTx := &entities.ETHTransaction{
		From:     msg.From().String(),
		Data:     string(tx.Data()),
		Gas:      fmt.Sprintf("%d", tx.Gas()),
		GasPrice: fmt.Sprintf("%d", tx.GasPrice()),
		Value:    tx.Value().String(),
		Nonce:    fmt.Sprintf("%d", tx.Nonce()),
		Hash:     tx.Hash().String(),
		Raw:      raw,
	}

	// If not contract creation
	if tx.To() != nil {
		jobTx.To = tx.To().String()
	}

	return jobTx, nil
}
