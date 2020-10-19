package ethereum

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/consensys/quorum/common/hexutil"

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
func (uc *createAccountUseCase) Execute(ctx context.Context, namespace, importedPrivKey string) (*entities.ETHAccount, error) {
	logger := log.WithContext(ctx).WithField("namespace", namespace)
	logger.Debug("creating new Ethereum account")

	var privKey = new(ecdsa.PrivateKey)
	var err error
	if importedPrivKey == "" {
		privKey, err = generatePrivKey(ctx)
	} else {
		privKey, err = retrievePrivKey(ctx, importedPrivKey)
	}
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
	}

	account := &entities.ETHAccount{
		Address:             crypto.PubkeyToAddress(privKey.PublicKey).Hex(),
		PublicKey:           hexutil.Encode(crypto.FromECDSAPub(&privKey.PublicKey)),
		CompressedPublicKey: hexutil.Encode(crypto.CompressPubkey(&privKey.PublicKey)),
		Namespace:           namespace,
	}

	err = uc.vault.Ethereum().Insert(ctx, account.Address, hex.EncodeToString(crypto.FromECDSA(privKey)), account.Namespace)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
	}

	logger.WithField("address", account.Address).Info("Ethereum account created successfully")
	return account, nil
}

func retrievePrivKey(ctx context.Context, importedPrivKey string) (*ecdsa.PrivateKey, error) {
	privKey, err := crypto.HexToECDSA(importedPrivKey)
	if err != nil {
		errMessage := "failed to import Ethereum private key, please verify that the provided private key is valid"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	return privKey, nil
}

func generatePrivKey(ctx context.Context) (*ecdsa.PrivateKey, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		errMessage := "failed to generate Ethereum private key"
		log.WithContext(ctx).WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	return privKey, nil
}
