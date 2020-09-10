package chains

import (
	"context"
	"encoding/json"
	"strings"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

const importChainComponent = "use-cases.import-chain"

type ImportChain interface {
	Execute(ctx context.Context, chainEncodeJSON string) error
}

// RegisterContract is a use case to register a new contract
type importChain struct {
	chainAgent store.ChainAgent
	ethClient  ethclient.Client
}

// NewGetCatalog creates a new GetCatalog
func NewImportChain(chainAgent store.ChainAgent, ec ethclient.Client) ImportChain {
	return &importChain{
		chainAgent: chainAgent,
		ethClient:  ec,
	}
}

func (uc *importChain) Execute(ctx context.Context, chainEncodeJSON string) error {
	logger := log.FromContext(ctx)
	logger.WithField("config", chainEncodeJSON).Debugf("import chain from configuration")
	chain := &models.Chain{}
	dec := json.NewDecoder(strings.NewReader(chainEncodeJSON))
	dec.DisallowUnknownFields() // Force errors if unknown fields
	err := dec.Decode(chain)
	if err != nil {
		return err
	}

	// Verifies URLs are valid and get chain ID
	err = utils.VerifyURLs(ctx, uc.ethClient, chain.URLs)
	if err != nil {
		return errors.FromError(err).ExtendComponent(importChainComponent)
	}

	chainID, err := utils.GetChainID(ctx, uc.ethClient, chain.URLs)
	if err != nil {
		return errors.FromError(err).ExtendComponent(importChainComponent)
	}
	chain.ChainID = chainID.String()

	// In case of not staring block, we use latest
	if chain.ListenerStartingBlock == nil {
		head := utils.GetChainTip(ctx, uc.ethClient, chain.URLs)
		chain.ListenerStartingBlock = &head
	}

	chain.SetDefault()
	err = uc.chainAgent.RegisterChain(ctx, chain)
	if err != nil {
		return errors.FromError(err).ExtendComponent(importChainComponent)
	}

	logger.WithFields(logrus.Fields{
		"chainName": chain.Name,
		"chainUUID": chain.UUID,
		"tenantID":  chain.TenantID,
	}).Infof("imported chain from configuration")

	return nil
}
