package signer

import (
	"context"
	"math/big"

	pkgcryto "github.com/consensys/orchestrate/pkg/crypto/ethereum"
	"github.com/consensys/orchestrate/pkg/encoding/rlp"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/formatters"
	"github.com/consensys/orchestrate/pkg/utils"
	qkmtypes "github.com/consensys/quorum-key-manager/src/stores/api/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	usecases "github.com/consensys/orchestrate/services/tx-sender/tx-sender/use-cases"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/client"
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

func (uc *signETHTransactionUseCase) Execute(ctx context.Context, job *entities.Job) (signedRaw hexutil.Bytes, txHash *ethcommon.Hash, err error) {
	logger := uc.logger.WithContext(ctx).WithField("one_time_key", job.InternalData.OneTimeKey)

	transaction := formatters.ETHTransactionToTransaction(job.Transaction, job.InternalData.ChainID)
	if job.InternalData.OneTimeKey {
		signedRaw, txHash, err = uc.signWithOneTimeKey(ctx, transaction, job.InternalData.ChainID)
	} else {
		signedRaw, txHash, err = uc.signWithAccount(ctx, job, transaction, job.InternalData.ChainID)
	}
	if err != nil {
		return nil, nil, errors.FromError(err).ExtendComponent(signTransactionComponent)
	}

	logger.WithField("tx_hash", txHash.Hex()).Debug("ethereum transaction signed successfully")
	return signedRaw, txHash, nil
}

func (uc *signETHTransactionUseCase) signWithOneTimeKey(ctx context.Context, transaction *types.Transaction, chainID *big.Int) (signedRaw hexutil.Bytes, txHash *ethcommon.Hash, err error) {
	logger := uc.logger.WithContext(ctx)
	privKey, err := crypto.GenerateKey()
	if err != nil {
		errMessage := "failed to generate Ethereum private one time key"
		logger.WithError(err).Error(errMessage)
		return nil, nil, errors.CryptoOperationError(errMessage)
	}

	signer := types.NewEIP155Signer(chainID)
	decodedSignature, err := pkgcryto.SignTransaction(transaction, privKey, signer)
	if err != nil {
		logger.WithError(err).Error("failed to sign Ethereum transaction")
		return nil, nil, err
	}

	signedTx, err := transaction.WithSignature(signer, decodedSignature)
	if err != nil {
		errMessage := "failed to set transaction signature"
		logger.WithError(err).Error(errMessage)
		return nil, nil, errors.InvalidParameterError(errMessage).ExtendComponent(signTransactionComponent)
	}

	signedRawB, err := rlp.Encode(signedTx)
	if err != nil {
		errMessage := "failed to RLP encode signed transaction"
		logger.WithError(err).Error(errMessage)
		return nil, nil, errors.CryptoOperationError(errMessage).ExtendComponent(signTransactionComponent)
	}

	return signedRawB, utils.ToPtr(signedTx.Hash()).(*ethcommon.Hash), nil
}

func (uc *signETHTransactionUseCase) signWithAccount(ctx context.Context, job *entities.Job, tx *types.Transaction,
	chainID *big.Int) (signedRaw hexutil.Bytes, txHash *ethcommon.Hash, err error) {
	logger := uc.logger.WithContext(ctx)

	signedRawStr, err := uc.keyManagerClient.SignTransaction(ctx, job.InternalData.StoreID, job.Transaction.From.Hex(), &qkmtypes.SignETHTransactionRequest{
		Nonce:           hexutil.Uint64(tx.Nonce()),
		To:              tx.To(),
		Data:            tx.Data(),
		ChainID:         hexutil.Big(*chainID),
		Value:           hexutil.Big(*tx.Value()),
		GasPrice:        hexutil.Big(*tx.GasPrice()),
		GasLimit:        hexutil.Uint64(tx.Gas()),
		GasFeeCap:       utils.ToPtr(hexutil.Big(*tx.GasFeeCap())).(*hexutil.Big),
		GasTipCap:       utils.ToPtr(hexutil.Big(*tx.GasTipCap())).(*hexutil.Big),
		AccessList:      job.Transaction.AccessList,
		TransactionType: string(job.Transaction.TransactionType),
	})
	if err != nil {
		errMsg := "failed to sign ethereum transaction using key manager"
		logger.WithError(err).Error(errMsg)
		return nil, nil, errors.DependencyFailureError(errMsg).AppendReason(err.Error())
	}

	signedRaw, err = hexutil.Decode(signedRawStr)
	if err != nil {
		errMessage := "failed to decode raw signature"
		logger.WithError(err).Error(errMessage)
		return nil, nil, errors.EncodingError(errMessage)
	}

	err = tx.UnmarshalBinary(signedRaw)
	if err != nil {
		errMessage := "failed to decode tx signature"
		logger.WithError(err).Error(errMessage)
		return nil, nil, errors.EncodingError(errMessage)
	}

	return signedRaw, utils.ToPtr(tx.Hash()).(*ethcommon.Hash), nil
}
