package ethereum

import (
	"context"

	"github.com/consensys/quorum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"
)

const signPayloadComponent = "use-cases.ethereum.sign-payload"

// signPayloadUseCase is a use case to sign an arbitrary payload usign an existing Ethereum account
type signPayloadUseCase struct {
	vault store.Vault
}

// NewSignUseCase creates a new SignUseCase
func NewSignUseCase(vault store.Vault) SignUseCase {
	return &signPayloadUseCase{
		vault: vault,
	}
}

// Execute signs an arbitrary payload using an existing Ethereum account
func (uc *signPayloadUseCase) Execute(ctx context.Context, address, namespace, data string) (string, error) {
	logger := log.WithContext(ctx).WithField("namespace", namespace).WithField("address", address)
	logger.Debug("signing message")

	retrievedPrivKey, err := uc.vault.Ethereum().FindOne(ctx, address, namespace)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signPayloadComponent)
	}

	privKey, err := NewECDSAFromPrivKey(retrievedPrivKey)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signPayloadComponent)
	}

	signature, err := crypto.Sign(crypto.Keccak256([]byte(data)), privKey)
	if err != nil {
		errMessage := "failed to sign payload using ECDSA"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage).ExtendComponent(signPayloadComponent)
	}

	logger.Info("payload signed successfully")
	return hexutil.Encode(signature), nil
}
