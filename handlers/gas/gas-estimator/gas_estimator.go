package gasestimator

import (
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"

	"github.com/ethereum/go-ethereum"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

// EnvelopeToCallMsg enrich an ethereum.CallMsg with Builder information
func EnvelopeToCallMsg(b *tx.Builder, call *ethereum.CallMsg) {
	call.To = b.GetTo()
	call.From = b.MustGetFromAddress()
	call.Value = b.GetValue()
	call.Data = b.MustGetDataBytes()
}

// Estimator creates an handler that set Gas Limit
func Estimator(p ethclient.GasEstimator) engine.HandlerFunc {
	pool := &sync.Pool{
		New: func() interface{} { return &ethereum.CallMsg{} },
	}

	return func(txctx *engine.TxContext) {

		if txctx.Builder.Gas == nil {
			// Retrieve re-cycled CallMsg
			call := pool.Get().(*ethereum.CallMsg)
			defer pool.Put(call)

			// Estimate gas
			EnvelopeToCallMsg(txctx.Builder, call)

			url, err := proxy.GetURL(txctx)
			if err != nil {
				return
			}

			g, err := p.EstimateGas(txctx.Context(), url, call)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("gas-estimator: could not estimate gas limit")
				return
			}

			// Set gas limit on context
			_ = txctx.Builder.SetGas(g)

			// Enrich logger
			txctx.Logger = txctx.Logger.WithFields(log.Fields{
				"tx.gas": g,
			})
			txctx.Logger.Debugf("gas-estimator: gas limit set")
		}

		// Enrich logger
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.gas": txctx.Builder.GetGas(),
		})
	}
}
