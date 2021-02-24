package signer

import (
	"context"
	"fmt"

	pkgcryto "github.com/ConsenSys/orchestrate/pkg/crypto/ethereum"
	"github.com/ConsenSys/orchestrate/pkg/encoding/rlp"
	"github.com/ConsenSys/orchestrate/pkg/log"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/pkg/utils"
	"github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/parsers"
	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	usecases "github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/use-cases"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/types/keymanager/ethereum"
	"github.com/ConsenSys/orchestrate/services/key-manager/client"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const signQuorumPrivateTransactionComponent = "use-cases.sign-quorum-private-transaction"

// signQuorumPrivateTransactionUseCase is a use case to sign a quorum private transaction
type signQuorumPrivateTransactionUseCase struct {
	keyManagerClient client.KeyManagerClient
	logger           *log.Logger
}

// NewSignQuorumPrivateTransactionUseCase creates a new signQuorumPrivateTransactionUseCase
func NewSignQuorumPrivateTransactionUseCase(keyManagerClient client.KeyManagerClient) usecases.SignQuorumPrivateTransactionUseCase {
	return &signQuorumPrivateTransactionUseCase{
		keyManagerClient: keyManagerClient,
		logger:           log.NewLogger().SetComponent(signQuorumPrivateTransactionComponent),
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
		decodedSignature, err = uc.signWithAccount(ctx, job, transaction)
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

func (uc *signQuorumPrivateTransactionUseCase) signWithAccount(ctx context.Context, job *entities.Job, tx *quorumtypes.Transaction) ([]byte, error) {
	logger := uc.logger.WithContext(ctx)
	request := &ethereum.SignQuorumPrivateTransactionRequest{
		Namespace: job.TenantID,
		Nonce:     tx.Nonce(),
		Amount:    tx.Value().String(),
		GasPrice:  tx.GasPrice().String(),
		GasLimit:  tx.Gas(),
		Data:      hexutil.Encode(tx.Data()),
	}
	if tx.To() != nil {
		request.To = tx.To().Hex()
	}

	tenants := utils.AllowedTenants(job.TenantID)
	for _, tenant := range tenants {
		request.Namespace = tenant
		sig, err := uc.keyManagerClient.ETHSignQuorumPrivateTransaction(ctx, job.Transaction.From, request)
		if err != nil && errors.IsNotFoundError(err) {
			continue
		}
		if err != nil {
			logger.WithError(err).Error("failed to sign quorum private transaction using key manager")
			return nil, errors.FromError(err)
		}

		decodedSignature, err := hexutil.Decode(sig)
		if err != nil {
			errMessage := "failed to decode quorum signature"
			logger.WithError(err).Error(errMessage)
			return nil, errors.EncodingError(errMessage)
		}

		return decodedSignature, nil
	}

	errMessage := fmt.Sprintf("account %s was not found on key-manager", job.Transaction.From)
	logger.WithField("from_account", job.Transaction.From).WithField("tenants", tenants).Error(errMessage)
	return nil, errors.InvalidParameterError(errMessage)
}
