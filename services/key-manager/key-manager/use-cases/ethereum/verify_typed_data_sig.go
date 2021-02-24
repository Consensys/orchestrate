package ethereum

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/log"
	"github.com/ConsenSys/orchestrate/services/key-manager/key-manager/use-cases/ethereum/utils"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	usecases "github.com/ConsenSys/orchestrate/services/key-manager/key-manager/use-cases"
	signer "github.com/ethereum/go-ethereum/signer/core"
)

const verifyTypedDataSignatureComponent = "use-cases.eth.verify-typed-data-signature"

// verifyTypedDataSignatureUseCase is a use case to verify the signature of a typed payload using an existing Ethereum account
type verifyTypedDataSignatureUseCase struct {
	verifySignatureUC usecases.VerifyETHSignatureUseCase
	logger            *log.Logger
}

// NewVerifyTypedDataSignatureUseCase creates a new VerifyTypedDataSignatureUseCase
func NewVerifyTypedDataSignatureUseCase(verifySignatureUC usecases.VerifyETHSignatureUseCase) usecases.VerifyTypedDataSignatureUseCase {
	return &verifyTypedDataSignatureUseCase{
		verifySignatureUC: verifySignatureUC,
		logger:            log.NewLogger().SetComponent(verifyTypedDataSignatureComponent),
	}
}

// Execute verifies the signature of a typed payload using an existing Ethereum account
func (uc *verifyTypedDataSignatureUseCase) Execute(ctx context.Context, address, signature string, typedData *signer.TypedData) error {
	logger := uc.logger.WithContext(ctx).WithField("address", address).WithField("signature", signature)

	encodedData, err := utils.GetEIP712EncodedData(typedData)
	if err != nil {
		logger.WithError(err).Error("failed to get typed encoded data")
		return errors.FromError(err).ExtendComponent(verifyTypedDataSignatureComponent)
	}

	return uc.verifySignatureUC.Execute(ctx, address, signature, encodedData)
}
