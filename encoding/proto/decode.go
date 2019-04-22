package proto

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

// Unmarshaller assumes that input message is a Envelope
type Unmarshaller struct{}

// Unmarshal message expected to be a Envelope
func (u *Unmarshaller) Unmarshal(msg interface{}, pb proto.Message) error {
	// Cast message into protobuffer
	e, ok := msg.(*envelope.Envelope)
	if !ok {
		return fmt.Errorf("message does not match expected format")
	}

	// Load Envelope from protobuffer
	proto.Merge(pb, e)

	return nil

}
