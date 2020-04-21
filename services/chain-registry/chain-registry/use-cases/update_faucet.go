package usecases

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

type UpdateFaucet interface {
	Execute(ctx context.Context, uuid string, faucet *models.Faucet) error
}

// RegisterContract is a use case to register a new contract
type updateFaucet struct {
	faucetAgent store.FaucetAgent
}

// NewGetCatalog creates a new GetCatalog
func NewUpdateFaucet(faucetAgent store.FaucetAgent) UpdateFaucet {
	return &updateFaucet{
		faucetAgent: faucetAgent,
	}
}

func (uc *updateFaucet) Execute(ctx context.Context, uuid string, faucet *models.Faucet) error {
	logger := log.FromContext(ctx)

	if faucet.UUID != "" && faucet.UUID != uuid {
		return errors.ConstraintViolatedError("update faucet UUID is not allowed")
	}

	err := uc.faucetAgent.UpdateFaucetByUUID(ctx, uuid, faucet)
	if err != nil {
		return err
	}

	logger.WithFields(logrus.Fields{
		"faucet.name":   faucet.Name,
		"faucet.uuid":   faucet.UUID,
		"faucet.tenant": faucet.TenantID,
	}).Infof("updated faucet")

	return nil
}
