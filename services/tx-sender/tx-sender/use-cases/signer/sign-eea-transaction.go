package signer

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/encoding/rlp"
	qkm "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager"
	qkmtypes "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager/types"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"

	pkgcryto "github.com/ConsenSys/orchestrate/pkg/crypto/ethereum"

	"github.com/ConsenSys/orchestrate/pkg/utils"

	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/parsers"
	"github.com/ethereum/go-ethereum/crypto"

	usecases "github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/use-cases"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/quorum-key-manager/client"
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

func (uc *signEEATransactionUseCase) Execute(ctx context.Context, job *entities.Job) (raw, txHash string, err error) {
	logger := uc.logger.WithContext(ctx).WithField("one_time_key", job.InternalData.OneTimeKey)

	signer := pkgcryto.GetEIP155Signer(job.InternalData.ChainID)
	transaction := parsers.ETHTransactionToTransaction(job.Transaction)
	privateArgs := &entities.PrivateETHTransactionParams{
		PrivateFrom:    job.Transaction.PrivateFrom,
		PrivateFor:     job.Transaction.PrivateFor,
		PrivacyGroupID: job.Transaction.PrivacyGroupID,
		PrivateTxType:  entities.PrivateTxTypeRestricted,
	}

	var decodedSignature []byte
	if job.InternalData.OneTimeKey {
		decodedSignature, err = uc.signWithOneTimeKey(ctx, transaction, privateArgs, job.InternalData.ChainID)
	} else {
		decodedSignature, err = uc.signWithAccount(ctx, job, privateArgs, transaction)
	}
	if err != nil {
		return "", "", errors.FromError(err).ExtendComponent(signEEATransactionComponent)
	}

	signedRaw, err := uc.getSignedRawEEATransaction(ctx, transaction, privateArgs, decodedSignature, signer)
	if err != nil {
		return "", "", errors.FromError(err).ExtendComponent(signEEATransactionComponent)
	}

	logger.Debug("eea transaction signed successfully")

	// transaction hash of EEA transactions cannot be computed
	return hexutil.Encode(signedRaw), "", nil
}

func (uc *signEEATransactionUseCase) signWithOneTimeKey(ctx context.Context, transaction *types.Transaction,
	privateArgs *entities.PrivateETHTransactionParams, chainID string) ([]byte, error) {
	logger := uc.logger.WithContext(ctx)
	privKey, err := crypto.GenerateKey()
	if err != nil {
		errMessage := "failed to generate ethereum private key"
		logger.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	sign, err := pkgcryto.SignEEATransaction(transaction, privateArgs, chainID, privKey)
	if err != nil {
		logger.WithError(err).Error("failed to sign EEA transaction")
		return nil, err
	}

	return sign, nil
}

func (uc *signEEATransactionUseCase) signWithAccount(ctx context.Context, job *entities.Job,
	privateArgs *entities.PrivateETHTransactionParams, tx *types.Transaction) ([]byte, error) {
	logger := uc.logger.WithContext(ctx)

	tenants := utils.AllowedTenants(job.TenantID)
	isAllowed, err := qkm.IsTenantAllowed(ctx, uc.keyManagerClient, tenants, uc.storeName, job.Transaction.From)
	if err != nil {
		errMsg := "failed to to sign eea transaction, cannot fetch account"
		uc.logger.WithField("address", job.Transaction.From).WithError(err).Error(errMsg)
		return nil, errors.DependencyFailureError(errMsg).AppendReason(err.Error())
	}

	if !isAllowed {
		errMessage := "failed to to sign eea transaction, tenant is not allowed"
		logger.WithField("address", job.Transaction.From).WithField("tenants", tenants).Error(errMessage)
		return nil, errors.InvalidAuthenticationError(errMessage)
	}

	txData, err := pkgcryto.EEATransactionPayload(tx, privateArgs, job.InternalData.ChainID)
	if err != nil {
		errMsg := "failed to build eea transaction payload"
		uc.logger.WithField("address", job.Transaction.From).WithError(err).Error(errMsg)
		return nil, errors.FromError(err)
	}

	sig, err := uc.keyManagerClient.SignEth1Data(ctx, uc.storeName, job.Transaction.From, &qkmtypes.SignHexPayloadRequest{
		Data: txData,
	})
	if err != nil {
		errMsg := "failed to sign eea transaction using key manager"
		logger.Error(errMsg)
		return nil, errors.DependencyFailureError(errMsg).AppendReason(err.Error())
	}

	decodedSignature, err := hexutil.Decode(sig)
	if err != nil {
		errMessage := "failed to decode signature for eea transaction"
		logger.WithError(err).Error(errMessage)
		return nil, errors.EncodingError(errMessage)
	}

	return decodedSignature, nil
}

func (uc *signEEATransactionUseCase) getSignedRawEEATransaction(ctx context.Context, transaction *types.Transaction,
	privateArgs *entities.PrivateETHTransactionParams, signature []byte, signer types.Signer) ([]byte, error) {
	logger := uc.logger.WithContext(ctx)
	privateFromEncoded, err := pkgcryto.GetEncodedPrivateFrom(privateArgs.PrivateFrom)
	if err != nil {
		return nil, err
	}

	privateRecipientEncoded, err := pkgcryto.GetEncodedPrivateRecipient(privateArgs.PrivacyGroupID, privateArgs.PrivateFor)
	if err != nil {
		return nil, err
	}

	signedTx, err := transaction.WithSignature(signer, signature)
	if err != nil {
		errMessage := "failed to set eea transaction signature"
		logger.WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}
	v, r, s := signedTx.RawSignatureValues()

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
