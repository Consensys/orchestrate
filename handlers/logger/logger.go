package logger

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
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

		txctx.Logger = txctx.Logger.
			WithFields(log.Fields{
				"latency": fmt.Sprintf("%vms", time.Since(start).Milliseconds()),
			})

		switch {
		case len(txctx.Envelope.GetErrors()) > 0 && txctx.Envelope.OnlyWarnings():
			txctx.Logger.
				WithError(fmt.Errorf("%q", txctx.Envelope.GetErrors())).
				Warn("message processed with warning")
		case len(txctx.Envelope.GetErrors()) > 0:
			txctx.Logger.
				WithError(fmt.Errorf("%q", txctx.Envelope.GetErrors())).
				Error("message processed with error")
		default:
			txctx.Logger.
				Log(logLevel, "message processed")
		}
	}
}
