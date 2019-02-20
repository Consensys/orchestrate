package handlers

import (
	"time"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
)

// LoggerHandler ...
func LoggerHandler(ctx *types.Context) {
	start := time.Now()
	msg := ctx.Msg.(*sarama.ConsumerMessage)
	ctx.Logger = log.WithFields(infSarama.ConsumerMessageFields(msg))

	ctx.Logger.WithFields(log.Fields{
		"start": start,
	}).Debug("New message")

	ctx.Next()

	latency := time.Now().Sub(start)
	if len(ctx.T.Errors) > 0 {
		ctx.Logger.WithFields(log.Fields{
			"latency": latency,
		}).Errorf("Error processing message %v", ctx.T.Errors)
	} else {
		ctx.Logger.WithFields(log.Fields{
			"latency": latency,
		}).Debug("Message processed")
	}
}
