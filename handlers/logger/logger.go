package logger

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Logger creates a handler middleware that log basic information about tx execution
func Logger(txctx *engine.TxContext) {
	// TODO: logger enrichment should append in Loader middleware, so we remove dependency to sarama in Logger
	// msg := txctx.Msg.(*sarama.ConsumerMessage)
	// txctx.Logger = log.WithFields(infsarama.ConsumerMessageFields(msg))

	txctx.Logger.Trace("logger: new message")
	start := time.Now()

	txctx.Next()

	txctx.Logger.
		WithFields(log.Fields{
			"latency": time.Since(start),
		}).
		WithError(fmt.Errorf("%q", txctx.Envelope.GetErrors())).
		Info("logger: message processed")
}
