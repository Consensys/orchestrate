package usecases

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

type DeleteFaucet interface {
	Execute(ctx context.Context, uuid string, tenants []string) error
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

func (uc *deleteFaucet) Execute(ctx context.Context, uuid string, tenants []string) error {
	logger := log.FromContext(ctx)

	err := uc.faucetAgent.DeleteFaucet(ctx, uuid, tenants)

	logger.WithFields(logrus.Fields{
		"uuid":    uuid,
		"tenants": tenants,
	}).Infof("deleted faucet")

	return err
}
