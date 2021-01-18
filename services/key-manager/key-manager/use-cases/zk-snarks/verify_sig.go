package zksnarks

import (
	"context"

	log "github.com/sirupsen/logrus"
	zksnarks "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/crypto/zk-snarks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases"
)

const verifySignatureComponent = "use-cases.zks.verify-signature"

type verifySignatureUseCase struct{}

func NewVerifySignatureUseCase() usecases.VerifyZKSSignatureUseCase {
	return &verifySignatureUseCase{}
}

func (uc *verifySignatureUseCase) Execute(ctx context.Context, publicKey, signature, payload string) error {
	logger := log.WithContext(ctx).
		WithField("public_key", publicKey).
		WithField("signature", utils.ShortString(signature, 10))
	logger.Debug("verifying signature")

	verified, err := zksnarks.VerifyZKSMessage(publicKey, signature, []byte(payload))
	if err != nil || !verified {
		errMessage := "failed to verify signature: publicKey does not match the expected one or payload is malformed"
		logger.WithError(err).Error(errMessage)
		return errors.InvalidParameterError(errMessage).ExtendComponent(verifySignatureComponent)
	}

	logger.Info("signature verified successfully")
	return nil
}
