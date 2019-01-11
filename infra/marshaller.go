package infra

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf/trace"
)

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
