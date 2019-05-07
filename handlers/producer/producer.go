package producer

import (
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

// PrepareMsg function should prepare a sarama.ProducerMessage from a TxContext
type PrepareMsg func(*engine.TxContext, *sarama.ProducerMessage) error

// Producer creates a producer handler
func Producer(p sarama.SyncProducer, prepareMsg PrepareMsg) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Next()

		// Prepare Message
		msg := &sarama.ProducerMessage{}
		err := prepareMsg(txctx, msg)
		if err != nil {
			_ = txctx.AbortWithError(err)
			txctx.Logger.WithError(err).Errorf("producer: could not prepare message")
			return
		}

		// Send message
		partition, offset, err := p.SendMessage(msg)
		if err != nil {
			_ = txctx.AbortWithError(err)
			txctx.Logger.WithError(err).Errorf("producer: could not produce message")
			return
		}

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"kafka.out.partition": partition,
			"kafka.out.offset":    offset,
			"kafka.out.topic":     msg.Topic,
		})
		txctx.Logger.Tracef("producer: message produced")
	}
}

// MultiProducer creates a multi producer handler
func MultiProducer(p sarama.SyncProducer, prepareMsg PrepareMsg) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Next()

		envelopes := txctx.Get("envelopes").([]*envelope.Envelope)
		for _, e := range envelopes {
			// Prepare Message
			subTxctx := &engine.TxContext{
				Envelope: e,
				Logger:   txctx.Logger,
			}
			subTxctx.WithContext(txctx.Context())
			subTxctx.Set("envelopes", nil)

			Producer(p, prepareMsg)(subTxctx)
		}
	}
}
