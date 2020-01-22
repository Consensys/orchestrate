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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

type MockChainStateReader struct{}

func (r *MockChainStateReader) BalanceAt(ctx context.Context, endpoint string, account ethcommon.Address, blockNumber *big.Int) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (r *MockChainStateReader) StorageAt(ctx context.Context, endpoint string, account ethcommon.Address, key ethcommon.Hash, blockNumber *big.Int) ([]byte, error) {
	return []byte{}, nil
}

func (r *MockChainStateReader) CodeAt(ctx context.Context, endpoint string, account ethcommon.Address, blockNumber *big.Int) ([]byte, error) {
	return []byte{}, nil
}

func (r *MockChainStateReader) NonceAt(ctx context.Context, endpoint string, account ethcommon.Address, blockNumber *big.Int) (uint64, error) {
	return 0, nil
}

func (r *MockChainStateReader) PendingBalanceAt(ctx context.Context, endpoint string, account ethcommon.Address) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (r *MockChainStateReader) PendingStorageAt(ctx context.Context, endpoint string, account ethcommon.Address, key ethcommon.Hash) ([]byte, error) {
	return []byte{}, nil
}

func (r *MockChainStateReader) PendingCodeAt(ctx context.Context, endpoint string, account ethcommon.Address) ([]byte, error) {
	return []byte{}, nil
}

func (r *MockChainStateReader) PendingNonceAt(ctx context.Context, endpoint string, account ethcommon.Address) (uint64, error) {
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

func (h *header) Add(key, value string) {}
func (h *header) Del(key string)        {}
func (h *header) Get(key string) string { return "" }
func (h *header) Set(key, value string) {}

func makeNonceContext(endpoint, key string, expectedNonce uint64, expectedErrorCount int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	txctx.Envelope.From = &ethereum.Account{Raw: []byte{}}
	txctx.Envelope.Tx = &ethereum.Transaction{TxData: &ethereum.TxData{}}
	txctx.In = mockMsg(key)
	txctx.WithContext(proxy.With(txctx.Context(), endpoint))

	txctx.Set("expectedErrorCount", expectedErrorCount)
	txctx.Set("expectedNonce", expectedNonce)

	return txctx
}

func assertTxContext(t *testing.T, txctx *engine.TxContext) {
	assert.Len(t, txctx.Envelope.GetErrors(), txctx.Get("expectedErrorCount").(int), "Error count should be correct")
	assert.Equal(t, txctx.Get("expectedNonce").(uint64), txctx.Envelope.GetTx().GetTxData().GetNonce(), "Nonce should be correct")
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
	txctx.Envelope.Metadata = &envelope.Metadata{Extra: map[string]string{"nonce.recovering.expected": "5"}}
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
