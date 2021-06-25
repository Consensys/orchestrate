package sender

import (
	"context"
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/sdk/client"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/ethclient"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/pkg/utils"
	usecases "github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/use-cases"
	utils2 "github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

const sendETHRawTxComponent = "use-cases.send-eth-raw-tx"

type sendETHRawTxUseCase struct {
	ec               ethclient.TransactionSender
	chainRegistryURL string
	jobClient        client.JobClient
	logger           *log.Logger
}

func NewSendETHRawTxUseCase(ec ethclient.TransactionSender, jobClient client.JobClient,
	chainRegistryURL string) usecases.SendETHRawTxUseCase {
	return &sendETHRawTxUseCase{
		jobClient:        jobClient,
		chainRegistryURL: chainRegistryURL,
		ec:               ec,
		logger:           log.NewLogger().SetComponent(sendETHRawTxComponent),
	}
}

func (uc *sendETHRawTxUseCase) Execute(ctx context.Context, job *entities.Job) error {
	ctx = log.With(log.WithFields(
		ctx,
		log.Field("job", job.UUID),
		log.Field("tenant_id", job.TenantID),
		log.Field("schedule_uuid", job.ScheduleUUID),
	), uc.logger)
	logger := uc.logger.WithContext(ctx)
	logger.Debug("processing ethereum raw transaction job")

	var err error
	job.Transaction, err = uc.rawTxDecoder(job.Transaction.Raw)
	if err != nil {
		logger.WithError(err).Error("failed to decode transaction")
		return errors.FromError(err).ExtendComponent(sendETHRawTxComponent)
	}

	if job.InternalData.ParentJobUUID == job.UUID || job.Status == entities.StatusPending || job.Status == entities.StatusResending {
		err = utils2.UpdateJobStatus(ctx, uc.jobClient, job, entities.StatusResending, "", job.Transaction)
	} else {
		err = utils2.UpdateJobStatus(ctx, uc.jobClient, job, entities.StatusPending, "", job.Transaction)
	}

	if err != nil {
		return errors.FromError(err).ExtendComponent(sendETHRawTxComponent)
	}

	txHash, err := uc.sendTx(ctx, job)
	if err != nil {
		return errors.FromError(err).ExtendComponent(sendETHRawTxComponent)
	}

	if txHash != job.Transaction.Hash {
		warnMessage := fmt.Sprintf("expected transaction hash %s, but got %s. Overriding", job.Transaction.Hash, txHash)
		job.Transaction.Hash = txHash
		err = utils2.UpdateJobStatus(ctx, uc.jobClient, job, entities.StatusWarning, warnMessage, job.Transaction)
		if err != nil {
			return errors.FromError(err).ExtendComponent(sendETHRawTxComponent)
		}
	}

	logger.Info("ethereum raw transaction job was processed successfully")
	return nil
}

func (uc *sendETHRawTxUseCase) sendTx(ctx context.Context, job *entities.Job) (string, error) {
	proxyURL := utils.GetProxyURL(uc.chainRegistryURL, job.ChainUUID)
	txHash, err := uc.ec.SendRawTransaction(ctx, proxyURL, job.Transaction.Raw)
	if err != nil {
		errMsg := "cannot send ethereum raw transaction"
		uc.logger.WithContext(ctx).WithError(err).Error(errMsg)
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
		Data:     hexutil.Encode(tx.Data()),
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
