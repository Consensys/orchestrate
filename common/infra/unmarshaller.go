package infra

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protobuf"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protobuf/trace"
)

// TracePbUnmarshaller assumes that input message is a protobuf
type TracePbUnmarshaller struct{}

// Unmarshal message expected to be a trace protobuffer
func (u *TracePbUnmarshaller) Unmarshal(msg interface{}, t *types.Trace) error {
	// Cast message into protobuffer
	pb, ok := msg.(*tracepb.Trace)
	if !ok {
		return fmt.Errorf("Message does not match expected format")
	}

	// Load trace from protobuffer
	protobuf.LoadTrace(pb, t)

	return nil
}
