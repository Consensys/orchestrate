package gaspricer

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
)

// Pricer creates a handler that set a Gas Price
func Pricer(p ethclient.GasPricer) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.Envelope.GasPrice == nil {
			url, err := proxy.GetURL(txctx)
			if err != nil {
				return
			}

			// Envelope a gas price suggestion
			p, err := p.SuggestGasPrice(txctx.Context(), url)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("gas-pricer: could not suggest gas price")
				return
			}

			// Set gas price
			_ = txctx.Envelope.SetGasPrice(p)
			txctx.Logger.Debugf("gas-pricer: gas price set")
		}

		// Enrich logger
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"gasPrice": txctx.Envelope.GetGasPriceString(),
		})
	}
}
