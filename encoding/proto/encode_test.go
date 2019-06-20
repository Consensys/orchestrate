package proto

import (
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

func newEnvelope() *envelope.Envelope {
	return &envelope.Envelope{
		From: &ethereum.Account{
			Raw: hexutil.MustDecode("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"),
		},
	}
}

func TestEnvelopeMarshaller(t *testing.T) {
	pbs := make([]*envelope.Envelope, 0)
	rounds := 1000

	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		pbs = append(pbs, &envelope.Envelope{})
		wg.Add(1)
		go func(pb *envelope.Envelope) {
			defer wg.Done()
			_ = Marshal(newEnvelope(), pb)
		}(pbs[len(pbs)-1])
	}
	wg.Wait()

	for _, pb := range pbs {
		if pb.GetFrom().Hex() != "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C" {
			t.Errorf("EnvelopeMarshaller: expected %q but got %q", "abcde", pb.GetFrom().Hex())
		}
	}
}
