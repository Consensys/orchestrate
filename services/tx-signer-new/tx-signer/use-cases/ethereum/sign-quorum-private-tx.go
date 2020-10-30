package ethereum

import (
	"context"

	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/crypto/ethereum/signing"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/rlp"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer-new/tx-signer/parsers"

	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer-new/tx-signer/use-cases"

	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"
)

const signQuorumPrivateTransactionComponent = "use-cases.sign-quorum-private-transaction"

// signQuorumPrivateTransactionUseCase is a use case to sign a quorum private transaction
type signQuorumPrivateTransactionUseCase struct {
	keyManagerClient client.KeyManagerClient
}

// NewSignQuorumPrivateTransactionUseCase creates a new signQuorumPrivateTransactionUseCase
func NewSignQuorumPrivateTransactionUseCase(keyManagerClient client.KeyManagerClient) usecases.SignQuorumPrivateTransactionUseCase {
	return &signQuorumPrivateTransactionUseCase{
		keyManagerClient: keyManagerClient,
	}
}

// Execute signs a quorum private transaction
func (uc *signQuorumPrivateTransactionUseCase) Execute(ctx context.Context, job *entities.Job) (raw, txHash string, err error) {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("signing quorum private transaction")

	signer := signing.GetQuorumPrivateTxSigner()
	transaction := parsers.ETHTransactionToQuorumTransaction(job.Transaction)
	transaction.SetPrivate()

	var decodedSignature []byte
	if job.InternalData.OneTimeKey {
		decodedSignature, err = uc.signWithOneTimeKey(transaction, signer)
	} else {
		decodedSignature, err = uc.signWithAccount(ctx, job, transaction)
	}
	if err != nil {
		return "", "", errors.FromError(err).ExtendComponent(signQuorumPrivateTransactionComponent)
	}

	signedTx, err := transaction.WithSignature(signer, decodedSignature)
	if err != nil {
		errMessage := "failed to set quorum private transaction signature"
		log.WithError(err).Error(errMessage)
		return "", "", errors.InvalidParameterError(errMessage).ExtendComponent(signQuorumPrivateTransactionComponent)
	}

	signedRaw, err := rlp.Encode(signedTx)
	if err != nil {
		errMessage := "failed to RLP encode signed quorum private transaction"
		log.WithError(err).Error(errMessage)
		return "", "", errors.CryptoOperationError(errMessage).ExtendComponent(signQuorumPrivateTransactionComponent)
	}
	txHash = transaction.Hash().Hex()

	logger.WithField("txHash", txHash).Info("quorum private transaction signed successfully")
	return hexutil.Encode(signedRaw), txHash, nil
}

func (*signQuorumPrivateTransactionUseCase) signWithOneTimeKey(transaction *quorumtypes.Transaction, signer quorumtypes.Signer) ([]byte, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		errMessage := "failed to generate Ethereum account"
		log.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	return signing.SignQuorumPrivateTransaction(transaction, privKey, signer)
}

func (uc *signQuorumPrivateTransactionUseCase) signWithAccount(ctx context.Context, job *entities.Job, tx *quorumtypes.Transaction) ([]byte, error) {
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

	sig, err := uc.keyManagerClient.ETHSignQuorumPrivateTransaction(ctx, job.Transaction.From, request)
	if err != nil {
		log.WithError(err).Error("failed to sign quorum private transaction using key manager")
		return nil, errors.FromError(err)
	}

	decodedSignature, err := hexutil.Decode(sig)
	if err != nil {
		errMessage := "failed to decode quorum signature"
		log.WithField("encoded_signature", sig).WithError(err).Error(errMessage)
		return nil, errors.EncodingError(errMessage)
	}

	return decodedSignature, nil
}
