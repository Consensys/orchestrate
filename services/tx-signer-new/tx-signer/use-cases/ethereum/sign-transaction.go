package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer-new/tx-signer/parsers"

	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer-new/tx-signer/use-cases"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"
)

const signTransactionComponent = "use-cases.sign-transaction"

// signTransactionUseCase is a use case to sign a public Ethereum transaction
type signTransactionUseCase struct {
	keyManagerClient client.KeyManagerClient
}

// NewSignTransactionUseCase creates a new SignTransactionUseCase
func NewSignTransactionUseCase(keyManagerClient client.KeyManagerClient) usecases.SignTransactionUseCase {
	return &signTransactionUseCase{
		keyManagerClient: keyManagerClient,
	}
}

// Execute signs a public Ethereum transaction
func (uc *signTransactionUseCase) Execute(ctx context.Context, job *entities.Job) (raw, txHash string, err error) {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("signing ethereum transaction")

	signer := getSigner(job.InternalData.ChainID)
	transaction := parsers.ETHTransactionToTransaction(job.Transaction)

	var decodedSignature []byte
	if job.InternalData.OneTimeKey {
		decodedSignature, err = signWithOneTimeKey(transaction, signer)
	} else {
		decodedSignature, err = uc.signWithAccount(ctx, job, transaction)
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

	signedRaw, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		errMessage := "failed to RLP encode signed transaction"
		logger.WithError(err).Error(errMessage)
		return "", "", errors.CryptoOperationError(errMessage).ExtendComponent(signTransactionComponent)
	}
	txHash = transaction.Hash().Hex()

	logger.WithField("txHash", txHash).Info("ethereum transaction signed successfully")
	return hexutil.Encode(signedRaw), txHash, nil
}

func getSigner(chainID string) types.Signer {
	chainIDBigInt := new(big.Int)
	chainIDBigInt, _ = chainIDBigInt.SetString(chainID, 10)
	return types.NewEIP155Signer(chainIDBigInt)
}

func signWithOneTimeKey(transaction *types.Transaction, signer types.Signer) ([]byte, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		errMessage := "failed to generate Ethereum private key"
		log.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	h := signer.Hash(transaction)
	decodedSignature, err := crypto.Sign(h[:], privKey)
	if err != nil {
		errMessage := "failed to sign ethereum transaction"
		log.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	return decodedSignature, nil
}

func (uc *signTransactionUseCase) signWithAccount(ctx context.Context, job *entities.Job, tx *types.Transaction) ([]byte, error) {
	request := &ethereum.SignETHTransactionRequest{
		Namespace: job.TenantID,
		Nonce:     tx.Nonce(),
		Amount:    tx.Value().String(),
		GasPrice:  tx.GasPrice().String(),
		GasLimit:  tx.Gas(),
		Data:      hexutil.Encode(tx.Data()),
		ChainID:   job.InternalData.ChainID,
	}
	if tx.To() != nil {
		request.To = tx.To().Hex()
	}

	sig, err := uc.keyManagerClient.ETHSignTransaction(ctx, job.Transaction.From, request)
	if err != nil {
		log.WithError(err).Error("failed to sign ethereum transaction using key manager")
		return nil, errors.FromError(err).ExtendComponent(signTransactionComponent)
	}

	decodedSignature, err := hexutil.Decode(sig)
	if err != nil {
		errMessage := "failed to decode signature"
		log.WithField("encoded_signature", sig).WithError(err).Error(errMessage)
		return nil, errors.EncodingError(errMessage).ExtendComponent(signTransactionComponent)
	}

	return decodedSignature, nil
}
