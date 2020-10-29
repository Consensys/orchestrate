package ethereum

import (
	"context"

	"github.com/consensys/quorum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/crypto/ethereum/signing"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store"
)

const signTransactionComponent = "use-cases.ethereum.sign-transaction"

// signTxUseCase is a use case to sign an ethereum transaction using an existing account
type signTxUseCase struct {
	vault store.Vault
}

// NewSignTransactionUseCase creates a new signTxUseCase
func NewSignTransactionUseCase(vault store.Vault) SignTransactionUseCase {
	return &signTxUseCase{
		vault: vault,
	}
}

// Execute signs an ethereum transaction
func (uc *signTxUseCase) Execute(ctx context.Context, address, namespace, chainID string, tx *ethtypes.Transaction) (string, error) {
	logger := log.WithContext(ctx).WithField("namespace", namespace).WithField("address", address)
	logger.Debug("signing ethereum transaction")

	retrievedPrivKey, err := uc.vault.Ethereum().FindOne(ctx, address, namespace)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signTransactionComponent)
	}

	privKey, err := NewECDSAFromPrivKey(retrievedPrivKey)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signTransactionComponent)
	}

	signature, err := signing.SignETHTransaction(tx, privKey, signing.GetEIP155Signer(chainID))
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signTransactionComponent)
	}

	logger.Info("ethereum transaction signed successfully")
	return hexutil.Encode(signature), nil
}
