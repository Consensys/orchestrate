package infra

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/envelope"
)

// EnvelopeUnmarshaller assumes that input message is a Envelope
type EnvelopeUnmarshaller struct{}

// Unmarshal message expected to be a Envelope
func (u *EnvelopeUnmarshaller) Unmarshal(msg interface{}, t *envelope.Envelope) error {
	// Cast message into protobuffer
	pb, ok := msg.(*envelope.Envelope)
	if !ok {
		return fmt.Errorf("Message does not match expected format")
	}

	// Load Envelope from protobuffer
	proto.Merge(t, pb)

	return nil
}
