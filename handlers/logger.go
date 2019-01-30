package handlers

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Logger to log context elements before and after the worker
func Logger(ctx *types.Context) {
	msg := ctx.Msg.(*sarama.ConsumerMessage)

	log.WithFields(log.Fields{
		"Offset":  msg.Offset,
		"ChainID": ctx.T.Chain().ID.Text(10),
		"TxHash":  ctx.T.Tx().Hash().Hex(),
	}).Infof("Crafter [IN]\nRaw: %v\n", hexutil.Encode(ctx.T.Tx().Raw()))

	ctx.Next()

	log.WithFields(log.Fields{
		"Offset": msg.Offset,
	}).Infof("Crafter [OUT]\nErrors: %v\n", ctx.T.Errors)

	errors := ctx.T.Errors
	if len(errors) > 0 {
		// TODO: change to log
		fmt.Printf("Error: %v\n", errors)
	}
}
