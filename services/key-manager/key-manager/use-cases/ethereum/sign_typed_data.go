package ethereum

import (
	"context"
	"fmt"

	signer "github.com/ethereum/go-ethereum/signer/core"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/service/formatters"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"
)

const signTypedDataComponent = "use-cases.sign-typed-data"

// signTypedDataUseCase is a use case to sign an arbitrary typed payload usign an existing Ethereum account
type signTypedDataUseCase struct {
	vaultClient store.Vault
}

// NewSignTypedDataUseCase creates a new SignTypedDataUseCase
func NewSignTypedDataUseCase(vaultClient store.Vault) usecases.SignTypedDataUseCase {
	return &signTypedDataUseCase{
		vaultClient: vaultClient,
	}
}

// Execute signs an arbitrary payload using an existing Ethereum account
func (uc *signTypedDataUseCase) Execute(ctx context.Context, address, namespace string, typedData *signer.TypedData) (string, error) {
	logger := log.WithContext(ctx).WithField("namespace", namespace).WithField("address", address)
	logger.Debug("signing typed data")

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		errMessage := "invalid typed data message"
		logger.WithError(err).Error(errMessage)
		return "", errors.InvalidParameterError(fmt.Sprintf("%s: %s", errMessage, err.Error())).ExtendComponent(signTypedDataComponent)
	}

	domainSeparatorHash, err := typedData.HashStruct(formatters.DomainLabel, typedData.Domain.Map())
	if err != nil {
		errMessage := "invalid domain separator"
		logger.WithError(err).Error(errMessage)
		return "", errors.InvalidParameterError(fmt.Sprintf("%s: %s", errMessage, err.Error())).ExtendComponent(signTypedDataComponent)
	}

	encodedData := fmt.Sprintf("\x19\x01%s%s", domainSeparatorHash, typedDataHash)
	signature, err := uc.vaultClient.ETHSign(address, namespace, encodedData)
	if err != nil {
		errMessage := "failed to sign typed data"
		logger.WithError(err).Error(errMessage)
		return "", errors.FromError(err).ExtendComponent(signTypedDataComponent)
	}

	logger.Info("typed data signed successfully")
	return signature, nil
}
