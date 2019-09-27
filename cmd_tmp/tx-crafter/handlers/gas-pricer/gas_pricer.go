package gaspricer

import (
	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient"
)

// Pricer creates a handler that set a Gas Price
func Pricer(p ethclient.GasPricer) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.Envelope.GetTx().GetTxData().GetGasPrice() == nil {
			// Request a gas price suggestion
			p, err := p.SuggestGasPrice(txctx.Context(), txctx.Envelope.Chain.ID())
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
