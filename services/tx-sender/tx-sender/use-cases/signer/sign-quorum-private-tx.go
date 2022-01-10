package signer

import (
	"context"

	pkgcryto "github.com/consensys/orchestrate/pkg/crypto/ethereum"
	"github.com/consensys/orchestrate/pkg/encoding/rlp"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/formatters"
	"github.com/consensys/orchestrate/pkg/utils"
	qkmtypes "github.com/consensys/quorum-key-manager/src/stores/api/types"
	quorumtypes "github.com/consensys/quorum/core/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	usecases "github.com/consensys/orchestrate/services/tx-sender/tx-sender/use-cases"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/client"
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
func (uc *signQuorumPrivateTransactionUseCase) Execute(ctx context.Context, job *entities.Job) (signedRaw hexutil.Bytes, txHash *ethcommon.Hash, err error) {
	logger := uc.logger.WithContext(ctx).WithField("one_time_key", job.InternalData.OneTimeKey)

	transaction := formatters.ETHTransactionToQuorumTransaction(job.Transaction)
	transaction.SetPrivate()
	if job.InternalData.OneTimeKey {
		signedRaw, txHash, err = uc.signWithOneTimeKey(ctx, transaction)
	} else {
		signedRaw, txHash, err = uc.signWithAccount(ctx, job, transaction)
	}

	if err != nil {
		return nil, nil, errors.FromError(err).ExtendComponent(signQuorumPrivateTransactionComponent)
	}

	logger.WithField("tx_hash", txHash).Debug("quorum private transaction signed successfully")
	return signedRaw, txHash, nil
}

func (uc *signQuorumPrivateTransactionUseCase) signWithOneTimeKey(ctx context.Context, transaction *quorumtypes.Transaction) (
	signedRaw hexutil.Bytes, txHash *ethcommon.Hash, err error) {
	logger := uc.logger.WithContext(ctx)
	privKey, err := crypto.GenerateKey()
	if err != nil {
		errMessage := "failed to generate Ethereum account"
		logger.WithError(err).Error(errMessage)
		return nil, nil, errors.CryptoOperationError(errMessage)
	}

	signer := pkgcryto.GetQuorumPrivateTxSigner()
	decodedSignature, err := pkgcryto.SignQuorumPrivateTransaction(transaction, privKey, signer)
	if err != nil {
		logger.WithError(err).Error("failed to sign private transaction")
		return nil, nil, err
	}

	signedTx, err := transaction.WithSignature(signer, decodedSignature)
	if err != nil {
		errMessage := "failed to set quorum private transaction signature"
		logger.WithError(err).Error(errMessage)
		return nil, nil, errors.InvalidParameterError(errMessage).ExtendComponent(signQuorumPrivateTransactionComponent)
	}

	signedRawB, err := rlp.Encode(signedTx)
	if err != nil {
		errMessage := "failed to RLP encode signed quorum private transaction"
		logger.WithError(err).Error(errMessage)
		return nil, nil, errors.CryptoOperationError(errMessage).ExtendComponent(signQuorumPrivateTransactionComponent)
	}

	return signedRawB, utils.ToPtr(signedTx.Hash()).(*ethcommon.Hash), nil
}

func (uc *signQuorumPrivateTransactionUseCase) signWithAccount(ctx context.Context, job *entities.Job, tx *quorumtypes.Transaction) (
	signedRaw hexutil.Bytes, txHash *ethcommon.Hash, err error) {
	logger := uc.logger.WithContext(ctx)

	signedRawStr, err := uc.keyManagerClient.SignQuorumPrivateTransaction(ctx, job.InternalData.StoreID, job.Transaction.From.Hex(), &qkmtypes.SignQuorumPrivateTransactionRequest{
		Nonce:    hexutil.Uint64(tx.Nonce()),
		To:       tx.To(),
		Value:    hexutil.Big(*tx.Value()),
		GasPrice: hexutil.Big(*tx.GasPrice()),
		GasLimit: hexutil.Uint64(tx.Gas()),
		Data:     tx.Data(),
	})
	if err != nil {
		errMsg := "failed to sign quorum private transaction using key manager"
		logger.WithError(err).Error(errMsg)
		return nil, nil, errors.DependencyFailureError(errMsg).AppendReason(err.Error())
	}

	signedRaw, err = hexutil.Decode(signedRawStr)
	if err != nil {
		errMessage := "failed to decode quorum raw signature"
		logger.WithError(err).Error(errMessage)
		return nil, nil, errors.EncodingError(errMessage)
	}

	err = rlp.Decode(signedRaw, &tx)
	if err != nil {
		errMessage := "failed to decode quorum transaction"
		logger.WithError(err).Error(errMessage)
		return nil, nil, errors.EncodingError(errMessage)
	}

	return signedRaw, utils.ToPtr(tx.Hash()).(*ethcommon.Hash), nil
}
