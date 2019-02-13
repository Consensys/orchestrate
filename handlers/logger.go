package handlers

import (
	"fmt"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// LoggerHandler ...
func LoggerHandler(ctx *types.Context) {
	msg := ctx.Msg.(*sarama.ConsumerMessage)

	log.WithFields(log.Fields{
		"Offset": msg.Offset,
	}).Info("Logger [IN]")

	ctx.Next()

	for _, logData := range ctx.T.Receipt().Logs {
		log.WithFields(log.Fields{
			"Offset": msg.Offset,
			"Log":    logData.DecodedData,
		}).Info("Logger [OUT]")
	}

	errors := ctx.T.Errors
	if len(errors) > 0 {
		// TODO: change to log
		fmt.Printf("Error: %v\n", errors)
	}
}
