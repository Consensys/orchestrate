package proto

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

func TestEnvelopeMarshaller(t *testing.T) {
	b, _ := Marshal(msg)
	var evlp *envelope.Envelope
	_ = proto.Unmarshal(b, evlp)

	if evlp.GetFrom().Hex() != "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C" {
		t.Errorf("EnvelopeMarshaller: expected %q but got %q", "abcde", evlp.GetFrom().Hex())
	}
}
