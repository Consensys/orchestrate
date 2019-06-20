package proto

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newProtoMessage() *envelope.Envelope {
	return &envelope.Envelope{
		From: &ethereum.Account{
			Raw: hexutil.MustDecode("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"),
		},
	}
}

func TestUnmarshaller(t *testing.T) {
	envelopes := make([]*envelope.Envelope, 0)
	rounds := 1000
	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		envelopes = append(envelopes, &envelope.Envelope{})
		wg.Add(1)
		go func(t *envelope.Envelope) {
			defer wg.Done()
			_ = Unmarshal(newProtoMessage(), t)
		}(envelopes[len(envelopes)-1])
	}
	wg.Wait()

	for _, tr := range envelopes {
		if tr.GetFrom().Hex() != newProtoMessage().GetFrom().Hex() {
			t.Errorf("EnvelopeUnmarshaller: expected %q but got %q", "abcd", tr.GetFrom().Hex())
		}
	}
}
