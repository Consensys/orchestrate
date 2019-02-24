package handlers

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Crafter creates a crafter handler
func Crafter(r services.ABIRegistry, c services.Crafter) types.HandlerFunc {
	return func(ctx *types.Context) {
		// Retrieve method identifier from trace
		methodID := ctx.T.Call().MethodID

		if methodID == "" || len(ctx.T.Tx().Data()) > 0 {
			// Nothing to craft
			return
		}

		// Retrieve  args from trace
		args := ctx.T.Call().Args
		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"crafter.method": methodID,
			"crafter.args":   args,
		})

		// Retrieve method ABI object
		method, err := r.GetMethodByID(methodID)
		if err != nil {
			e := types.Error{
				Err:  err,
				Type: 0, // TODO: add an error type ErrorTypeABIGet
			}
			// Abort execution
			ctx.Logger.WithError(err).Errorf("crafter: could not retrieve method ABI")
			ctx.AbortWithError(e)
			return
		}

		// Craft transaction payload
		payload, err := c.Craft(method, args...)
		if err != nil {
			e := types.Error{
				Err:  err,
				Type: 0, // TODO: add an error type ErrorTypeCraft
			}
			// Abort execution
			ctx.Logger.WithError(err).Errorf("crafter: could not craft tx data payload")
			ctx.AbortWithError(e)
			return
		}

		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"tx.data":        hexutil.Encode(payload),
		})

		// Update Trace
		ctx.T.Tx().SetData(payload)

		ctx.Logger.Debugf("crafter: tx data payload set")
	}
}
