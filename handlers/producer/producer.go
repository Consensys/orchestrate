package producer

import (
	"fmt"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
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

// MultiProducer creates a multi producer handler
func MultiProducer(p sarama.SyncProducer, prepareMsg PrepareMsg) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Next()

		// Pass through the handler if no envelopes to send
		if txctx.Get("envelopes") == nil {
			return
		}

		// Test if able to cast txctx into []*envelope.Envelope
		envelopes, ok := txctx.Get("envelopes").([]*envelope.Envelope)
		if !ok {
			err := fmt.Errorf("not able to cast envelopes %q", envelopes)
			_ = txctx.AbortWithError(err)
			txctx.Logger.WithError(err).Errorf("multiProducer: could not produce messages")
			return
		}

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
