package offset

import (
	"github.com/Shopify/sarama"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Marker is a Middleware handler that marks offsets
func Marker(txctx *engine.TxContext) {
	// Marker only executes in second phase of middleware
	txctx.Next()

	// Extract sarama ConsumerGroupSession from context
	s, _ := broker.GetConsumerGroupSessionAndClaim(txctx.Context())
	if s != nil {
		// Cast message
		msg, ok := txctx.Msg.(*broker.Msg)
		if !ok {
			return
		}
		s.MarkMessage((*sarama.ConsumerMessage)(msg), "")
	}
}
