package gasestimator

import (
	"sync"

	"github.com/ethereum/go-ethereum"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

// EnvelopeToCallMsg enrich an ethereum.CallMsg with Envelope information
func EnvelopeToCallMsg(e *envelope.Envelope, call *ethereum.CallMsg) {
	to, _ := e.GetTx().GetTxData().ToAddress()
	call.To = &to
	call.From, _ = e.GetSender().Address()
	call.Value = e.GetTx().GetTxData().ValueBig()
	call.Data = e.GetTx().GetTxData().DataBytes()
}

// Estimator creates an handler that set Gas Limit
func Estimator(p services.GasEstimator) engine.HandlerFunc {
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
			g, err := p.EstimateGas(txctx.Context(), txctx.Envelope.GetChain().ID(), *call)
			if err != nil {
				// TODO: handle error
				txctx.Logger.WithError(err).Errorf("gas-estimator: could not estimate gas limit")
				_ = txctx.AbortWithError(err)
			} else {
				// Set gas limit on context
				txctx.Envelope.GetTx().GetTxData().SetGas(g)

				// Enrich logger
				txctx.Logger = txctx.Logger.WithFields(log.Fields{
					"tx.gas": g,
				})
				txctx.Logger.Debugf("gas-estimator: gas limit set")
			}
		} else {
			// Enrich logger
			txctx.Logger = txctx.Logger.WithFields(log.Fields{
				"tx.gas": txctx.Envelope.GetTx().GetTxData().GetGas(),
			})
		}
	}
}
