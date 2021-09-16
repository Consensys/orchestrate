package signer

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/encoding/rlp"
	qkm "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/pkg/utils"
	qkmtypes "github.com/consensys/quorum-key-manager/src/stores/api/types"

	pkgcryto "github.com/ConsenSys/orchestrate/pkg/crypto/ethereum"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/parsers"
	"github.com/ethereum/go-ethereum/crypto"

	usecases "github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/use-cases"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

const signTransactionComponent = "use-cases.sign-eth-transaction"

// signETHTransactionUseCase is a use case to sign a public Ethereum transaction
type signETHTransactionUseCase struct {
	keyManagerClient client.KeyManagerClient
	logger           *log.Logger
	storeName        string
}

// NewSignETHTransactionUseCase creates a new SignTransactionUseCase
func NewSignETHTransactionUseCase(keyManagerClient client.KeyManagerClient) usecases.SignETHTransactionUseCase {
	return &signETHTransactionUseCase{
		keyManagerClient: keyManagerClient,
		logger:           log.NewLogger().SetComponent(signTransactionComponent),
		storeName:        qkm.GlobalStoreName(),
	}
}

func (uc *signETHTransactionUseCase) Execute(ctx context.Context, job *entities.Job) (signedRaw, txHash string, err error) {
	logger := uc.logger.WithContext(ctx).WithField("one_time_key", job.InternalData.OneTimeKey)

	transaction := parsers.ETHTransactionToTransaction(job.Transaction)
	if job.InternalData.OneTimeKey {
		signedRaw, txHash, err = uc.signWithOneTimeKey(ctx, transaction, job.InternalData.ChainID)
	} else {
		signedRaw, txHash, err = uc.signWithAccount(ctx, job, transaction, job.InternalData.ChainID)
	}
	if err != nil {
		return "", "", errors.FromError(err).ExtendComponent(signTransactionComponent)
	}

	logger.WithField("tx_hash", txHash).Debug("ethereum transaction signed successfully")
	return signedRaw, txHash, nil
}

func (uc *signETHTransactionUseCase) signWithOneTimeKey(ctx context.Context, transaction *types.Transaction, chainID string) (signedRaw, txHash string, err error) {
	logger := uc.logger.WithContext(ctx)
	privKey, err := crypto.GenerateKey()
	if err != nil {
		errMessage := "failed to generate Ethereum private one time key"
		logger.WithError(err).Error(errMessage)
		return "", "", errors.CryptoOperationError(errMessage)
	}

	signer := pkgcryto.GetEIP155Signer(chainID)
	decodedSignature, err := pkgcryto.SignTransaction(transaction, privKey, signer)
	if err != nil {
		logger.WithError(err).Error("failed to sign Ethereum transaction")
		return "", "", err
	}

	signedTx, err := transaction.WithSignature(signer, decodedSignature)
	if err != nil {
		errMessage := "failed to set transaction signature"
		logger.WithError(err).Error(errMessage)
		return "", "", errors.InvalidParameterError(errMessage).ExtendComponent(signTransactionComponent)
	}

	signedRawB, err := rlp.Encode(signedTx)
	if err != nil {
		errMessage := "failed to RLP encode signed transaction"
		logger.WithError(err).Error(errMessage)
		return "", "", errors.CryptoOperationError(errMessage).ExtendComponent(signTransactionComponent)
	}

	return hexutil.Encode(signedRawB), signedTx.Hash().Hex(), nil
}

func (uc *signETHTransactionUseCase) signWithAccount(ctx context.Context, job *entities.Job, tx *types.Transaction, chainID string) (signedRaw, txHash string, err error) {
	logger := uc.logger.WithContext(ctx)

	tenants := utils.AllowedTenants(job.TenantID)
	isAllowed, err := qkm.IsTenantAllowed(ctx, uc.keyManagerClient, tenants, uc.storeName, job.Transaction.From)
	if err != nil {
		errMsg := "failed to to sign transaction, cannot fetch account"
		uc.logger.WithField("address", job.Transaction.From).WithError(err).Error(errMsg)
		return "", "", errors.DependencyFailureError(errMsg).AppendReason(err.Error())
	}

	if !isAllowed {
		errMessage := "failed to to sign transaction, tenant is not allowed"
		logger.WithField("address", job.Transaction.From).WithField("tenants", tenants).Error(errMessage)
		return "", "", errors.UnauthorizedError(errMessage)
	}

	signedRaw, err = uc.keyManagerClient.SignTransaction(ctx, uc.storeName, job.Transaction.From, &qkmtypes.SignETHTransactionRequest{
		Nonce:           hexutil.Uint64(tx.Nonce()),
		To:              tx.To(),
		Data:            tx.Data(),
		ChainID:         hexutil.Big(*utils.MustEncodeBigInt(chainID)),
		Value:           hexutil.Big(*tx.Value()),
		GasPrice:        hexutil.Big(*tx.GasPrice()),
		GasLimit:        hexutil.Uint64(tx.Gas()),
		TransactionType: qkmtypes.LegacyTxType,
	})
	if err != nil {
		errMsg := "failed to sign ethereum transaction using key manager"
		logger.WithError(err).Error(errMsg)
		return "", "", errors.DependencyFailureError(errMsg).AppendReason(err.Error())
	}

	err = rlp.Decode(hexutil.MustDecode(signedRaw), tx)
	if err != nil {
		errMessage := "failed to decode signature"
		logger.WithError(err).Error(errMessage)
		return "", "", errors.EncodingError(errMessage)
	}

	return signedRaw, tx.Hash().Hex(), nil
}
