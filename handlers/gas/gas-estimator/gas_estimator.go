package gasestimator

import (
	"sync"

	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
)

// EnvelopeToCallMsg enrich an ethereum.CallMsg with Envelope information
func EnvelopeToCallMsg(b *tx.Envelope, call *ethereum.CallMsg) {
	call.To = b.GetTo()
	if b.IsOneTimeKeySignature() {
		// Generate a dummy eth address to enforce estimation
		call.From = ethcommon.HexToAddress("0x1")
	} else {
		call.From = b.MustGetFromAddress()
	}
	call.Value = b.GetValue()
	call.Data = b.MustGetDataBytes()
}

// Estimator creates an handler that set Gas Limit
func Estimator(p ethclient.GasEstimator) engine.HandlerFunc {
	pool := &sync.Pool{
		New: func() interface{} { return &ethereum.CallMsg{} },
	}

	return func(txctx *engine.TxContext) {

		if txctx.Envelope.GetGas() == nil {
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
			_ = txctx.Envelope.SetGas(g)

			// Enrich logger
			txctx.Logger = txctx.Logger.WithFields(log.Fields{
				"tx.gas": g,
			})
			txctx.Logger.Debugf("gas-estimator: gas limit set")
		}

		// Enrich logger
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"tx.gas": txctx.Envelope.MustGetGasUint64(),
		})
	}
}
