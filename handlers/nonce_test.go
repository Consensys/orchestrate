package handlers

import (
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

func newNonceTest(a common.Address) (*SafeNonce, error) {
	return &SafeNonce{10, &sync.Mutex{}}, nil
}

func TestNonceHandler(t *testing.T) {
	// Create handler
	m := NewCacheNonceManager(newNonceTest)
	handler := NonceHandler(m)

	// Create a context
	ctx := infra.NewContext()
	ctx.Init([]infra.HandlerFunc{handler})

	// Set Trace values
	a := common.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684")
	ctx.T.Chain().ID = "abc"
	ctx.T.Sender().Address = &a

	// Execute handler
	ctx.Next()

	if ctx.T.Tx().Nonce() != 10 {
		t.Errorf("NonceHandler: expected nonce to be %v but got %v", 10, ctx.T.Tx().Nonce())
	}
}
