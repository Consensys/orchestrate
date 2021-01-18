package signer

import (
	"context"
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/rlp"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"

	"github.com/ethereum/go-ethereum/crypto"
	pkgcryto "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/crypto/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/tx-sender/parsers"

	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/tx-sender/use-cases"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
)

const signTransactionComponent = "use-cases.sign-eth-transaction"

// signETHTransactionUseCase is a use case to sign a public Ethereum transaction
type signETHTransactionUseCase struct {
	keyManagerClient client.KeyManagerClient
}

// NewSignETHTransactionUseCase creates a new SignTransactionUseCase
func NewSignETHTransactionUseCase(keyManagerClient client.KeyManagerClient) usecases.SignETHTransactionUseCase {
	return &signETHTransactionUseCase{
		keyManagerClient: keyManagerClient,
	}
}

func (uc *signETHTransactionUseCase) Execute(ctx context.Context, job *entities.Job) (raw, txHash string, err error) {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID).WithField("one_time_key", job.InternalData.OneTimeKey)
	logger.Debug("signing ethereum transaction")

	signer := pkgcryto.GetEIP155Signer(job.InternalData.ChainID)
	transaction := parsers.ETHTransactionToTransaction(job.Transaction)

	var decodedSignature []byte
	if job.InternalData.OneTimeKey {
		decodedSignature, err = uc.signWithOneTimeKey(transaction, signer)
	} else {
		decodedSignature, err = uc.signWithAccount(ctx, job, transaction)
	}
	if err != nil {
		return "", "", errors.FromError(err).ExtendComponent(signTransactionComponent)
	}

	signedTx, err := transaction.WithSignature(signer, decodedSignature)
	if err != nil {
		errMessage := "failed to set transaction signature"
		log.WithError(err).Error(errMessage)
		return "", "", errors.InvalidParameterError(errMessage).ExtendComponent(signTransactionComponent)
	}

	signedRaw, err := rlp.Encode(signedTx)
	if err != nil {
		errMessage := "failed to RLP encode signed transaction"
		log.WithError(err).Error(errMessage)
		return "", "", errors.CryptoOperationError(errMessage).ExtendComponent(signTransactionComponent)
	}
	txHash = signedTx.Hash().Hex()

	logger.WithField("txHash", txHash).Info("ethereum transaction signed successfully")
	return hexutil.Encode(signedRaw), txHash, nil
}

func (*signETHTransactionUseCase) signWithOneTimeKey(transaction *types.Transaction, signer types.Signer) ([]byte, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		errMessage := "failed to generate Ethereum private key"
		log.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	return pkgcryto.SignTransaction(transaction, privKey, signer)
}

func (uc *signETHTransactionUseCase) signWithAccount(ctx context.Context, job *entities.Job, tx *types.Transaction) ([]byte, error) {
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
			log.WithError(err).Error("failed to sign ethereum transaction using key manager")
			return nil, errors.FromError(err)
		}

		decodedSignature, err := hexutil.Decode(sig)
		if err != nil {
			errMessage := "failed to decode signature"
			log.WithField("encoded_signature", sig).WithError(err).Error(errMessage)
			return nil, errors.EncodingError(errMessage)
		}

		return decodedSignature, nil
	}

	errMessage := fmt.Sprintf("account %s was not found on key-manager", job.Transaction.From)
	log.WithField("from_account", job.Transaction.From).WithField("tenants", tenants).Error(errMessage)
	return nil, errors.InvalidParameterError(errMessage)
}
