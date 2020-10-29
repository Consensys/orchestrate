package ethereum

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/crypto/ethereum/signing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store"
)

const signEEATransactionComponent = "use-cases.ethereum.sign-eea-transaction"

// signEEATxUseCase is a use case to sign a Quorum private transaction using an existing account
type signEEATxUseCase struct {
	vault store.Vault
}

// NewSignEEATransactionUseCase creates a new signEEATxUseCase
func NewSignEEATransactionUseCase(vault store.Vault) SignEEATransactionUseCase {
	return &signEEATxUseCase{
		vault: vault,
	}
}

// Execute signs a Quorum private transaction
func (uc *signEEATxUseCase) Execute(
	ctx context.Context,
	address, namespace, chainID string,
	tx *ethtypes.Transaction,
	privateArgs *entities.PrivateETHTransactionParams,
) (string, error) {
	logger := log.WithContext(ctx).WithField("namespace", namespace).WithField("address", address)
	logger.Debug("signing eea private transaction")

	retrievedPrivKey, err := uc.vault.Ethereum().FindOne(ctx, address, namespace)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signEEATransactionComponent)
	}

	privKey, err := NewECDSAFromPrivKey(retrievedPrivKey)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signEEATransactionComponent)
	}

	signature, err := signing.SignEEATransaction(tx, privateArgs, chainID, privKey)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signEEATransactionComponent)
	}

	logger.Info("eea private transaction signed successfully")
	return hexutil.Encode(signature), nil
}
