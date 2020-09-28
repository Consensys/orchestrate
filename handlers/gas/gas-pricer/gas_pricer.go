package gaspricer

import (
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
)

// Pricer creates a handler that set a Gas Price
func Pricer(p ethclient.GasPricer) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger.WithField("envelope_id", txctx.Envelope.GetID()).Debugf("pricer handler starts")
		if txctx.Envelope.GasPrice == nil {
			url, err := proxy.GetURL(txctx)
			if err != nil {
				return
			}

			p, err := p.SuggestGasPrice(txctx.Context(), url)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("gas-pricer: could not suggest gas price")
				return
			}

			// Set gas price
			_ = txctx.Envelope.SetGasPrice(applyPriorityCoefficient(p, txctx.Envelope.GetContextLabelsValue(tx.PriorityLabel)))
			txctx.Logger.Debugf("gas-pricer: gas price set")
		}

		// Enrich logger
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"gasPrice": txctx.Envelope.GetGasPriceString(),
		})
	}
}

func applyPriorityCoefficient(initialPrice *big.Int, priority string) *big.Int {
	switch priority {
	case utils.PriorityVeryLow:
		return initialPrice.Mul(initialPrice, big.NewInt(6)).Div(initialPrice, big.NewInt(10))
	case utils.PriorityLow:
		return initialPrice.Mul(initialPrice, big.NewInt(8)).Div(initialPrice, big.NewInt(10))
	case utils.PriorityMedium:
		return initialPrice
	case utils.PriorityHigh:
		return initialPrice.Mul(initialPrice, big.NewInt(12)).Div(initialPrice, big.NewInt(10))
	case utils.PriorityVeryHigh:
		return initialPrice.Mul(initialPrice, big.NewInt(14)).Div(initialPrice, big.NewInt(10))
	default:
		return initialPrice
	}
}
