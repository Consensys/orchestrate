package proto

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

// Marshaller assumes that message is an Envelope
type Marshaller struct{}

// Marshal a proto into a message assumed to be an Envelope
func (u *Marshaller) Marshal(pb proto.Message, msg interface{}) error {
	// Cast message into Envelope
	e, ok := msg.(*envelope.Envelope)
	if !ok {
		return fmt.Errorf("message does not match expected format")
	}

	// Merge msg into Envelope
	proto.Merge(e, pb)

	return nil
}
