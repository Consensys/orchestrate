package chaininjector

import (
	"fmt"

	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// ChainInjector enrich the envelope with the chainUUID, chainName and inject in the input.Context the proxy URL
func ChainInjector(multitenancyEnabled bool, r registry.Client, chainRegistryURL string) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// check chainUUID and inject if not present
		tenantID := multitenancy.DefaultTenantIDName
		if multitenancyEnabled {
			if tenantID = multitenancy.TenantIDFromContext(txctx.Context()); tenantID == "" {
				e := txctx.AbortWithError(errors.InternalError("invalid tenantID not found")).ExtendComponent(component)
				txctx.Logger.Error(e)
				return
			}
		}

		// Check if chain exist
		chain, err := getChain(txctx, r, tenantID)
		if err != nil {
			e := txctx.AbortWithError(errors.FromError(err)).ExtendComponent(component)
			txctx.Logger.Error(e)
			return
		}

		// Re-inject chain Name and chain UUID from chain registry
		txctx.Envelope.GetChain().SetName(chain.Name)
		txctx.Envelope.GetChain().SetUUID(chain.UUID)

		// Inject chain proxy path as /tenantID/chainName
		proxyURL := fmt.Sprintf("%s/%s", chainRegistryURL, proxy.PathByChainName(tenantID, chain.Name))
		txctx.WithContext(proxy.With(txctx.Context(), proxyURL))
	}
}

// getChain retrieves chain in the chain registry
func getChain(txctx *engine.TxContext, r registry.Client, tenantID string) (*types.Chain, error) {
	var n *types.Chain
	var err error
	chainUUID := txctx.Envelope.GetChain().GetUuid()
	chainName := txctx.Envelope.GetChain().GetName()
	switch {
	case chainUUID != "":
		n, err = r.GetChainByTenantAndUUID(txctx.Context(), tenantID, chainUUID)
	case chainName != "":
		n, err = r.GetChainByTenantAndName(txctx.Context(), tenantID, chainName)
	default:
		return nil, errors.InternalError("invalid envelope - no chain uuid or chain name are filled - cannot retrieve chain data")
	}
	if err != nil {
		return nil, err
	}
	return n, nil
}

// ChainIDInjector enrich the envelope with the chain UUID retrieved from the chain proxy
func ChainIDInjector(ec ethclient.ChainSyncReader) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		chainProxyURL := proxy.FromContext(txctx.Context())
		chainID, err := ec.Network(txctx.Context(), chainProxyURL)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("injector: could not retrieve chain id from %s", chainProxyURL)
			return
		}
		txctx.Envelope.GetChain().SetChainID(chainID)
		txctx.Logger.Debugf("injector: chain id %s injected", chainID.String())
	}
}
