package usecases

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

const registerChainComponent = "use-cases.register-chain"

type RegisterChain interface {
	Execute(ctx context.Context, chain *models.Chain) error
}

// RegisterContract is a use case to register a new contract
type registerChain struct {
	chainAgent store.ChainAgent
	ethClient  ethclient.Client
}

// NewGetCatalog creates a new GetCatalog
func NewRegisterChain(chainAgent store.ChainAgent, ec ethclient.Client) RegisterChain {
	return &registerChain{
		chainAgent: chainAgent,
		ethClient:  ec,
	}
}

func (uc *registerChain) Execute(ctx context.Context, chain *models.Chain) error {
	logger := log.WithContext(ctx).WithField("chain_name", chain.Name)
	logger.Debug("registering new chain")

	// Verifies URLs are valid and get chain ID
	err := utils.VerifyURLs(ctx, uc.ethClient, chain.URLs)
	if err != nil {
		return errors.FromError(err).ExtendComponent(registerChainComponent)
	}

	chainID, err := utils.GetChainID(ctx, uc.ethClient, chain.URLs)
	if err != nil {
		return errors.FromError(err).ExtendComponent(registerChainComponent)
	}
	chain.ChainID = chainID.String()

	// If no starting block provided, we use the latest
	if chain.ListenerStartingBlock == nil {
		head := utils.GetChainTip(ctx, uc.ethClient, chain.URLs)
		chain.ListenerStartingBlock = &head
	}

	chain.SetDefault()
	err = uc.chainAgent.RegisterChain(ctx, chain)
	if err != nil {
		return errors.FromError(err).ExtendComponent(registerChainComponent)
	}

	logger.WithField("chain_uuid", chain.UUID).Info("chain successfully registered")

	return nil
}
