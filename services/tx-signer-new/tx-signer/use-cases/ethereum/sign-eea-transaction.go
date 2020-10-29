package ethereum

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/crypto/ethereum/signing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"

	"github.com/ethereum/go-ethereum/crypto"
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

const signEEATransactionComponent = "use-cases.sign-eea-transaction"

// signEEATransactionUseCase is a use case to sign a public Ethereum transaction
type signEEATransactionUseCase struct {
	keyManagerClient client.KeyManagerClient
}

// NewSignEEATransactionUseCase creates a new SignEEATransactionUseCase
func NewSignEEATransactionUseCase(keyManagerClient client.KeyManagerClient) usecases.SignEEATransactionUseCase {
	return &signEEATransactionUseCase{
		keyManagerClient: keyManagerClient,
	}
}

// Execute signs a public Ethereum transaction
func (uc *signEEATransactionUseCase) Execute(ctx context.Context, job *entities.Job) (raw, txHash string, err error) {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("signing EEA transaction")

	signer := signing.GetEIP155Signer(job.InternalData.ChainID)
	transaction := parsers.ETHTransactionToTransaction(job.Transaction)

	var decodedSignature []byte
	if job.InternalData.OneTimeKey {
		decodedSignature, err = uc.signWithOneTimeKey(transaction, &entities.PrivateETHTransactionParams{
			PrivateFrom:    job.Transaction.PrivateFrom,
			PrivateFor:     job.Transaction.PrivateFor,
			PrivacyGroupID: job.Transaction.PrivacyGroupID,
			PrivateTxType:  utils.PrivateTxTypeRestricted,
		}, job.InternalData.ChainID)
	} else {
		decodedSignature, err = uc.signWithAccount(ctx, job, transaction)
	}
	if err != nil {
		return "", "", errors.FromError(err).ExtendComponent(signEEATransactionComponent)
	}

	signedRaw, err := GetSignedRawTransaction(transaction, decodedSignature, signer)
	if err != nil {
		return "", "", errors.FromError(err).ExtendComponent(signEEATransactionComponent)
	}
	txHash = transaction.Hash().Hex()

	logger.WithField("txHash", txHash).Info("eea transaction signed successfully")
	return signedRaw, txHash, nil
}

func (uc *signEEATransactionUseCase) signWithOneTimeKey(
	transaction *types.Transaction,
	privateArgs *entities.PrivateETHTransactionParams,
	chainID string,
) ([]byte, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		errMessage := "failed to generate ethereum private key"
		log.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	return signing.SignEEATransaction(transaction, privateArgs, chainID, privKey)
}

func (uc *signEEATransactionUseCase) signWithAccount(ctx context.Context, job *entities.Job, tx *types.Transaction) ([]byte, error) {
	request := &ethereum.SignEEATransactionRequest{
		Namespace:      job.TenantID,
		Nonce:          tx.Nonce(),
		Amount:         tx.Value().String(),
		GasPrice:       tx.GasPrice().String(),
		GasLimit:       tx.Gas(),
		Data:           hexutil.Encode(tx.Data()),
		ChainID:        job.InternalData.ChainID,
		PrivateFrom:    job.Transaction.PrivateFrom,
		PrivateFor:     job.Transaction.PrivateFor,
		PrivacyGroupID: job.Transaction.PrivacyGroupID,
	}
	if tx.To() != nil {
		request.To = tx.To().Hex()
	}

	sig, err := uc.keyManagerClient.ETHSignEEATransaction(ctx, job.Transaction.From, request)
	if err != nil {
		log.WithError(err).Error("failed to sign eea transaction using key manager")
		return nil, err
	}

	decodedSignature, err := hexutil.Decode(sig)
	if err != nil {
		errMessage := "failed to decode signature for eea transaction"
		log.WithField("encoded_signature", sig).WithError(err).Error(errMessage)
		return nil, errors.EncodingError(errMessage)
	}

	return decodedSignature, nil
}
