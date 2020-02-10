package producer

import (
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

// PrepareMsg function should prepare a sarama.ProducerMessage from a Builder
type PrepareMsg func(*engine.TxContext, *sarama.ProducerMessage) error

// Producer creates a producer handler
func Producer(p sarama.SyncProducer, prepareMsg PrepareMsg) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Next()

		// Prepare Message
		msg := &sarama.ProducerMessage{}
		err := prepareMsg(txctx, msg)
		if err != nil {
			txctx.Logger.WithError(err).Fatalf("producer: could not prepare message")
		}

		if msg.Topic != "" {
			// Send message
			partition, offset, err := p.SendMessage(msg)
			if err != nil {
				e := txctx.AbortWithError(err).ExtendComponent(component)
				txctx.Logger.WithError(e).Errorf("producer: could not produce message")
				return
			}

			txctx.Logger = txctx.Logger.WithFields(log.Fields{
				"kafka.out.partition": partition,
				"kafka.out.offset":    offset,
				"kafka.out.topic":     msg.Topic,
			})

			txctx.Logger.Tracef("producer: message produced")
		} else {
			txctx.Logger.Tracef("producer: no message produced")
		}
	}
}
