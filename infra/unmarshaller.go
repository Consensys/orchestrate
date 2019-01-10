package infra

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

// TraceProtoUnmarshaller assumes that input message is a protobuf
type TraceProtoUnmarshaller struct{}

// Unmarshal message
func (u *TraceProtoUnmarshaller) Unmarshal(msg interface{}, pb *tracepb.Trace) error {
	// Cast message into a sarama.ConsumerMessage
	var cast, ok = msg.(*tracepb.Trace)
	if !ok {
		return fmt.Errorf("Message does not match expected format")
	}
	pb.Reset()
	proto.Merge(pb, cast)
	return nil
}

// SaramaUnmarshaller assumes that input messages is a Sarama message
type SaramaUnmarshaller struct{}

// Unmarshal message
func (u *SaramaUnmarshaller) Unmarshal(msg interface{}, pb *tracepb.Trace) error {
	// Cast message into a sarama.ConsumerMessage
	var cast, ok = msg.(*sarama.ConsumerMessage)
	if !ok {
		return fmt.Errorf("Message does not match expected format")
	}

	// Unmarshal Sarama message to protobuffer
	err := proto.Unmarshal(cast.Value, pb)
	if err != nil {
		return err
	}

	return nil
}
