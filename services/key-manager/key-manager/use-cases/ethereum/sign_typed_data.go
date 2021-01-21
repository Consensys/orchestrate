package ethereum

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases/ethereum/utils"

	signer "github.com/ethereum/go-ethereum/signer/core"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"
)

const signTypedDataComponent = "use-cases.eth.sign-typed-data"

type signTypedDataUseCase struct {
	vaultClient store.Vault
	logger      *log.Logger
}

func NewSignTypedDataUseCase(vaultClient store.Vault) usecases.SignTypedDataUseCase {
	return &signTypedDataUseCase{
		vaultClient: vaultClient,
		logger:      log.NewLogger().SetComponent(signTypedDataComponent),
	}
}

// Execute signs an arbitrary payload using an existing Ethereum account
func (uc *signTypedDataUseCase) Execute(ctx context.Context, address, namespace string, typedData *signer.TypedData) (string, error) {
	logger := uc.logger.WithContext(ctx).WithField("namespace", namespace).WithField("address", address)

	encodedData, err := utils.GetEIP712EncodedData(typedData)
	if err != nil {
		logger.WithError(err).Error("failed to get typed encoded data")
		return "", errors.FromError(err).ExtendComponent(signTypedDataComponent)
	}

	signature, err := uc.vaultClient.ETHSign(address, namespace, encodedData)
	if err != nil {
		logger.WithError(err).Error("failed to sign typed data")
		return "", errors.FromError(err).ExtendComponent(signTypedDataComponent)
	}

	logger.Info("typed data signed successfully")
	return signature, nil
}
