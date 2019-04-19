package proto

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

// EnvelopeMarshaller assumes that message is an Envelope
type EnvelopeMarshaller struct{}

// Marshal Envelope into a message assumed to be a protobuf
func (u *EnvelopeMarshaller) Marshal(t *envelope.Envelope, msg interface{}) error {
	// Cast message into Envelope
	pb, ok := msg.(*envelope.Envelope)
	if !ok {
		return fmt.Errorf("Message does not match expected format")
	}

	// Merge msg into Envelope
	proto.Merge(pb, t)

	return nil
}
