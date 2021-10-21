package signer

import (
	"context"

	"github.com/consensys/orchestrate/pkg/encoding/rlp"
	qkm "github.com/consensys/orchestrate/pkg/quorum-key-manager"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	qkmtypes "github.com/consensys/quorum-key-manager/src/stores/api/types"

	pkgcryto "github.com/consensys/orchestrate/pkg/crypto/ethereum"

	"github.com/consensys/orchestrate/pkg/utils"

	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/tx-sender/tx-sender/parsers"
	"github.com/ethereum/go-ethereum/crypto"

	usecases "github.com/consensys/orchestrate/services/tx-sender/tx-sender/use-cases"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

const signEEATransactionComponent = "use-cases.sign-eea-transaction"

// signEEATransactionUseCase is a use case to sign a public Ethereum transaction
type signEEATransactionUseCase struct {
	keyManagerClient client.KeyManagerClient
	logger           *log.Logger
	storeName        string
}

// NewSignEEATransactionUseCase creates a new SignEEATransactionUseCase
func NewSignEEATransactionUseCase(keyManagerClient client.KeyManagerClient) usecases.SignEEATransactionUseCase {
	return &signEEATransactionUseCase{
		keyManagerClient: keyManagerClient,
		logger:           log.NewLogger().SetComponent(signEEATransactionComponent),
		storeName:        qkm.GlobalStoreName(),
	}
}

func (uc *signEEATransactionUseCase) Execute(ctx context.Context, job *entities.Job) (signedRaw, txHash string, err error) {
	logger := uc.logger.WithContext(ctx).WithField("one_time_key", job.InternalData.OneTimeKey)

	transaction := parsers.ETHTransactionToTransaction(job.Transaction, job.InternalData.ChainID)
	privateArgs := &entities.PrivateETHTransactionParams{
		PrivateFrom:    job.Transaction.PrivateFrom,
		PrivateFor:     job.Transaction.PrivateFor,
		PrivacyGroupID: job.Transaction.PrivacyGroupID,
		PrivateTxType:  entities.PrivateTxTypeRestricted,
	}

	if job.InternalData.OneTimeKey {
		signedRaw, err = uc.signWithOneTimeKey(ctx, transaction, privateArgs, job.InternalData.ChainID)
	} else {
		signedRaw, err = uc.signWithAccount(ctx, job, privateArgs, transaction, job.InternalData.ChainID)
	}
	if err != nil {
		return "", "", errors.FromError(err).ExtendComponent(signEEATransactionComponent)
	}

	logger.Debug("eea transaction signed successfully")

	// transaction hash of EEA transactions cannot be computed
	return signedRaw, "", nil
}

func (uc *signEEATransactionUseCase) signWithOneTimeKey(ctx context.Context, transaction *types.Transaction,
	privateArgs *entities.PrivateETHTransactionParams, chainID string) (string, error) {
	logger := uc.logger.WithContext(ctx)
	privKey, err := crypto.GenerateKey()
	if err != nil {
		errMessage := "failed to generate ethereum private key"
		logger.WithError(err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	decodedSignature, err := pkgcryto.SignEEATransaction(transaction, privateArgs, chainID, privKey)
	if err != nil {
		logger.WithError(err).Error("failed to sign EEA transaction")
		return "", err
	}

	signedRaw, err := uc.getSignedRawEEATransaction(ctx, transaction, privateArgs, decodedSignature, chainID)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signEEATransactionComponent)
	}

	logger.Debug("eea transaction signed successfully")

	return hexutil.Encode(signedRaw), nil
}

func (uc *signEEATransactionUseCase) signWithAccount(ctx context.Context, job *entities.Job,
	privateArgs *entities.PrivateETHTransactionParams, tx *types.Transaction, chainID string) (string, error) {
	logger := uc.logger.WithContext(ctx)

	signedRaw, err := uc.keyManagerClient.SignEEATransaction(ctx, uc.storeName, job.Transaction.From, &qkmtypes.SignEEATransactionRequest{
		Nonce:          hexutil.Uint64(tx.Nonce()),
		To:             tx.To(),
		Data:           tx.Data(),
		Value:          hexutil.Big(*tx.Value()),
		GasPrice:       hexutil.Big(*tx.GasPrice()),
		GasLimit:       hexutil.Uint64(tx.Gas()),
		ChainID:        hexutil.Big(*utils.MustEncodeBigInt(chainID)),
		PrivateFrom:    privateArgs.PrivateFrom,
		PrivateFor:     privateArgs.PrivateFor,
		PrivacyGroupID: privateArgs.PrivacyGroupID,
	})

	if err != nil {
		errMsg := "failed to sign eea transaction using key manager"
		logger.WithError(err).Error(errMsg)
		return "", errors.DependencyFailureError(errMsg).AppendReason(err.Error())
	}

	return signedRaw, nil
}

func (uc *signEEATransactionUseCase) getSignedRawEEATransaction(ctx context.Context, transaction *types.Transaction,
	privateArgs *entities.PrivateETHTransactionParams, signature []byte, chainID string) ([]byte, error) {
	logger := uc.logger.WithContext(ctx)

	signedTx, err := transaction.WithSignature(pkgcryto.GetEIP155Signer(chainID), signature)
	if err != nil {
		errMessage := "failed to set eea transaction signature"
		logger.WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}
	v, r, s := signedTx.RawSignatureValues()

	privateFromEncoded, err := pkgcryto.GetEncodedPrivateFrom(privateArgs.PrivateFrom)
	if err != nil {
		return nil, err
	}

	privateRecipientEncoded, err := pkgcryto.GetEncodedPrivateRecipient(privateArgs.PrivacyGroupID, privateArgs.PrivateFor)
	if err != nil {
		return nil, err
	}

	signedRaw, err := rlp.Encode([]interface{}{
		transaction.Nonce(),
		transaction.GasPrice(),
		transaction.Gas(),
		transaction.To(),
		transaction.Value(),
		transaction.Data(),
		v,
		r,
		s,
		privateFromEncoded,
		privateRecipientEncoded,
		privateArgs.PrivateTxType,
	})
	if err != nil {
		errMessage := "failed to RLP encode signed eea transaction"
		logger.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	return signedRaw, nil
}
