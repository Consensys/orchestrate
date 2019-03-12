package infra

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

// TracePbUnmarshaller assumes that input message is a protobuf
type TracePbUnmarshaller struct{}

// Unmarshal message expected to be a trace protobuffer
func (u *TracePbUnmarshaller) Unmarshal(msg interface{}, t *trace.Trace) error {
	// Cast message into protobuffer
	pb, ok := msg.(*trace.Trace)
	if !ok {
		return fmt.Errorf("Message does not match expected format")
	}

	// Load trace from protobuffer
	proto.Merge(t, pb)

	return nil
}
