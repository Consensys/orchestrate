package ethereum

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/crypto/ethereum"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/log"
	"github.com/ConsenSys/orchestrate/pkg/utils"
	usecases "github.com/ConsenSys/orchestrate/services/key-manager/key-manager/use-cases"
)

const verifySignatureComponent = "use-cases.eth.verify-signature"

// verifySignatureUseCase is a use case to verify the signature of a payload using an existing Ethereum account
type verifySignatureUseCase struct {
	logger *log.Logger
}

// NewVerifySignatureUseCase creates a new VerifyETHSignatureUseCase
func NewVerifySignatureUseCase() usecases.VerifyETHSignatureUseCase {
	return &verifySignatureUseCase{
		logger: log.NewLogger().SetComponent(verifySignatureComponent),
	}
}

// Execute verifies the signature of a payload using an existing Ethereum account
func (uc *verifySignatureUseCase) Execute(ctx context.Context, address, signature, payload string) error {
	logger := uc.logger.WithContext(ctx).WithField("address", address).
		WithField("signature", utils.ShortString(signature, 10))

	recoveredAddress, err := ethereum.GetSignatureSender(signature, payload)
	if err != nil {
		logger.WithError(err).Error("failed to signature extract sender")
		return errors.InvalidParameterError(err.Error()).ExtendComponent(verifySignatureComponent)
	}

	if address != recoveredAddress.Hex() {
		errMessage := "failed to verify signature: recovered address does not match the expected one or payload is malformed"
		logger.WithField("recovered_address", recoveredAddress.Hex()).Error(errMessage)
		return errors.InvalidParameterError(errMessage).ExtendComponent(verifySignatureComponent)
	}

	logger.Info("signature verified successfully")
	return nil
}
