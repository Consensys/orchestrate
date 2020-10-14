package ethereum

import (
	"context"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store"
)

const createAccountComponent = "use-cases.ethereum.create-account"

// createAccountUseCase is a use case to create a new Ethereum account
type createAccountUseCase struct {
	vault store.Vault
}

// NewCreateAccountUseCase creates a new CreateAccountUseCase
func NewCreateAccountUseCase(vault store.Vault) CreateAccountUseCase {
	return &createAccountUseCase{
		vault: vault,
	}
}

// Execute creates a new Ethereum account and stores it in the Vault
func (uc *createAccountUseCase) Execute(ctx context.Context, account *entities.ETHAccount) (*entities.ETHAccount, error) {
	logger := log.WithContext(ctx).
		WithField("key_type", account.KeyType).
		WithField("namespace", account.Namespace)
	logger.Debug("creating new Ethereum account")

	// TODO: Verify keyType here and branch between sub use cases to create different keys given the elliptic curve
	// TODO: Currently not needed as only keyType implementation is Secp256k1 for ETH1
	privKey, err := crypto.GenerateKey()
	if err != nil {
		errMessage := "failed to generate Ethereum private key"
		logger.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage).ExtendComponent(createAccountComponent)
	}

	account.PublicKey = hex.EncodeToString(crypto.FromECDSAPub(&privKey.PublicKey))
	account.CompressedPublicKey = hex.EncodeToString(crypto.CompressPubkey(&privKey.PublicKey))
	account.Address = crypto.PubkeyToAddress(privKey.PublicKey).Hex()

	err = uc.vault.Ethereum().Insert(ctx, account.Address, hex.EncodeToString(crypto.FromECDSA(privKey)), account.Namespace)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
	}

	logger.WithField("address", account.Address).Info("Ethereum account created successfully")
	return account, nil
}
