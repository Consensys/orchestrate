package chaininjector

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
)

func ChainUUIDHandler(r registry.ChainRegistryClient, chainRegistryURL string) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		err := chainUUIDInjector(txctx, r, chainRegistryURL)
		if err != nil {
			_ = txctx.AbortWithError(err).ExtendComponent(component)
			return
		}
	}
}

func ChainUUIDHandlerWithoutAbort(r registry.ChainRegistryClient, chainRegistryURL string) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		err := chainUUIDInjector(txctx, r, chainRegistryURL)
		if err != nil {
			_ = txctx.Error(err).ExtendComponent(component)
			return
		}
	}
}

func chainUUIDInjector(txctx *engine.TxContext, r registry.ChainRegistryClient, chainRegistryURL string) error {
	// Check if chain exist
	if txctx.Envelope.GetChainName() == "" {
		return errors.DataError("no chain name found")
	}

	if txctx.Envelope.GetChainUUID() == "" {
		chain, err := r.GetChainByName(txctx.Context(), txctx.Envelope.GetChainName())
		if err != nil {
			return errors.FromError(err)
		}

		// chainUUIDInjector chain UUID from chain registry
		_ = txctx.Envelope.SetChainUUID(chain.UUID)
	}

	// chainUUIDInjector chain proxy path as /chainUUID
	proxyURL := fmt.Sprintf("%s/%s", chainRegistryURL, txctx.Envelope.GetChainUUID())
	txctx.WithContext(proxy.With(txctx.Context(), proxyURL))
	return nil
}

// ChainIDInjectorHandler enrich the envelope with the chain UUID retrieved from the chain proxy
func ChainIDInjectorHandler(ec ethclient.ChainSyncReader) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.Envelope.GetChainID() != nil {
			return
		}

		chainProxyURL := proxy.FromContext(txctx.Context())
		chainID, err := ec.Network(txctx.Context(), chainProxyURL)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("injector: could not retrieve chain id from %s", chainProxyURL)
			return
		}
		_ = txctx.Envelope.SetChainID(chainID)
		txctx.Logger.Debugf("injector: chain id %s injected", chainID.String())
	}
}
