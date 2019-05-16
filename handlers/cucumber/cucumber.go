package logger

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Logger creates a handler middleware that log basic information about tx execution
func Cucumber() engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		start := time.Now()

		txctx.Next()

		txctx.Logger.
			WithFields(log.Fields{
				"latency": time.Since(start),
				"chainId": txctx.Envelope.GetChain().GetId(),
			}).
			Info("cucumber: message processed")
		}
}
