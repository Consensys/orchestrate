package handlers

import "gitlab.com/ConsenSys/client/fr/core-stack/core/types"

// ErrorHandler return an Handler for error handling
func ErrorHandler() types.HandlerFunc {
	return func(ctx *types.Context) {
		// TODO: process errors before handling

		ctx.Next()

		// TODO: process errors after handling
	}
}
