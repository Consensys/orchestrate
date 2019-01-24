package handlers

import (
	"fmt"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Logger to log context elements before and after the worker
func Logger(ctx *types.Context) {
	msg := ctx.Msg.(*sarama.ConsumerMessage)

	log.WithFields(log.Fields{
		"Offset": msg.Offset,
		"Nonce":  ctx.T.Tx().Nonce(),
	}).Info("Nonce [IN]")

	ctx.Next()

	log.WithFields(log.Fields{
		"Offset": msg.Offset,
		"Nonce":  ctx.T.Tx().Nonce(),
	}).Info("Nonce [OUT]")

	errors := ctx.T.Errors
	if len(errors) > 0 {
		// TODO: change to log
		fmt.Printf("Error: %v\n", errors)
	}
}
