package handlers

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Logger to log context elements before and after the worker
func Logger(ctx *types.Context) {

	ctx.Logger = log.WithFields(log.Fields{
		"eth.chain":       ctx.T.Chain().ID.Text(16),
		"eth.blockNumber": ctx.T.Receipt().BlockNumber,
		"eth.txIndex":     ctx.T.Receipt().TxIndex,
		"eth.txHash":      ctx.T.Receipt().TxHash.Hex(),
	})

	ctx.Logger.Debug("worker: new receipt")
	start := time.Now()

	ctx.Next()

	latency := time.Now().Sub(start)

	if len(ctx.T.Errors) != 0 {
		ctx.Logger.WithFields(log.Fields{
			"latency": latency,
		}).Errorf("worker: Errors: %v", ctx.T.Errors)
	} else {
		ctx.Logger.WithFields(log.Fields{
			"latency": latency,
		}).Info("worker: message processed")
	}
}
