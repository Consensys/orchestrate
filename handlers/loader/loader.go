package loader

import (
	log "github.com/sirupsen/logrus"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Loader is an handler that Load sarama.ConsumerGroup messages
func Loader(txctx *engine.TxContext) {
	// Cast message into sarama.ConsumerMessage
	msg, ok := txctx.In.(*broker.Msg)
	if !ok {
		txctx.Logger.Fatalf("loader: expected a sarama.ConsumerMessage")
	}

	// Enrich Logger
	txctx.Logger = txctx.Logger.WithFields(log.Fields{
		"kafka.in.topic":     msg.Topic,
		"kafka.in.offset":    msg.Offset,
		"kafka.in.partition": msg.Partition,
	})

	err := encoding.Unmarshal(msg, txctx.Envelope)
	if err != nil {
		e := txctx.AbortWithError(err).ExtendComponent(component)
		txctx.Logger.WithError(e).Errorf("loader: error unmarshalling")
		return
	}

	txctx.Logger.Tracef("loader: message loaded: %v", txctx.Envelope.String())
}
