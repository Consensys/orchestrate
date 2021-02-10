package signer

import (
	"context"
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/rlp"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"

	pkgcryto "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/crypto/ethereum"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"

	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/tx-sender/parsers"

	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/tx-sender/use-cases"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
)

const signEEATransactionComponent = "use-cases.sign-eea-transaction"

// signEEATransactionUseCase is a use case to sign a public Ethereum transaction
type signEEATransactionUseCase struct {
	keyManagerClient client.KeyManagerClient
	logger           *log.Logger
}

// NewSignEEATransactionUseCase creates a new SignEEATransactionUseCase
func NewSignEEATransactionUseCase(keyManagerClient client.KeyManagerClient) usecases.SignEEATransactionUseCase {
	return &signEEATransactionUseCase{
		keyManagerClient: keyManagerClient,
		logger:           log.NewLogger().SetComponent(signEEATransactionComponent),
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
		decodedSignature, err = uc.signWithAccount(ctx, job, transaction)
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
	tx *types.Transaction) ([]byte, error) {
	logger := uc.logger.WithContext(ctx)
	request := &ethereum.SignEEATransactionRequest{
		Namespace:      job.TenantID,
		Nonce:          tx.Nonce(),
		Data:           hexutil.Encode(tx.Data()),
		ChainID:        job.InternalData.ChainID,
		PrivateFrom:    job.Transaction.PrivateFrom,
		PrivateFor:     job.Transaction.PrivateFor,
		PrivacyGroupID: job.Transaction.PrivacyGroupID,
	}
	if tx.To() != nil {
		request.To = tx.To().Hex()
	}

	tenants := utils.AllowedTenants(job.TenantID)
	for _, tenant := range tenants {
		request.Namespace = tenant
		sig, err := uc.keyManagerClient.ETHSignEEATransaction(ctx, job.Transaction.From, request)
		if err != nil && errors.IsNotFoundError(err) {
			continue
		}
		if err != nil {
			logger.Error("failed to sign eea transaction using key manager")
			return nil, errors.FromError(err)
		}

		decodedSignature, err := hexutil.Decode(sig)
		if err != nil {
			errMessage := "failed to decode signature for eea transaction"
			logger.WithError(err).Error(errMessage)
			return nil, errors.EncodingError(errMessage)
		}

		return decodedSignature, nil
	}

	errMessage := fmt.Sprintf("account %s was not found on key-manager", job.Transaction.From)
	logger.WithField("from_account", job.Transaction.From).WithField("tenants", tenants).Error(errMessage)
	return nil, errors.InvalidParameterError(errMessage)
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
