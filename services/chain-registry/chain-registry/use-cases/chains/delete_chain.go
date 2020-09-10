package chains

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

type DeleteChain interface {
	Execute(ctx context.Context, uuid string, tenants []string) error
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

func (uc *deleteChain) Execute(ctx context.Context, uuid string, tenants []string) error {
	logger := log.FromContext(ctx)

	err := uc.chainAgent.DeleteChain(ctx, uuid, tenants)

	logger.WithFields(logrus.Fields{
		"chainUUID": uuid,
		"tenantIDs": tenants,
	}).Infof("deleted chain")

	return err
}
