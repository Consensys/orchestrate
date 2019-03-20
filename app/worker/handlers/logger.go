package handlers

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// Logger to log context elements before and after the worker
func Logger(ctx *worker.Context) {
	ctx.Logger.Debug("logger: new message")
	start := time.Now()

	ctx.Next()

	latency := time.Now().Sub(start)
	ctx.Logger.WithFields(log.Fields{
		"latency": latency,
	}).WithError(fmt.Errorf("%q", ctx.T.GetErrors())).Info("logger: message processed")
}
