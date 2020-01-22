package gasestimator

import (
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"

	"github.com/ethereum/go-ethereum"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
)

// EnvelopeToCallMsg enrich an ethereum.CallMsg with Envelope information
func EnvelopeToCallMsg(e *envelope.Envelope, call *ethereum.CallMsg) {
	to := e.GetTx().GetTxData().Receiver()
	call.To = &to
	call.From = e.Sender()
	call.Value = e.GetTx().GetTxData().GetValueBig()
	call.Data = e.GetTx().GetTxData().GetDataBytes()
}

// Estimator creates an handler that set Gas Limit
func Estimator(p ethclient.GasEstimator) engine.HandlerFunc {
	pool := &sync.Pool{
		New: func() interface{} { return &ethereum.CallMsg{} },
	}

	return func(txctx *engine.TxContext) {

		if txctx.Envelope.GetTx().GetTxData().GetGas() == 0 {
			// Retrieve re-cycled CallMsg
			call := pool.Get().(*ethereum.CallMsg)
			defer pool.Put(call)

			// Estimate gas
			EnvelopeToCallMsg(txctx.Envelope, call)

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
			txctx.Envelope.GetTx().GetTxData().SetGas(g)

			// Enrich logger
			txctx.Logger = txctx.Logger.WithFields(log.Fields{
				"tx.gas": g,
			})
			txctx.Logger.Debugf("gas-estimator: gas limit set")
		}

		// Enrich logger
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.gas": txctx.Envelope.GetTx().GetTxData().GetGas(),
		})
	}
}
