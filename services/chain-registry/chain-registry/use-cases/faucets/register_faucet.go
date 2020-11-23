package faucets

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
)

type RegisterFaucet interface {
	Execute(ctx context.Context, faucet *models.Faucet) error
}

// RegisterContract is a use case to register a new contract
type registerFaucet struct {
	faucetAgent store.FaucetAgent
}

// NewGetCatalog creates a new GetCatalog
func NewRegisterFaucet(faucetAgent store.FaucetAgent) RegisterFaucet {
	return &registerFaucet{
		faucetAgent: faucetAgent,
	}
}

func (uc *registerFaucet) Execute(ctx context.Context, faucet *models.Faucet) error {
	logger := log.FromContext(ctx)

	faucet.SetDefault()
	err := uc.faucetAgent.RegisterFaucet(ctx, faucet)
	if err != nil {
		return err
	}

	logger.WithFields(logrus.Fields{
		"faucet.name":   faucet.Name,
		"faucet.uuid":   faucet.UUID,
		"faucet.tenant": faucet.TenantID,
	}).Infof("registered a new faucet")

	return nil
}
