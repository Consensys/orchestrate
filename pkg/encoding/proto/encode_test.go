// +build unit

package proto

import (
	"testing"

	"github.com/consensys/orchestrate/pkg/types/tx"

	"google.golang.org/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestMarshaller(t *testing.T) {
	uuid := "957e4d77-5dbb-484e-998c-e44c0dd8891d"

	msg := &tx.TxRequest{
		Id: uuid,
	}
	b, _ := Marshal(msg)

	var evlp = &tx.TxRequest{}
	_ = proto.Unmarshal(b, evlp)

	if evlp.GetId() != uuid {
		t.Errorf("EnvelopeMarshaller: expected %q but got %q", uuid, evlp.GetId())
	}
}

func TestMarshallerError(t *testing.T) {
	b, err := Marshal(nil)
	assert.NoError(t, err, "Marshal should not error on nil")
	assert.Equal(t, []byte(nil), b, "Marshal nil should be correct")
}
