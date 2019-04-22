package loader

import (
	"fmt"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

func init() {
	unmarshaller = encoding.NewUnmarshaller()
}

var unmarshaller *encoding.Unmarshaller

// Loader is a Middleware enginer.HandlerFunc that Load sarama.ConsumerGroup messages
func Loader(txctx *engine.TxContext) {
	// Cast message into sarama.ConsumerMessage
	msg, ok := txctx.Msg.(*sarama.ConsumerMessage)
	if !ok {
		txctx.Logger.Errorf("loader: expected a sarama.ConsumerMessage")
		_ = txctx.AbortWithError(fmt.Errorf("invalid input message format"))
		return
	}

	err := unmarshaller.Unmarshal(msg, txctx.Envelope)
	if err != nil {
		// TODO: handle error
		txctx.Logger.WithError(err).Errorf("loader: error unmarshalling")
		_ = txctx.AbortWithError(err)
		return
	}

	// Enrich Logger
	txctx.Logger = txctx.Logger.WithFields(log.Fields{
		"kafka.in.topic":     msg.Topic,
		"kafka.in.offset":    msg.Offset,
		"kafka.in.partition": msg.Partition,
	})

	txctx.Logger.Tracef("loader: message loaded: %v", txctx.Envelope.String())
}
