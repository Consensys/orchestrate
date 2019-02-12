package handlers

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Logger to log context elements before and after the worker
func Logger(ctx *types.Context) {
	log.WithFields(log.Fields{
		"Chain":       ctx.T.Chain().ID.Text(16),
		"BlockNumber": ctx.T.Receipt().BlockNumber,
		"TxIndex":     ctx.T.Receipt().TxIndex,
		"TxHash":      ctx.T.Receipt().TxHash.Hex(),
	}).Debug("tx-listener-worker: new receipt")

	ctx.Next()

	if len(ctx.T.Errors) != 0 {
		log.WithFields(log.Fields{
			"Chain":  ctx.T.Chain().ID.Text(16),
			"TxHash": ctx.T.Receipt().TxHash.Hex(),
		}).Errorf("tx-listener-worker: Errors: %v", ctx.T.Errors)
	}
}
