package infra

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf/trace"
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
