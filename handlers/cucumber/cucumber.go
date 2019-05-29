package cucumber

import (
	"fmt"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/tests/e2e.git/cucumber/chanregistry"
)

// Cucumber
func Cucumber(c *chanregistry.ChanRegistry) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {

		txctx.Next()

		msg, ok := txctx.Msg.(*sarama.ConsumerMessage)
		if !ok {
			txctx.Logger.Errorf("loader: expected a sarama.ConsumerMessage")
			_ = txctx.AbortWithError(fmt.Errorf("invalid input message format"))
			return
		}

		extra := txctx.Envelope.GetMetadata().GetExtra()

		if extra["ScenarioID"] != "" {
			envelopeChan := c.GetEnvelopeChan(extra["ScenarioID"], msg.Topic)
			if envelopeChan != nil {
				envelopeChan <- txctx.Envelope
			}
		} else {
			txctx.Logger.
				WithFields(log.Fields{
					"ScenarioID": txctx.Envelope.GetMetadata().GetId(),
					"msg.Topic":  msg.Topic,
				}).
				Error("cucumber: received unknown envelope")
		}

	}

}
