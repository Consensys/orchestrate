package handlers

import (
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

var testPKeys = []struct {
	prv string
	a   string
}{
	{"5FBB50BFF6DFAD35C4A374C9237BA2F7EAED9C6868E0108CB259B62D68029B1A", "0xdbb881a51CD4023E4400CEF3ef73046743f08da3"},
	{"86B021CCB810F26A30445B85F71E4C1596A11A97DDF9B9E348AC93D1DA6735BC", "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff"},
	{"DD614C3B343E1B6DBD1B2811D4F146CC90337DEEF96AB97C353578E871B19D5E", "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"},
	{"425D92F63A836F890F1690B34B6A25C2971EF8D035CD8EA8592FD1069BD151C6", "0xffbBa394DEf3Ff1df0941c6429887107f58d4e9b"},
}

var testChains = []struct {
	ID       string
	IsEIP155 bool
}{
	{"0x1ae3", true},
	{"0x3", true},
	{"0xbf6e", false},
}

type TestSignerMsg struct {
	chainID  string
	IsEIP155 bool
	a        string
}

func newSignerTestMessage(i int) *TestSignerMsg {
	return &TestSignerMsg{testChains[i%3].ID, testChains[i%3].IsEIP155, testPKeys[i%4].a}
}

func testSignerLoader() infra.HandlerFunc {
	return func(ctx *infra.Context) {
		msg := ctx.Msg.(*TestSignerMsg)
		ctx.Pb.Chain = &tracepb.Chain{Id: msg.chainID, IsEIP155: msg.IsEIP155}
		ctx.Pb.Sender = &tracepb.Account{Address: msg.a}

		// Load Trace from protobuffer
		protobuf.LoadTrace(ctx.Pb, ctx.T)
	}
}

func TestSigner(t *testing.T) {
	pKeys := []string{}
	for _, p := range testPKeys {
		pKeys = append(pKeys, p.prv)
	}
	txSigner := NewStaticSigner(pKeys)
	// Create signer handler
	h := Signer(txSigner)

	// Create new worker
	w := infra.NewWorker(100)
	w.Use(testSignerLoader())
	w.Use(h)

	testH := &testCrafterHandler{
		mux:     &sync.Mutex{},
		handled: []*infra.Context{},
	}
	w.Use(testH.Handler(50, t))

	// Create a Sarama message channel
	in := make(chan interface{})

	// Run worker
	go w.Run(in)

	// Feed sarama channel and then close it
	rounds := 1000
	for i := 1; i <= rounds; i++ {
		in <- newSignerTestMessage(i)
	}
	close(in)

	// Wait for worker to be done
	<-w.Done()

	// Run worker
	go w.Run(in)

	signers := []string{}
	txSigner.signers.Range(
		func(key, value interface{}) bool {
			signers = append(signers, key.(string))
			return true
		},
	)

	if len(signers) != len(testChains) {
		t.Errorf("Signer: Expected %v signers but got %v", len(testChains), len(signers))
	}

	if len(testH.handled) != rounds {
		t.Errorf("Signer: expected %v rounds but got %v", rounds, len(testH.handled))
	}

	for _, c := range testH.handled {
		if len(c.T.Tx().Raw()) < 95 {
			t.Errorf("Signer: Expected tx to be signed but got %q", hexutil.Encode(c.T.Tx().Raw()))
		}

		if c.T.Tx().Hash().Hex() == "0x0000000000000000000000000000000000000000000000000000000000000000" {
			t.Errorf("Signer: Expected tx hash to be set")
		}
	}
}
