package handlers

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// Logger to log context elements before and after the worker
func Logger(ctx *worker.Context) {
	msg := ctx.Msg.(*sarama.ConsumerMessage)
	ctx.Logger = log.WithFields(infSarama.ConsumerMessageFields(msg))

	ctx.Logger.Debug("logger: new message")
	start := time.Now()

	ctx.Next()

	latency := time.Now().Sub(start)
	ctx.Logger.WithFields(log.Fields{
		"latency": latency,
	}).WithError(fmt.Errorf("%q", ctx.T.GetErrors())).Info("logger: message processed")
}
