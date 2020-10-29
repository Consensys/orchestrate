package ethereum

import (
	"context"
	"math/big"
	"strconv"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer-new/tx-signer/parsers"

	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer-new/tx-signer/use-cases"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"
)

const signTransactionComponent = "use-cases.sign-transaction"

// signTransactionUseCase is a use case to sign a public Ethereum transaction
type signTransactionUseCase struct {
	keyManagerClient client.KeyManagerClient
}

// NewSignTransactionUseCase creates a new SignTransactionUseCase
func NewSignTransactionUseCase(keyManagerClient client.KeyManagerClient) usecases.SignTransactionUseCase {
	return &signTransactionUseCase{
		keyManagerClient: keyManagerClient,
	}
}

// Execute signs a public Ethereum transaction
func (uc *signTransactionUseCase) Execute(ctx context.Context, job *entities.Job) (raw, txHash string, err error) {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("signing ethereum transaction")

	if job.InternalData.OneTimeKey {
		// TODO: Sign one time key tx
		return "", "", errors.InternalError("Not implemented yet")
	}

	// TODO: Temporary before types alignment
	nonce, _ := strconv.ParseUint(job.Transaction.Nonce, 10, 64)
	gasLimit, _ := strconv.ParseUint(job.Transaction.Gas, 10, 64)
	sig, err := uc.keyManagerClient.ETHSignTransaction(ctx, job.Transaction.From, &ethereum.SignETHTransactionRequest{
		Namespace: job.TenantID,
		Nonce:     nonce,
		Amount:    job.Transaction.Value,
		GasPrice:  job.Transaction.GasPrice,
		GasLimit:  gasLimit,
		Data:      job.Transaction.Data,
		To:        job.Transaction.To,
		ChainID:   job.InternalData.ChainID,
	})
	if err != nil {
		logger.WithError(err).Error("failed to sign ethereum transaction using key manager")
		return "", "", errors.FromError(err).ExtendComponent(signTransactionComponent)
	}

	decodedSignature, err := hexutil.Decode(sig)
	if err != nil {
		errMessage := "failed to decode signature"
		logger.WithField("encoded_signature", sig).WithError(err).Error(errMessage)
		return "", "", errors.EncodingError(errMessage).ExtendComponent(signTransactionComponent)
	}

	transaction := parsers.ETHTransactionToTransaction(job.Transaction)
	chainID := new(big.Int)
	chainID, _ = chainID.SetString(job.InternalData.ChainID, 10)
	signedTx, err := transaction.WithSignature(types.NewEIP155Signer(chainID), decodedSignature)
	if err != nil {
		errMessage := "failed to set transaction signature"
		logger.WithError(err).Error(errMessage)
		return "", "", errors.InvalidParameterError(errMessage).ExtendComponent(signTransactionComponent)
	}

	signedRaw, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		errMessage := "failed to RLP encode signed transaction"
		logger.WithError(err).Error(errMessage)
		return "", "", errors.CryptoOperationError(errMessage).ExtendComponent(signTransactionComponent)
	}
	txHash = transaction.Hash().Hex()

	logger.WithField("txHash", txHash).Info("ethereum transaction signed successfully")
	return hexutil.Encode(signedRaw), txHash, nil
}
