package signer

import (
	"context"
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/encoding/rlp"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/pkg/utils"

	pkgcryto "github.com/ConsenSys/orchestrate/pkg/crypto/ethereum"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/parsers"
	"github.com/ethereum/go-ethereum/crypto"

	usecases "github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/use-cases"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/types/keymanager/ethereum"
	"github.com/ConsenSys/orchestrate/services/key-manager/client"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

const signTransactionComponent = "use-cases.sign-eth-transaction"

// signETHTransactionUseCase is a use case to sign a public Ethereum transaction
type signETHTransactionUseCase struct {
	keyManagerClient client.KeyManagerClient
	logger           *log.Logger
}

// NewSignETHTransactionUseCase creates a new SignTransactionUseCase
func NewSignETHTransactionUseCase(keyManagerClient client.KeyManagerClient) usecases.SignETHTransactionUseCase {
	return &signETHTransactionUseCase{
		keyManagerClient: keyManagerClient,
		logger:           log.NewLogger().SetComponent(signTransactionComponent),
	}
}

func (uc *signETHTransactionUseCase) Execute(ctx context.Context, job *entities.Job) (raw, txHash string, err error) {
	logger := uc.logger.WithContext(ctx).WithField("one_time_key", job.InternalData.OneTimeKey)

	signer := pkgcryto.GetEIP155Signer(job.InternalData.ChainID)
	transaction := parsers.ETHTransactionToTransaction(job.Transaction)

	var decodedSignature []byte
	if job.InternalData.OneTimeKey {
		decodedSignature, err = uc.signWithOneTimeKey(ctx, transaction, signer)
	} else {
		decodedSignature, err = uc.signWithAccount(ctx, job, transaction)
	}
	if err != nil {
		return "", "", errors.FromError(err).ExtendComponent(signTransactionComponent)
	}

	signedTx, err := transaction.WithSignature(signer, decodedSignature)
	if err != nil {
		errMessage := "failed to set transaction signature"
		logger.WithError(err).Error(errMessage)
		return "", "", errors.InvalidParameterError(errMessage).ExtendComponent(signTransactionComponent)
	}

	signedRaw, err := rlp.Encode(signedTx)
	if err != nil {
		errMessage := "failed to RLP encode signed transaction"
		logger.WithError(err).Error(errMessage)
		return "", "", errors.CryptoOperationError(errMessage).ExtendComponent(signTransactionComponent)
	}
	txHash = signedTx.Hash().Hex()

	logger.WithField("tx_hash", txHash).Debug("ethereum transaction signed successfully")
	return hexutil.Encode(signedRaw), txHash, nil
}

func (uc *signETHTransactionUseCase) signWithOneTimeKey(ctx context.Context, transaction *types.Transaction,
	signer types.Signer) ([]byte, error) {
	logger := uc.logger.WithContext(ctx)
	privKey, err := crypto.GenerateKey()
	if err != nil {
		errMessage := "failed to generate Ethereum private one time key"
		logger.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	sign, err := pkgcryto.SignTransaction(transaction, privKey, signer)
	if err != nil {
		logger.WithError(err).Error("failed to sign Ethereum transaction")
		return nil, err
	}

	return sign, nil
}

func (uc *signETHTransactionUseCase) signWithAccount(ctx context.Context, job *entities.Job, tx *types.Transaction) ([]byte, error) {
	logger := uc.logger.WithContext(ctx)
	request := &ethereum.SignETHTransactionRequest{
		Nonce:    tx.Nonce(),
		Amount:   tx.Value().String(),
		GasPrice: tx.GasPrice().String(),
		GasLimit: tx.Gas(),
		Data:     hexutil.Encode(tx.Data()),
		ChainID:  job.InternalData.ChainID,
	}
	if tx.To() != nil {
		request.To = tx.To().Hex()
	}

	tenants := utils.AllowedTenants(job.TenantID)
	for _, tenant := range tenants {
		request.Namespace = tenant
		sig, err := uc.keyManagerClient.ETHSignTransaction(ctx, job.Transaction.From, request)
		if err != nil && errors.IsNotFoundError(err) {
			continue
		}
		if err != nil {
			logger.WithError(err).Error("failed to sign ethereum transaction using key manager")
			return nil, errors.FromError(err)
		}

		decodedSignature, err := hexutil.Decode(sig)
		if err != nil {
			errMessage := "failed to decode signature"
			logger.WithError(err).Error(errMessage)
			return nil, errors.EncodingError(errMessage)
		}

		return decodedSignature, nil
	}

	errMessage := fmt.Sprintf("account %s was not found on key-manager", job.Transaction.From)
	logger.WithField("from_account", job.Transaction.From).WithField("tenants", tenants).Error(errMessage)
	return nil, errors.InvalidParameterError(errMessage)
}
