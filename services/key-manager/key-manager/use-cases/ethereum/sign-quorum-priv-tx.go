package ethereum

import (
	"context"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/consensys/quorum/common/hexutil"
	quorumtypes "github.com/consensys/quorum/core/types"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"
)

const signQuorumPrivateTransactionComponent = "use-cases.ethereum.sign-quorum-private-transaction"

// signQuorumPrivateTxUseCase is a use case to sign a Quorum private transaction using an existing account
type signQuorumPrivateTxUseCase struct {
	vault store.Vault
}

// NewSignQuorumPrivateTransactionUseCase creates a new signQuorumPrivateTxUseCase
func NewSignQuorumPrivateTransactionUseCase(vault store.Vault) SignQuorumPrivateTransactionUseCase {
	return &signQuorumPrivateTxUseCase{
		vault: vault,
	}
}

// Execute signs a Quorum private transaction
func (uc *signQuorumPrivateTxUseCase) Execute(ctx context.Context, address, namespace string, tx *quorumtypes.Transaction) (string, error) {
	logger := log.WithContext(ctx).WithField("namespace", namespace).WithField("address", address)
	logger.Debug("signing quorum private transaction")

	retrievedPrivKey, err := uc.vault.Ethereum().FindOne(ctx, address, namespace)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signQuorumPrivateTransactionComponent)
	}

	privKey, err := NewECDSAFromPrivKey(retrievedPrivKey)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signQuorumPrivateTransactionComponent)
	}

	signer := quorumtypes.QuorumPrivateTxSigner{}
	tx.SetPrivate()
	h := signer.Hash(tx)
	signature, err := crypto.Sign(h[:], privKey)
	if err != nil {
		errMessage := "failed to sign quorum private transaction"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage).ExtendComponent(signQuorumPrivateTransactionComponent)
	}

	logger.Info("quorum private transaction signed successfully")
	return hexutil.Encode(signature), nil
}
