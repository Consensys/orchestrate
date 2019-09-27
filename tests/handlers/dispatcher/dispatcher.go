package dispatcher

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/tests/e2e.git/service/chanregistry"
)

// Dispacher is dispatching envelopes to registered channels
func Dispacher(c chanregistry.ChanRegistry) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {

		txctx.Next()

		msg := txctx.In
		if msg == nil {
			txctx.Logger.Error("dispacher: received invalid message")
			_ = txctx.AbortWithError(fmt.Errorf("invalid input message format"))
			return
		}

		extra := txctx.Envelope.GetMetadata().GetExtra()

		if extra["ScenarioID"] != "" {
			envelopeChan := c.GetEnvelopeChan(extra["ScenarioID"], msg.Entrypoint())
			if envelopeChan != nil {
				envelopeChan <- txctx.Envelope
				return
			}
		}

		txctx.Logger.
			WithFields(log.Fields{
				"MetadataID": txctx.Envelope.GetMetadata().GetId(),
				"msg.Topic":  msg.Entrypoint(),
			}).
			Error("dispacher: received unknown envelope")
		_ = txctx.AbortWithError(fmt.Errorf("scenarioID unknown, envelope not dispatched"))
	}
}
