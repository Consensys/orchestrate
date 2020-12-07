package ethereum

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases/ethereum/utils"

	signer "github.com/ethereum/go-ethereum/signer/core"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases"
)

const verifyTypedDataSignatureComponent = "use-cases.verify-typed-data-signature"

// verifyTypedDataSignatureUseCase is a use case to verify the signature of a typed payload using an existing Ethereum account
type verifyTypedDataSignatureUseCase struct {
	verifySignatureUC usecases.VerifySignatureUseCase
}

// NewVerifyTypedDataSignatureUseCase creates a new VerifyTypedDataSignatureUseCase
func NewVerifyTypedDataSignatureUseCase(verifySignatureUC usecases.VerifySignatureUseCase) usecases.VerifyTypedDataSignatureUseCase {
	return &verifyTypedDataSignatureUseCase{verifySignatureUC: verifySignatureUC}
}

// Execute verifies the signature of a typed payload using an existing Ethereum account
func (uc *verifyTypedDataSignatureUseCase) Execute(ctx context.Context, address, signature string, typedData *signer.TypedData) error {
	logger := log.WithContext(ctx).
		WithField("address", address).
		WithField("signature", signature)
	logger.Debug("verifying typed data signature")

	encodedData, err := utils.GetEIP712EncodedData(typedData)
	if err != nil {
		return errors.FromError(err).ExtendComponent(verifyTypedDataSignatureComponent)
	}

	return uc.verifySignatureUC.Execute(ctx, address, signature, encodedData)
}
