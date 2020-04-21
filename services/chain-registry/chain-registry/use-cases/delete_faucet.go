package usecases

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

type DeleteFaucet interface {
	Execute(ctx context.Context, uuid, tenantID string) error
}

// RegisterContract is a use case to register a new contract
type deleteFaucet struct {
	faucetAgent store.FaucetAgent
}

// NewGetCatalog creates a new GetCatalog
func NewDeleteFaucet(faucetAgent store.FaucetAgent) DeleteFaucet {
	return &deleteFaucet{
		faucetAgent: faucetAgent,
	}
}

func (uc *deleteFaucet) Execute(ctx context.Context, uuid, tenantID string) error {
	logger := log.FromContext(ctx)

	var err error
	if tenantID == "" || tenantID == multitenancy.DefaultTenantIDName {
		err = uc.faucetAgent.DeleteFaucetByUUID(ctx, uuid)
	} else {
		err = uc.faucetAgent.DeleteFaucetByUUIDAndTenant(ctx, uuid, tenantID)
	}

	logger.WithFields(logrus.Fields{
		"faucet.uuid":   uuid,
		"faucet.tenant": tenantID,
	}).Infof("deleted faucet")

	return err
}
