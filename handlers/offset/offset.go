package offset

import (
	broker "github.com/consensys/orchestrate/pkg/broker/sarama"
	"github.com/consensys/orchestrate/pkg/engine"
)

// Marker is a Middleware handler that marks offsets
func Marker(txctx *engine.TxContext) {
	// Marker only executes in second phase of middleware
	txctx.Next()

	// In case of msg retrying, skip offset marker
	if err := txctx.HasRetryMsgErr(); err != nil {
		txctx.Logger.WithError(err).Warn("marker - skip offset marking, retrying message...")
		return
	}

	// Extract sarama ConsumerGroupSession from context
	s, _ := broker.GetConsumerGroupSessionAndClaim(txctx.Context())
	if s != nil {
		// Cast message
		msg, ok := txctx.In.(*broker.Msg)
		if !ok {
			txctx.Logger.Fatalf("marker - expected a sarama.ConsumerMessage")
		}

		txctx.Logger.Debug("marker - commit offset")
		s.MarkMessage(&msg.ConsumerMessage, "")
		s.Commit()
	}
}
