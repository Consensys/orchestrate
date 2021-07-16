package signer

import (
	"context"

	pkgcryto "github.com/ConsenSys/orchestrate/pkg/crypto/ethereum"
	"github.com/ConsenSys/orchestrate/pkg/encoding/rlp"
	qkm "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/pkg/utils"
	"github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/parsers"
	qkmtypes "github.com/consensys/quorum-key-manager/src/stores/api/types"
	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	usecases "github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/use-cases"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const signQuorumPrivateTransactionComponent = "use-cases.sign-quorum-private-transaction"

// signQuorumPrivateTransactionUseCase is a use case to sign a quorum private transaction
type signQuorumPrivateTransactionUseCase struct {
	keyManagerClient client.KeyManagerClient
	logger           *log.Logger
	storeName        string
}

// NewSignQuorumPrivateTransactionUseCase creates a new signQuorumPrivateTransactionUseCase
func NewSignQuorumPrivateTransactionUseCase(keyManagerClient client.KeyManagerClient) usecases.SignQuorumPrivateTransactionUseCase {
	return &signQuorumPrivateTransactionUseCase{
		keyManagerClient: keyManagerClient,
		logger:           log.NewLogger().SetComponent(signQuorumPrivateTransactionComponent),
		storeName:        qkm.GlobalStoreName(),
	}
}

// Execute signs a quorum private transaction
func (uc *signQuorumPrivateTransactionUseCase) Execute(ctx context.Context, job *entities.Job) (raw, txHash string, err error) {
	logger := uc.logger.WithContext(ctx).WithField("one_time_key", job.InternalData.OneTimeKey)

	signer := pkgcryto.GetQuorumPrivateTxSigner()
	transaction := parsers.ETHTransactionToQuorumTransaction(job.Transaction)
	transaction.SetPrivate()

	var decodedSignature []byte
	if job.InternalData.OneTimeKey {
		decodedSignature, err = uc.signWithOneTimeKey(ctx, transaction, signer)
	} else {
		decodedSignature, err = uc.signWithAccount(ctx, job, transaction, signer)
	}
	if err != nil {
		return "", "", errors.FromError(err).ExtendComponent(signQuorumPrivateTransactionComponent)
	}

	signedTx, err := transaction.WithSignature(signer, decodedSignature)
	if err != nil {
		errMessage := "failed to set quorum private transaction signature"
		logger.WithError(err).Error(errMessage)
		return "", "", errors.InvalidParameterError(errMessage).ExtendComponent(signQuorumPrivateTransactionComponent)
	}

	signedRaw, err := rlp.Encode(signedTx)
	if err != nil {
		errMessage := "failed to RLP encode signed quorum private transaction"
		logger.WithError(err).Error(errMessage)
		return "", "", errors.CryptoOperationError(errMessage).ExtendComponent(signQuorumPrivateTransactionComponent)
	}
	txHash = signedTx.Hash().Hex()

	logger.WithField("tx_hash", txHash).Debug("quorum private transaction signed successfully")
	return hexutil.Encode(signedRaw), txHash, nil
}

func (uc *signQuorumPrivateTransactionUseCase) signWithOneTimeKey(ctx context.Context, transaction *quorumtypes.Transaction,
	signer quorumtypes.Signer) ([]byte, error) {
	logger := uc.logger.WithContext(ctx)
	privKey, err := crypto.GenerateKey()
	if err != nil {
		errMessage := "failed to generate Ethereum account"
		logger.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	sign, err := pkgcryto.SignQuorumPrivateTransaction(transaction, privKey, signer)
	if err != nil {
		logger.WithError(err).Error("failed to sign private transaction")
		return nil, err
	}

	return sign, nil
}

func (uc *signQuorumPrivateTransactionUseCase) signWithAccount(ctx context.Context, job *entities.Job,
	tx *quorumtypes.Transaction, signer quorumtypes.Signer) ([]byte, error) {
	logger := uc.logger.WithContext(ctx)
	tenants := utils.AllowedTenants(job.TenantID)
	isAllowed, err := qkm.IsTenantAllowed(ctx, uc.keyManagerClient, tenants, uc.storeName, job.Transaction.From)
	if err != nil {
		errMsg := "failed to to sign private quorum transaction, cannot fetch account"
		uc.logger.WithField("address", job.Transaction.From).WithError(err).Error(errMsg)
		return nil, errors.DependencyFailureError(errMsg).AppendReason(err.Error())
	}

	if !isAllowed {
		errMessage := "failed to to sign private quorum transaction, tenant is not allowed"
		logger.WithField("address", job.Transaction.From).WithField("tenants", tenants).Error(errMessage)
		return nil, errors.InvalidAuthenticationError(errMessage)
	}

	txData := signer.Hash(tx).Bytes()
	sig, err := uc.keyManagerClient.SignEth1Data(ctx, uc.storeName, job.Transaction.From, &qkmtypes.SignHexPayloadRequest{
		Data: txData,
	})
	if err != nil {
		errMsg := "failed to sign quorum private transaction using key manager"
		logger.Error(errMsg)
		return nil, errors.DependencyFailureError(errMsg).AppendReason(err.Error())
	}

	decodedSignature, err := hexutil.Decode(sig)
	if err != nil {
		errMessage := "failed to decode quorum signature"
		logger.WithError(err).Error(errMessage)
		return nil, errors.EncodingError(errMessage)
	}

	return decodedSignature, nil
}
