package proto

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var msg = &envelope.Envelope{
	From: &ethereum.Account{
		Raw: hexutil.MustDecode("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"),
	},
}
var buf, _ = proto.Marshal(msg)

func TestUnmarshaller(t *testing.T) {
	envelopes := make([]*envelope.Envelope, 0)
	rounds := 1000
	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		envelopes = append(envelopes, &envelope.Envelope{})
		wg.Add(1)
		go func(t *envelope.Envelope) {
			defer wg.Done()
			_ = Unmarshal(buf, t)
		}(envelopes[len(envelopes)-1])
	}
	wg.Wait()

	for _, tr := range envelopes {
		assert.Equal(t, msg.GetFrom().Hex(), tr.GetFrom().Hex(), "Expected correct unmarshalled message")
	}
}

func TestUnmarshalError(t *testing.T) {
	pb := &ethereum.TxData{}
	err := errors.FromError(Unmarshal([]byte{0xab, 0x00}, pb))
	assert.NotNil(t, err, "Unmarshal should error")
	assert.Equal(t, err.GetComponent(), "encoding.proto", "Error code should be correct")
}

func TestUnmarshalMergeError(t *testing.T) {
	pb := &ethereum.TxData{}
	err := errors.FromError(UnmarshalMerge([]byte{0xab, 0x00}, pb))
	assert.NotNil(t, err, "UnmarshalMerge should error")
	assert.Equal(t, err.GetComponent(), "encoding.proto", "Error code should be correct")
}
