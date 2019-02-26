package handlers

import (
	"time"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
)

// Logger to log context elements before and after the worker
func Logger(ctx *types.Context) {

	msg := ctx.Msg.(*sarama.ConsumerMessage)
	ctx.Logger = log.WithFields(infSarama.ConsumerMessageFields(msg))

	ctx.Logger.Debug("logger: new message")
	start := time.Now()

	ctx.Next()

	latency := time.Now().Sub(start)
	ctx.Logger.WithFields(log.Fields{
		"latency": latency,
	}).WithError(ctx.T.Errors).Info("logger: message processed")
}
