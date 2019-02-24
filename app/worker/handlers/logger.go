package handlers

import (
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Logger ...
func Logger(ctx *types.Context) {
	msg := ctx.Msg.(*sarama.ConsumerMessage)

	log.WithFields(log.Fields{
		"Offset": msg.Offset,
		"Sender": ctx.T.Sender().Address.Hex(),
	}).Info("Logger [IN]\n")

	ctx.Next()

	log.WithFields(log.Fields{
		"Offset": msg.Offset,
	}).Infof("Logger [OUT]\nRaw: %v\nHash: %v\nErrors: %v\n", hexutil.Encode(ctx.T.Tx().Raw()), ctx.T.Tx().Hash().Hex(), ctx.T.Errors)
}
