package offset

import (
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

// Marker is a Middleware handler that marks offsets
func Marker(txctx *engine.TxContext) {
	// Marker only executes in second phase of middleware
	txctx.Next()

	// Extract sarama ConsumerGroupSession from context
	s, _ := broker.GetConsumerGroupSessionAndClaim(txctx.Context())
	if s != nil {
		// Cast message
		msg, ok := txctx.In.(*broker.Msg)
		if !ok {
			txctx.Logger.Fatalf("marker: expected a sarama.ConsumerMessage")
		}
		s.MarkMessage(&msg.ConsumerMessage, "")
	}
}
