package infra

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf/trace"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// TracePbMarshaller assumes that message is a trace protobuf
type TracePbMarshaller struct{}

// Marshal Trace into a message assumed to be a protobuf
func (u *TracePbMarshaller) Marshal(t *types.Trace, msg interface{}) error {
	// Cast message into trace protobuf
	pb, ok := msg.(*tracepb.Trace)
	if !ok {
		return fmt.Errorf("Message does not match expected format")
	}

	// Dump trace into protobuffer
	protobuf.DumpTrace(t, pb)
	
	return nil
}
