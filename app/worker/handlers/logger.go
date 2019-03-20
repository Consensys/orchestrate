package handlers

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// Logger to log context elements before and after the worker
func Logger(ctx *worker.Context) {

	ctx.Logger = ctx.Logger.WithFields(log.Fields{
		"chain.id":        ctx.T.GetChain().GetId(),
		"receipt.txhash":  ctx.T.GetReceipt().GetTxHash(),
		"receipt.txIndex": ctx.T.GetReceipt().GetTxIndex(),
	})

	ctx.Logger.Debug("worker: new receipt")
	start := time.Now()

	ctx.Next()

	latency := time.Now().Sub(start)

	ctx.Logger.WithFields(log.Fields{
		"latency": latency,
	}).WithError(fmt.Errorf("%v", ctx.T.Errors)).Info("worker: receipt processed")
}
