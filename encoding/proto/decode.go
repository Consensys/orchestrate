package proto

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

// EnvelopeUnmarshaller assumes that input message is a Envelope
type EnvelopeUnmarshaller struct{}

// Unmarshal message expected to be a Envelope
func (u *EnvelopeUnmarshaller) Unmarshal(msg interface{}, e *envelope.Envelope) error {
	// Cast message into protobuffer
	pb, ok := msg.(*envelope.Envelope)
	if !ok {
		return fmt.Errorf("Message does not match expected format")
	}

	// Load Envelope from protobuffer
	proto.Merge(e, pb)

	return nil

}
