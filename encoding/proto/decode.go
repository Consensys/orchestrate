package proto

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

// Unmarshal message expected to be a Envelope
func Unmarshal(msg interface{}, pb proto.Message) error {
	// Cast message into protobuffer
	e, ok := msg.(*envelope.Envelope)
	if !ok {
		return fmt.Errorf("message does not match expected format")
	}

	// Load Envelope from protobuffer
	proto.Merge(pb, e)

	return nil

}
