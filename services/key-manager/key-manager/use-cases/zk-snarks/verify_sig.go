package zksnarks

import (
	"context"

	zksnarks "github.com/ConsenSys/orchestrate/pkg/crypto/zk-snarks"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/log"
	"github.com/ConsenSys/orchestrate/pkg/utils"
	usecases "github.com/ConsenSys/orchestrate/services/key-manager/key-manager/use-cases"
)

const verifySignatureComponent = "use-cases.zks.verify-signature"

type verifySignatureUseCase struct {
	logger *log.Logger
}

func NewVerifySignatureUseCase() usecases.VerifyZKSSignatureUseCase {
	return &verifySignatureUseCase{
		logger: log.NewLogger().SetComponent(verifySignatureComponent),
	}
}

func (uc *verifySignatureUseCase) Execute(ctx context.Context, publicKey, signature, payload string) error {
	logger := uc.logger.WithContext(ctx).
		WithField("component", verifySignatureComponent).
		WithField("public_key", publicKey).
		WithField("signature", utils.ShortString(signature, 10))

	verified, err := zksnarks.VerifyZKSMessage(publicKey, signature, []byte(payload))
	if err != nil || !verified {
		errMessage := "failed to verify signature"
		logger.WithError(err).Error(errMessage)
		return errors.InvalidParameterError(errMessage).ExtendComponent(verifySignatureComponent)
	}

	logger.Info("signature verified successfully")
	return nil
}
