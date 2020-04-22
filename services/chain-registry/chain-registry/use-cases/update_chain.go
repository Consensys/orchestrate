package usecases

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

type UpdateChain interface {
	Execute(ctx context.Context, uuid, chainName string, chain *models.Chain) error
}

// RegisterContract is a use case to register a new contract
type updateChain struct {
	chainAgent store.ChainAgent
}

// NewGetCatalog creates a new GetCatalog
func NewUpdateChain(chainAgent store.ChainAgent) UpdateChain {
	return &updateChain{
		chainAgent: chainAgent,
	}
}

func (uc *updateChain) Execute(ctx context.Context, uuid, chainName string, chain *models.Chain) error {
	logger := log.FromContext(ctx)

	if chain.UUID != "" && chain.UUID != uuid {
		return errors.ConstraintViolatedError("update chain UUID is not allowed")
	}

	// We need to insert UUID for the new PrivateTxManagers if so
	chain.SetPrivateTxManagersDefault()

	var err error
	if uuid != "" {
		err = uc.chainAgent.UpdateChainByUUID(ctx, uuid, chain)
	} else if chainName != "" {
		err = uc.chainAgent.UpdateChainByName(ctx, chainName, chain)
	}

	if err != nil {
		return err
	}

	logger.WithFields(logrus.Fields{
		"chain.name":   chain.Name,
		"chain.uuid":   chain.UUID,
		"chain.tenant": chain.TenantID,
	}).Infof("updated chain from configuration")

	return nil
}
