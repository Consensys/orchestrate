package proto

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/ethereum"
)

func TestMarshaller(t *testing.T) {
	msg := &envelope.Envelope{
		From: &ethereum.Account{
			Raw: hexutil.MustDecode("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"),
		},
	}
	b, _ := Marshal(msg)

	var evlp = &envelope.Envelope{}
	_ = proto.Unmarshal(b, evlp)

	if evlp.GetFrom().Hex() != "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C" {
		t.Errorf("EnvelopeMarshaller: expected %q but got %q", "abcde", evlp.GetFrom().Hex())
	}
}

func TestMarshallerError(t *testing.T) {
	b, err := Marshal(nil)
	assert.NotNil(t, err, "Marshal should error on nil")
	assert.Equal(t, []byte(nil), b, "Marshal nil should be correct")
}
