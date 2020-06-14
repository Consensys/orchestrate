package chains

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

//go:generate mockgen -source=get_chain_by_name.go -destination=mocks/get_chain_by_name.go -package=mocks

const getChainByNameComponent = "use-cases.get-chain-by-name"

type GetChainByNameUseCase interface {
	Execute(ctx context.Context, chainName, tenantID string) (*types.Chain, error)
}

// GetChainByNameUseCase is a use case to get a job
type getChainByNameUseCase struct {
	chainRegistryClient client.ChainRegistryClient
}

// NewGetChainByNameUseCase creates a new GetChainByNameUseCase
func NewGetChainByNameUseCase(chainRegistryClient client.ChainRegistryClient) GetChainByNameUseCase {
	return &getChainByNameUseCase{
		chainRegistryClient: chainRegistryClient,
	}
}

// Execute gets a job
func (uc *getChainByNameUseCase) Execute(ctx context.Context, chainName, tenantID string) (*types.Chain, error) {
	log.WithContext(ctx).
		WithField("chain_name", chainName).
		Debug("getting chain")

	// Validate that the chain exists
	chain, err := uc.chainRegistryClient.GetChainByName(ctx, chainName)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("cannot load '%s' chain", chainName)
		if errors.IsNotFoundError(err) {
			return nil, errors.InvalidParameterError(err.Error()).ExtendComponent(getChainByNameComponent)
		}
		return nil, errors.FromError(err).ExtendComponent(getChainByNameComponent)
	}

	// TODO: validate tenantID once #180 is done ()

	return parsers.NewChainFromModels(chain), nil
}
