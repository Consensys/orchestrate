package account

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/identity-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
)

const signPayloadComponent = "use-cases.sign-payload"

type signPayloadUseCase struct {
	keyManagerClient client.KeyManagerClient
}

func NewSignPayloadUseCase(keyManagerClient client.KeyManagerClient) usecases.SignPayloadUseCase {
	return &signPayloadUseCase{
		keyManagerClient: keyManagerClient,
	}
}

func (uc *signPayloadUseCase) Execute(ctx context.Context, address, payload, tenantID string) (string, error) {
	log.WithContext(ctx).WithField("address", address).
		Debug("signing payload")

	signature, err := uc.keyManagerClient.ETHSign(ctx, address, &keymanager.PayloadRequest{
		Data:      payload,
		Namespace: tenantID,
	})
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(signPayloadComponent)
	}

	log.WithContext(ctx).Debug("payload signed successfully")
	return signature, nil
}
