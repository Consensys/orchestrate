package chaininjector

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient"
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
			txctx.Logger.Warn(err)
		}
	}
}

func chainUUIDInjector(txctx *engine.TxContext, r registry.ChainRegistryClient, chainRegistryURL string) error {
	chainUUID := txctx.Envelope.GetChainUUID()
	chainName := txctx.Envelope.GetChainName()

	if chainUUID == "" && chainName == "" {
		return errors.DataError("no chain found")
	}

	if chainUUID == "" {
		chain, err := r.GetChainByName(txctx.Context(), chainName)
		if err != nil {
			return errors.FromError(err)
		}

		chainUUID = chain.UUID
		_ = txctx.Envelope.SetChainUUID(chainUUID)
	} else {
		_ = txctx.Envelope.SetChainUUID(chainUUID)
	}

	// chainUUIDInjector chain proxy path as /chainUUID
	proxyURL := fmt.Sprintf("%s/%s", chainRegistryURL, chainUUID)
	txctx.WithContext(proxy.With(txctx.Context(), proxyURL))
	return nil
}

// ChainIDInjectorHandler enrich the envelope with the chain UUID retrieved from the chain proxy
func ChainIDInjectorHandler(ec ethclient.ChainSyncReader) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger.WithField("envelope_id", txctx.Envelope.GetID()).Debugf("chainID injector handler starts")
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
