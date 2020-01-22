package gaspricer

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

// Pricer creates a handler that set a Gas Price
func Pricer(p ethclient.GasPricer) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.Envelope.GetTx().GetTxData().GetGasPrice() == nil {
			url, err := proxy.GetURL(txctx)
			if err != nil {
				return
			}

			// Request a gas price suggestion
			p, err := p.SuggestGasPrice(txctx.Context(), url)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("gas-pricer: could not suggest gas price")
				return
			}

			// Set gas price
			txctx.Envelope.Tx.TxData.SetGasPrice(p)
			txctx.Logger.Debugf("gas-pricer: gas price set")

			return
		}

		// Enrich logger
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.gas.price": txctx.Envelope.GetTx().GetTxData().GetGasPrice(),
		})
	}
}
