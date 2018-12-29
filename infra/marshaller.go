package infra

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

// Marshaller are responsible to marshal
type Marshaller interface {
	Marshal(pb *tracepb.Trace, msg interface{}) error
}

// TraceProtoMarshaller assumes that message is a protobuf
type TraceProtoMarshaller struct{}

// Marshal message
func (u *TraceProtoMarshaller) Marshal(pb *tracepb.Trace, msg interface{}) error {
	// Cast message into trace protobuf
	var cast, ok = msg.(*tracepb.Trace)
	if !ok {
		return fmt.Errorf("Message does not match expected format")
	}
	cast.Reset()
	proto.Merge(cast, pb)
	return nil
}

// SaramaMarshaller assumes that input messages is a Sarama message
type SaramaMarshaller struct{}

// Marshal message
func (u *SaramaMarshaller) Marshal(pb *tracepb.Trace, msg interface{}) error {
	// Cast message into a sarama.ConsumerMessage
	var cast, ok = msg.(*sarama.ProducerMessage)
	if !ok {
		return fmt.Errorf("Message does not match expected format")
	}

	// Marshal protobuffer
	b, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	cast.Value = sarama.ByteEncoder(b)
	return nil
}
