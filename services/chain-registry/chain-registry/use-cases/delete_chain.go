package usecases

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

type DeleteChain interface {
	Execute(ctx context.Context, uuid, tenantID string) error
}

// RegisterContract is a use case to register a new contract
type deleteChain struct {
	chainAgent store.ChainAgent
}

// NewGetCatalog creates a new GetCatalog
func NewDeleteChain(chainAgent store.ChainAgent) DeleteChain {
	return &deleteChain{
		chainAgent: chainAgent,
	}
}

func (uc *deleteChain) Execute(ctx context.Context, uuid, tenantID string) error {
	logger := log.FromContext(ctx)

	var err error
	if tenantID == "" || tenantID == multitenancy.DefaultTenantIDName {
		err = uc.chainAgent.DeleteChainByUUID(ctx, uuid)
	} else {
		err = uc.chainAgent.DeleteChainByUUIDAndTenant(ctx, uuid, tenantID)
	}

	logger.WithFields(logrus.Fields{
		"chainUUID": uuid,
		"tenantID":  tenantID,
	}).Infof("deleted chain")

	return err
}
