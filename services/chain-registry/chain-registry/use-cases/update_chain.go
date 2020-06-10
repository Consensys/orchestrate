package usecases

import (
	"context"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

type UpdateChain interface {
	Execute(ctx context.Context, uuid, chainName string, tenants []string, chain *models.Chain) error
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

func (uc *updateChain) Execute(ctx context.Context, uuid, chainName string, tenants []string, chain *models.Chain) error {
	logger := log.FromContext(ctx)

	if chain.UUID != "" && chain.UUID != uuid {
		return errors.ConstraintViolatedError("update chain UUID is not allowed")
	}

	// We need to insert UUID for the new PrivateTxManagers if so
	chain.SetPrivateTxManagersDefault()

	// @TODO: This is a provisional hack meanwhile task #PO-479 is implemented
	updateAt := time.Now().UTC()
	chain.UpdatedAt = &updateAt

	var err error
	if uuid != "" {
		err = uc.chainAgent.UpdateChain(ctx, uuid, tenants, chain)
	} else if chainName != "" {
		err = uc.chainAgent.UpdateChainByName(ctx, chainName, tenants, chain)
	}

	if err != nil {
		return err
	}

	logger.WithFields(logrus.Fields{
		"chainName": chain.Name,
		"chainUUID": chain.UUID,
		"tenantIDs": tenants,
	}).Infof("updated chain from configuration")

	return nil
}
