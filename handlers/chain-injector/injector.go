package chaininjector

import (
	"fmt"

	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// ChainInjector enrich the envelope with the chainUUID, chainName and inject in the input.Context the proxy URL
func ChainInjector(r registry.Client, chainRegistryURL string) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Check if chain exist
		if txctx.Builder.GetChainName() == "" {
			_ = txctx.AbortWithError(errors.DataError("no chainName found")).ExtendComponent(component)
			return
		}

		// retrieve tenantID and check if present
		tenantID := multitenancy.TenantIDFromContext(txctx.Context())
		if tenantID == "" {
			_ = txctx.AbortWithError(errors.DataError("no tenantID found")).ExtendComponent(component)
			return
		}

		if txctx.Builder.GetChainUUID() == "" {
			chain, err := r.GetChainByTenantAndName(txctx.Context(), tenantID, txctx.Builder.GetChainName())
			if err != nil {
				_ = txctx.AbortWithError(errors.FromError(err)).ExtendComponent(component)
				return
			}

			// Inject chain UUID from chain registry
			_ = txctx.Builder.SetChainUUID(chain.UUID)
		}

		// Inject chain proxy path as /tenantID/chainName
		proxyURL := fmt.Sprintf("%s/%s", chainRegistryURL, proxy.PathByChainName(tenantID, txctx.Builder.GetChainName()))
		txctx.WithContext(proxy.With(txctx.Context(), proxyURL))
	}
}

// ChainIDInjector enrich the envelope with the chain UUID retrieved from the chain proxy
func ChainIDInjector(ec ethclient.ChainSyncReader) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.Builder.GetChainID() != nil {
			return
		}

		chainProxyURL := proxy.FromContext(txctx.Context())
		chainID, err := ec.Network(txctx.Context(), chainProxyURL)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("injector: could not retrieve chain id from %s", chainProxyURL)
			return
		}
		_ = txctx.Builder.SetChainID(chainID)
		txctx.Logger.Debugf("injector: chain id %s injected", chainID.String())
	}
}
