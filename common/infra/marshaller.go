package infra

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

// TracePbMarshaller assumes that message is a trace protobuf
type TracePbMarshaller struct{}

// Marshal Trace into a message assumed to be a protobuf
func (u *TracePbMarshaller) Marshal(t *trace.Trace, msg interface{}) error {
	// Cast message into trace protobuf
	pb, ok := msg.(*trace.Trace)
	if !ok {
		return fmt.Errorf("Message does not match expected format")
	}

	// Merge msg into trace
	proto.Merge(pb, t)

	return nil
}
