package ethereum

import (
	"context"
	"fmt"

	"github.com/consensys/quorum/crypto"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases"
)

const verifySignatureComponent = "use-cases.verify-signature"

// verifySignatureUseCase is a use case to verify the signature of a payload using an existing Ethereum account
type verifySignatureUseCase struct{}

// NewVerifySignatureUseCase creates a new VerifySignatureUseCase
func NewVerifySignatureUseCase() usecases.VerifySignatureUseCase {
	return &verifySignatureUseCase{}
}

// Execute verifies the signature of a payload using an existing Ethereum account
func (uc *verifySignatureUseCase) Execute(ctx context.Context, address, signature, payload string) error {
	logger := log.WithContext(ctx).
		WithField("address", address).
		WithField("signature", signature)
	logger.Debug("verifying signature")

	signatureBytes, err := hexutil.Decode(signature)
	if err != nil {
		errMessage := "failed to decode signature"
		logger.WithError(err).Error(errMessage)
		return errors.InvalidParameterError(fmt.Sprintf("%s: %s", errMessage, err.Error())).ExtendComponent(verifySignatureComponent)
	}

	hash := crypto.Keccak256([]byte(payload))
	pubKey, err := crypto.SigToPub(hash, signatureBytes)
	if err != nil {
		errMessage := "failed to recover public key"
		logger.WithError(err).Error(errMessage)
		return errors.InvalidParameterError(fmt.Sprintf("%s: %s", errMessage, err.Error())).ExtendComponent(verifySignatureComponent)
	}

	recoveredAddress := crypto.PubkeyToAddress(*pubKey)
	if address != recoveredAddress.Hex() {
		errMessage := "failed to verify signature: recovered address does not match the expected one or payload is malformed"
		logger.WithField("recovered_address", recoveredAddress.Hex()).Error(errMessage)
		return errors.InvalidParameterError(errMessage).ExtendComponent(verifySignatureComponent)
	}

	logger.Info("signature verified successfully")
	return nil
}
