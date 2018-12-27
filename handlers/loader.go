package handlers

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

// CtxUnmarshaller are responsible to unmarshal an input message to protobuf
type CtxUnmarshaller interface {
	Unmarshal(ctx *infra.Context) error
}

// TraceProtoUnmarshaller assumes that input messages are protobuf
type TraceProtoUnmarshaller struct{}

// Unmarshal message
func (u *TraceProtoUnmarshaller) Unmarshal(ctx *infra.Context) error {
	// Cast message into a sarama.ConsumerMessage
	var ok bool
	ctx.Pb, ok = ctx.Msg.(*tracepb.Trace)

	if !ok {
		return fmt.Errorf("Input does not match expected format")
	}

	return nil
}

// SaramaUnmarshaller assumes that input messages is a Sarama message
type SaramaUnmarshaller struct{}

// Unmarshal message
func (u *SaramaUnmarshaller) Unmarshal(ctx *infra.Context) error {
	// Cast message into a sarama.ConsumerMessage
	msg, ok := ctx.Msg.(*sarama.ConsumerMessage)

	if !ok {
		return fmt.Errorf("Input does not match expected format")
	}

	// Unmarshal Sarama message to protobuffer
	err := proto.Unmarshal(msg.Value, ctx.Pb)
	if err != nil {
		return err
	}

	return nil

}

// Loader creates an handler loading input
func Loader(u CtxUnmarshaller) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		// Unmarshal message
		u.Unmarshal(ctx)

		// Load Trace from protobuffer
		protobuf.LoadTrace(ctx.Pb, ctx.T)
	}
}
