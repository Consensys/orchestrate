package ethereum

import (
	"context"

	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/crypto/ethereum/signing"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer-new/tx-signer/parsers"

	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer-new/tx-signer/use-cases"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
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

	signer := signing.GetEIP155Signer(job.InternalData.ChainID)
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

	signedRaw, err := GetSignedRawTransaction(transaction, decodedSignature, signer)
	if err != nil {
		return "", "", errors.FromError(err).ExtendComponent(signEEATransactionComponent)
	}
	txHash = transaction.Hash().Hex()

	logger.WithField("txHash", txHash).Info("ethereum transaction signed successfully")
	return signedRaw, txHash, nil
}

func signWithOneTimeKey(transaction *types.Transaction, signer types.Signer) ([]byte, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		errMessage := "failed to generate Ethereum private key"
		log.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	return signing.SignETHTransaction(transaction, privKey, signer)
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
