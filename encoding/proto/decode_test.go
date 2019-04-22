package proto

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newProtoMessage() *envelope.Envelope {
	return &envelope.Envelope{
		Sender: &common.Account{Id: "abcd"},
	}
}

func TestUnmarshaller(t *testing.T) {
	u := Unmarshaller{}
	envelopes := make([]*envelope.Envelope, 0)
	rounds := 1000
	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		envelopes = append(envelopes, &envelope.Envelope{})
		wg.Add(1)
		go func(t *envelope.Envelope) {
			defer wg.Done()
			_ = u.Unmarshal(newProtoMessage(), t)
		}(envelopes[len(envelopes)-1])
	}
	wg.Wait()

	for _, tr := range envelopes {
		if tr.Sender.Id != "abcd" {
			t.Errorf("EnvelopeUnmarshaller: expected %q but got %q", "abcd", tr.Sender.Id)
		}
	}
}
