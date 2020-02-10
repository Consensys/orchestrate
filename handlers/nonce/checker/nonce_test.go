package noncechecker

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"

	ethcommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/nonce/memory"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

type MockChainStateReader struct{}

const endpointError = "error"

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
	if endpoint == endpointError {
		// Simulate error
		return 0, fmt.Errorf("unknown chain")
	}

	return 10, nil
}

type MockNonceManager struct {
	memory.NonceManager
}

func (nm *MockNonceManager) GetLastSent(key string) (value uint64, ok bool, err error) {
	if strings.Contains(key, "error-on-get") {
		// Simulate error
		return 0, false, fmt.Errorf("could not get nonce")
	}
	return nm.NonceManager.GetLastSent(key)
}

func (nm *MockNonceManager) SetLastSent(key string, value uint64) error {
	if strings.Contains(key, "error-on-set") {
		// Simulate error
		return fmt.Errorf("could not set nonce")
	}
	_ = nm.NonceManager.SetLastSent(key, value)
	return nil
}

func (nm *MockNonceManager) IncrLastSent(key string) error {
	if strings.Contains(key, "error-on-incr") {
		// Simulate error
		return fmt.Errorf("could not increment nonce")
	}
	_ = nm.NonceManager.IncrLastSent(key)
	return nil
}

func (nm *MockNonceManager) IsRecovering(key string) (bool, error) {
	if strings.Contains(key, "error-on-recovering-is") {
		// Simulate error
		return false, fmt.Errorf("coult not load recovery status")
	}
	return nm.NonceManager.IsRecovering(key)
}

func (nm *MockNonceManager) SetRecovering(key string, status bool) error {
	if strings.Contains(key, "error-on-recovering-ser") {
		// Simulate error
		return fmt.Errorf("could not set recovery status")
	}
	return nm.NonceManager.SetRecovering(key, status)
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

func makeContext(
	endpoint,
	key string,
	invalid bool,
	nonce, expectedNonceInMetadata uint64,
	expectedRecoveryCount, expectedErrorCount int,
	errorOnSend string,
) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	txctx.Envelope.Tx = &ethereum.Transaction{TxData: &ethereum.TxData{
		Nonce: nonce,
	}}
	txctx.In = mockMsg(key)
	txctx.WithContext(proxy.With(txctx.Context(), endpoint))

	txctx.Set("expectedErrorCount", expectedErrorCount)
	txctx.Set("expectedNonceInMetadata", expectedNonceInMetadata)
	txctx.Set("expectedInvalid", invalid)
	txctx.Set("expectedRecoveryCount", expectedRecoveryCount)
	txctx.Set("errorOnSend", errorOnSend)

	return txctx
}

func assertTxContext(t *testing.T, txctx *engine.TxContext) {
	assert.Len(t, txctx.Envelope.GetErrors(), txctx.Get("expectedErrorCount").(int), "Error count should be correct")

	expectedNonceInMetadata := txctx.Get("expectedNonceInMetadata").(uint64)
	if expectedNonceInMetadata > 0 {
		v, ok := txctx.Envelope.GetMetadataValue("nonce.recovering.expected")
		assert.True(t, ok, "Signal for nonce recovery in envelope metadata should have been set")
		lastSent, _ := strconv.ParseUint(v, 10, 64)
		assert.Equal(t, expectedNonceInMetadata, lastSent, "Nonce in metadata should be correct")
	} else {
		_, ok := txctx.Envelope.GetMetadataValue("nonce.recovering.expected")
		assert.False(t, ok, "Signal for nonce recovery in envelope metadata should not have been set")
	}

	if txctx.Get("expectedInvalid").(bool) {
		v, ok := txctx.Get("invalid.nonce").(bool)
		assert.True(t, ok && v, "Nonce invalidity should be correct")
	} else {
		invalid, ok := txctx.Get("invalid.nonce").(bool)
		assert.False(t, ok || invalid, "Nonce invalidity should be correct")
	}

	var recoveryCount int
	v, ok := txctx.Envelope.GetMetadataValue("nonce.recovering.count")
	if ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		recoveryCount = i
	}
	assert.Equal(t, txctx.Get("expectedRecoveryCount").(int), recoveryCount, "Recovery count should be")
}

func MockSenderHandler(txctx *engine.TxContext) {
	if v, ok := txctx.Get("errorOnSend").(string); ok && v != "" {
		_ = txctx.Error(fmt.Errorf(v))
	}
}

func TestChecker(t *testing.T) {
	m := memory.NewNonceManager()
	nm := &MockNonceManager{*m}
	tracker := NewRecoveryTracker()
	conf := &Configuration{
		MaxRecovery: 5,
	}
	h := engine.CombineHandlers(
		RecoveryStatusSetter(nm, tracker),
		Checker(conf, nm, &MockChainStateReader{}, tracker),
		MockSenderHandler,
	)

	testKey1 := "key1"
	// On 1st execution envelope with nonce 10 should be valid (as the mock client returns always return pending nonce 10)
	txctx := makeContext("testURL", testKey1, false, 10, 0, 0, 0, "")
	h(txctx)
	assertTxContext(t, txctx)

	// On 2nd execution envelope with nonce 11 should be valid nonce manager should have incremented last sent
	txctx = makeContext("testURL", testKey1, false, 11, 0, 0, 0, "")
	h(txctx)
	assertTxContext(t, txctx)

	// On 3rd execution envelope with nonce 10 should be too low
	txctx = makeContext("testURL", testKey1, true, 10, 0, 1, 0, "")
	h(txctx)
	assertTxContext(t, txctx)

	// On 4th execution envelope with nonce 14 should be too high
	// Checker should signal in metadata
	txctx = makeContext("testURL", testKey1, true, 14, 12, 1, 0, "")
	h(txctx)
	assertTxContext(t, txctx)
	recovering := tracker.Recovering(testKey1) > 0
	assert.True(t, recovering, "NonceManager should be recovering")

	// On 5th execution envelope with nonce 15 should be too high
	// Checker should not signal in metadata
	txctx = makeContext("testURL", testKey1, true, 15, 0, 3, 0, "")
	txctx.Envelope.SetMetadataValue("nonce.recovering.count", fmt.Sprintf("%v", 2))
	h(txctx)
	assertTxContext(t, txctx)

	// On 6th execution envelope with nonce 12 be valid
	txctx = makeContext("testURL", testKey1, false, 12, 0, 0, 0, "")
	h(txctx)
	assertTxContext(t, txctx)
	recovering = tracker.Recovering(testKey1) > 0
	assert.False(t, recovering, "NonceManager should have stopped recovering")

	// On 7th execution envelope with nonce 14 but raw mode should be valid
	txctx = makeContext("testURL", testKey1, false, 14, 0, 0, 0, "")
	txctx.Envelope.SetMetadataValue("tx.mode", "raw")

	h(txctx)
	assertTxContext(t, txctx)
	recovering = tracker.Recovering(testKey1) > 0
	assert.False(t, recovering, "NonceManager should not be recovering")

	// Execution with invalid chain
	txctx = makeContext(endpointError, "key2", false, 10, 0, 0, 1, "")
	h(txctx)
	assertTxContext(t, txctx)

	// Execution with error on nonce manager
	txctx = makeContext("testURL", "key-error-on-get", false, 10, 0, 0, 1, "")
	h(txctx)
	assertTxContext(t, txctx)

	// Execution with error on nonce manager
	txctx = makeContext("testURL", "key-error-on-set", false, 10, 0, 0, 0, "")
	h(txctx)
	assertTxContext(t, txctx)

	// Execution with nonce too low on send
	txctx = makeContext("testURL", testKey1, true, 13, 0, 1, 0, "json-rpc: nonce too low")
	h(txctx)
	assertTxContext(t, txctx)
	v, _, _ := nm.GetLastSent(testKey1)
	assert.Equal(t, uint64(9), v, "Nonce should have been re-initialized")

	// Execution with recovery count exceeded
	txctx = makeContext("testURL", testKey1, false, 13, 0, 10, 1, "")
	txctx.Envelope.SetMetadataValue("nonce.recovering.count", fmt.Sprintf("%v", 10))
	h(txctx)
	assertTxContext(t, txctx)
}
