package handlers

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	coreWorker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// Crafter creates a crafter handler
func Crafter(r services.ABIRegistry, c services.Crafter) coreWorker.HandlerFunc {
	return func(ctx *coreWorker.Context) {
		// Retrieve method identifier from trace
		methodID := ctx.T.Call.Short()

		if methodID == "" || ctx.T.Tx.TxData.GetData() != "" {
			// Nothing to craft
			return
		}

		// Retrieve  args from trace
		args := ctx.T.Call.GetArgs()
		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"crafter.method": methodID,
			"crafter.args":   args,
		})

		// Retrieve method ABI object
		method, err := r.GetMethodByID(methodID)
		if err != nil {
			// e := commonpb.Error{
			// 	Message:  err.Error(),
			// 	Type: 0, // TODO: add an error type ErrorTypeABIGet
			// }
			// Abort execution
			ctx.Logger.WithError(err).Errorf("crafter: could not retrieve method ABI")
			ctx.AbortWithError(err)
			return
		}

		// Craft transaction payload
		payload, err := c.Craft(method, args...)
		if err != nil {
			// e := commonpb.Error{
			// 	Err:  err.Error(),
			// 	Type: 0, // TODO: add an error type ErrorTypeCraft
			// }
			// Abort execution
			ctx.Logger.WithError(err).Errorf("crafter: could not craft tx data payload")
			ctx.AbortWithError(err)
			return
		}

		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"tx.data":        hexutil.Encode(payload),
		})

		// Update Trace
		ctx.T.Tx.TxData.SetData(payload)

		ctx.Logger.Debugf("crafter: tx data payload set")
	}
}
