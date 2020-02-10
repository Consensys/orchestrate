package nonceattributor

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"

	ethcommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/nonce/memory"
)

type MockChainStateReader struct{}

func (r *MockChainStateReader) BalanceAt(_ context.Context, _ string, _ ethcommon.Address, _ *big.Int) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (r *MockChainStateReader) StorageAt(_ context.Context, _ string, _ ethcommon.Address, _ ethcommon.Hash, _ *big.Int) ([]byte, error) {
	return []byte{}, nil
}

func (r *MockChainStateReader) CodeAt(_ context.Context, _ string, _ ethcommon.Address, _ *big.Int) ([]byte, error) {
	return []byte{}, nil
}

func (r *MockChainStateReader) NonceAt(_ context.Context, _ string, _ ethcommon.Address, _ *big.Int) (uint64, error) {
	return 0, nil
}

func (r *MockChainStateReader) PendingBalanceAt(_ context.Context, _ string, _ ethcommon.Address) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (r *MockChainStateReader) PendingStorageAt(_ context.Context, _ string, _ ethcommon.Address, _ ethcommon.Hash) ([]byte, error) {
	return []byte{}, nil
}

func (r *MockChainStateReader) PendingCodeAt(_ context.Context, _ string, _ ethcommon.Address) ([]byte, error) {
	return []byte{}, nil
}

func (r *MockChainStateReader) PendingNonceAt(_ context.Context, endpoint string, _ ethcommon.Address) (uint64, error) {
	if endpoint == "0" {
		// Simulate error
		return 0, fmt.Errorf("unknown chain")
	}

	return 10, nil
}

type MockNonceManager struct {
	memory.NonceManager
}

func (nm *MockNonceManager) GetLastAttributed(key string) (value uint64, ok bool, err error) {
	if strings.Contains(key, "error-on-get") {
		// Simulate error
		return 0, false, fmt.Errorf("could not get nonce")
	}
	return nm.NonceManager.GetLastAttributed(key)
}

func (nm *MockNonceManager) SetLastAttributed(key string, value uint64) error {
	if strings.Contains(key, "error-on-set") {
		// Simulate error
		return fmt.Errorf("could not set nonce")
	}
	_ = nm.NonceManager.SetLastAttributed(key, value)
	return nil
}

func (nm *MockNonceManager) IncrLastAttributed(key string) error {
	if strings.Contains(key, "error-on-incr") {
		// Simulate error
		return fmt.Errorf("could not increment nonce")
	}
	_ = nm.NonceManager.IncrLastAttributed(key)
	return nil
}

type mockMsg string

func (m mockMsg) Entrypoint() string    { return "" }
func (m mockMsg) Value() []byte         { return []byte{} }
func (m mockMsg) Key() []byte           { return []byte(m) }
func (m mockMsg) Header() engine.Header { return &header{} }

type header struct{}

func (h *header) Add(_, _ string)     {}
func (h *header) Del(_ string)        {}
func (h *header) Get(_ string) string { return "" }
func (h *header) Set(_, _ string)     {}

func makeNonceContext(endpoint, key string, expectedNonce uint64, expectedErrorCount int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	txctx.In = mockMsg(key)
	_ = txctx.Builder.SetFrom(ethcommon.HexToAddress("0x1"))
	txctx.WithContext(proxy.With(txctx.Context(), endpoint))

	txctx.Set("expectedErrorCount", expectedErrorCount)
	txctx.Set("expectedNonce", expectedNonce)

	return txctx
}

func assertTxContext(t *testing.T, txctx *engine.TxContext) {
	assert.Len(t, txctx.Builder.GetErrors(), txctx.Get("expectedErrorCount").(int), "Error count should be correct")
	assert.Equal(t, txctx.Get("expectedNonce").(uint64), txctx.Builder.MustGetNonceUint64(), "Nonce should be correct")
}

func TestNonceHandler(t *testing.T) {
	m := memory.NewNonceManager()
	nm := &MockNonceManager{*m}
	h := Nonce(nm, &MockChainStateReader{})

	testKey1 := "key1"
	// On 1st execution nonce should be 10 (as the mock client returns always return pending nonce 10)
	txctx := makeNonceContext("1", testKey1, 10, 0)
	h(txctx)
	assertTxContext(t, txctx)

	// On 2nd execution nonce should be 11 (as nonce should be retrieved from cache)
	txctx = makeNonceContext("1", testKey1, 11, 0)
	h(txctx)
	assertTxContext(t, txctx)

	// On 3rd execution we signal a recovery from 5 so expected nonce should be 5
	txctx = makeNonceContext("1", testKey1, 5, 0)
	_ = txctx.Builder.SetInternalLabelsValue("nonce.recovering.expected", "5")
	h(txctx)
	assertTxContext(t, txctx)

	// NonceManager should trigger an error get
	txctx = makeNonceContext("1", "key-error-on-get", 0, 1)
	h(txctx)
	assertTxContext(t, txctx)

	// NonceManager should trigger an error on set
	txctx = makeNonceContext("1", "key-error-on-set", 10, 0)
	h(txctx)
	assertTxContext(t, txctx)

	// NonceManager should error when unknown chain
	txctx = makeNonceContext("0", "key", 0, 1)
	h(txctx)
	assertTxContext(t, txctx)
}
