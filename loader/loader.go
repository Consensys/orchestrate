package loader

import (
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
	// Unmarshal message
	err := unmarshaller.Unmarshal(txctx.Msg, txctx.Envelope)
	if err != nil {
		// TODO: handle error
		txctx.Logger.Errorf("Error unmarshalling: %v", err)
		txctx.AbortWithError(err)
		return
	}

	// Enrich Logger
	msg := txctx.Msg.(*sarama.ConsumerMessage)
	txctx.Logger = txctx.Logger.WithFields(log.Fields{
		"kafka.in.topic":     msg.Topic,
		"kafka.in.offset":    msg.Offset,
		"kafka.in.partition": msg.Partition,
	})

	txctx.Logger.Tracef("Message loaded: %v", txctx.Envelope.String())
}
