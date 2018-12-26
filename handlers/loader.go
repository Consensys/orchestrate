package handlers

import (
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

// SaramaLoader return an handler to load messages from sarama
func SaramaLoader() infra.HandlerFunc {
	return func(ctx *infra.Context) {
		// Cast message into a sarama.ConsumerMessage
		msg := ctx.Msg.(*sarama.ConsumerMessage)

		// Unmarshal Sarama message to protobuffer
		err := proto.Unmarshal(msg.Value, ctx.Pb)
		if err != nil {
			// Indicate error for a possible middleware to recover it
			e := &types.Error{
				Err:  err,
				Type: types.ErrorTypeLoad,
			}
			ctx.Error(e)
			return
		}

		// Load Trace from protobuffer
		protobuf.LoadTrace(ctx.Pb, ctx.T)
	}
}
