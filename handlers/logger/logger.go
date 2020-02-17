package logger

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

// Logger creates a handler middleware that log basic information about tx execution
func Logger(level string) engine.HandlerFunc {
	logLevel, err := log.ParseLevel(level)
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}

	return func(txctx *engine.TxContext) {
		txctx.Logger.Trace("logger: new message")
		start := time.Now()

		txctx.Next()

		txctx.Logger.
			WithFields(log.Fields{
				"latency": time.Since(start),
			}).
			WithError(fmt.Errorf("%q", txctx.Envelope.GetErrors())).
			Log(logLevel, "logger: message processed")
	}
}
