package handlers

import "gitlab.com/ConsenSys/client/fr/core-stack/core/infra"

// ErrorHandler return an Handler for error handling
func ErrorHandler() infra.HandlerFunc {
	return func(ctx *infra.Context) {
		// TODO: process errors before handling

		ctx.Next()

		// TODO: process errors after handling
	}
}
