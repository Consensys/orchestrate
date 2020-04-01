// +build unit

package proto

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const uuid = "957e4d77-5dbb-484e-998c-e44c0dd8891d"

var msg = &tx.TxRequest{
	Id: uuid,
}
var buf, _ = proto.Marshal(msg)

func TestUnmarshaller(t *testing.T) {
	envelopes := make([]*tx.TxRequest, 0)
	rounds := 1000
	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		envelopes = append(envelopes, &tx.TxRequest{})
		wg.Add(1)
		go func(t *tx.TxRequest) {
			defer wg.Done()
			_ = Unmarshal(buf, t)
		}(envelopes[len(envelopes)-1])
	}
	wg.Wait()

	for _, tr := range envelopes {
		assert.Equal(t, msg.GetId(), tr.GetId(), "Expected correct unmarshalled message")
	}
}

func TestUnmarshalError(t *testing.T) {
	pb := &tx.TxRequest{}
	err := errors.FromError(Unmarshal([]byte{0xab, 0x00}, pb))
	assert.Error(t, err, "Unmarshal should error")
	assert.Equal(t, err.GetComponent(), "encoding.proto", "Error code should be correct")
}

func TestUnmarshalMergeError(t *testing.T) {
	pb := &tx.TxRequest{}
	err := errors.FromError(UnmarshalMerge([]byte{0xab, 0x00}, pb))
	assert.Error(t, err, "UnmarshalMerge should error")
	assert.Equal(t, err.GetComponent(), "encoding.proto", "Error code should be correct")
}
