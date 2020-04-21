package usecases

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

type ImportChain interface {
	Execute(ctx context.Context, chainEncodeJSON string) error
}

// RegisterContract is a use case to register a new contract
type importChain struct {
	chainAgent store.ChainAgent
	ethClient  ethclient.ChainLedgerReader
}

// NewGetCatalog creates a new GetCatalog
func NewImportChain(chainAgent store.ChainAgent, ec ethclient.ChainLedgerReader) ImportChain {
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

	// In case of not staring block, we use latest
	if chain.ListenerStartingBlock == nil {
		var head uint64
		head, err = utils.GetChainTip(ctx, uc.ethClient, chain.URLs)
		if err != nil {
			logger.WithError(err).Errorf("could not import chain head block. Default 0")
			head = 0
		}

		chain.ListenerStartingBlock = &head
	}

	chain.SetDefault()
	err = uc.chainAgent.RegisterChain(ctx, chain)
	if err != nil {
		return err
	}

	logger.WithFields(logrus.Fields{
		"chain.name":   chain.Name,
		"chain.uuid":   chain.UUID,
		"chain.tenant": chain.TenantID,
	}).Infof("imported chain from configuration")

	return nil
}
