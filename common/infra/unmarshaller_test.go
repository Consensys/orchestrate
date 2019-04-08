package infra

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/envelope"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func newProtoMessage() *envelope.Envelope {
	return &envelope.Envelope{
		Sender: &common.Account{Id: "abcde"},
	}
}

func TestEnvelopeUnmarshaller(t *testing.T) {
	u := EnvelopeUnmarshaller{}
	envelopes := make([]*envelope.Envelope, 0)
	rounds := 1000
	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		envelopes = append(envelopes, &envelope.Envelope{})
		wg.Add(1)
		go func(t *envelope.Envelope) {
			defer wg.Done()
			u.Unmarshal(newProtoMessage(), t)
		}(envelopes[len(envelopes)-1])
	}
	wg.Wait()

	for _, tr := range envelopes {
		if tr.Sender.Id != "abcde" {
			t.Errorf("EnvelopeUnmarshaller: expected %q but got %q", "abcde", tr.Sender.Id)
		}
	}
}
