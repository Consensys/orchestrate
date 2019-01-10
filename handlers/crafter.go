package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
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

		// Retrieve method ABI object
		method, err := r.GetMethodByID(methodID)
		if err != nil {
			e := types.Error{
				Err:  err,
				Type: 0, // TODO: add an error type ErrorTypeABIGet
			}
			// Abort execution
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
			ctx.AbortWithError(e)
			return
		}

		// Update Trace
		ctx.T.Tx().SetData(payload)
	}
}
