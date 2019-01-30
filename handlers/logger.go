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
		"Address": ctx.T.Sender().Address.Hex(),
	}).Info("Crafter [IN]\n")

	ctx.Next()

	log.WithFields(log.Fields{
		"Offset": msg.Offset,
	}).Infof("Crafter [OUT]\n Data: %v\nGasPrice: %v\nGas Limit: %v\nErrors: %v",
		hexutil.Encode(ctx.T.Tx().Data()), hexutil.EncodeBig(ctx.T.Tx().GasPrice()), ctx.T.Tx().GasLimit(), ctx.T.Errors)

	errors := ctx.T.Errors
	if len(errors) > 0 {
		// TODO: change to log
		fmt.Printf("Error: %v\n", errors)
	}
}
