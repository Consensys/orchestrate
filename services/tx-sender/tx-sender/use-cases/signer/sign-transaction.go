package signer

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/encoding/rlp"
	qkm "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager"
	qkmtypes "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager/types"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/pkg/utils"

	pkgcryto "github.com/ConsenSys/orchestrate/pkg/crypto/ethereum"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/parsers"
	"github.com/ethereum/go-ethereum/crypto"

	usecases "github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/use-cases"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/quorum-key-manager/client"
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

func (uc *signETHTransactionUseCase) Execute(ctx context.Context, job *entities.Job) (raw, txHash string, err error) {
	logger := uc.logger.WithContext(ctx).WithField("one_time_key", job.InternalData.OneTimeKey)

	signer := pkgcryto.GetEIP155Signer(job.InternalData.ChainID)
	transaction := parsers.ETHTransactionToTransaction(job.Transaction)

	var decodedSignature []byte
	if job.InternalData.OneTimeKey {
		decodedSignature, err = uc.signWithOneTimeKey(ctx, transaction, signer)
	} else {
		decodedSignature, err = uc.signWithAccount(ctx, job, transaction, signer)
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

func (uc *signETHTransactionUseCase) signWithAccount(ctx context.Context, job *entities.Job, tx *types.Transaction, signer types.Signer) ([]byte, error) {
	logger := uc.logger.WithContext(ctx)

	tenants := utils.AllowedTenants(job.TenantID)
	isAllowed, err := qkm.IsTenantAllowed(ctx, uc.keyManagerClient, tenants, uc.storeName, job.Transaction.From)
	if err != nil {
		errMsg := "failed to to sign transaction, cannot fetch account"
		uc.logger.WithField("address", job.Transaction.From).WithError(err).Error(errMsg)
		return nil, errors.DependencyFailureError(errMsg).AppendReason(err.Error())
	}

	if !isAllowed {
		errMessage := "failed to to sign transaction, tenant is not allowed"
		logger.WithField("address", job.Transaction.From).WithField("tenants", tenants).Error(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	txData := signer.Hash(tx).Bytes()
	sig, err := uc.keyManagerClient.SignEth1Data(ctx, uc.storeName, job.Transaction.From, &qkmtypes.SignHexPayloadRequest{
		Data: txData,
	})
	if err != nil {
		errMsg := "failed to sign ethereum transaction using key manager"
		logger.Error(errMsg)
		return nil, errors.DependencyFailureError(errMsg).AppendReason(err.Error())
	}

	decodedSignature, err := hexutil.Decode(sig)
	if err != nil {
		errMessage := "failed to decode signature"
		logger.WithError(err).Error(errMessage)
		return nil, errors.EncodingError(errMessage)
	}

	return decodedSignature, nil
}
