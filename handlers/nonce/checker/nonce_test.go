package noncechecker

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/nonce/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/ethereum"
)

type MockChainStateReader struct{}

func (r *MockChainStateReader) BalanceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (r *MockChainStateReader) StorageAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, key ethcommon.Hash, blockNumber *big.Int) ([]byte, error) {
	return []byte{}, nil
}

func (r *MockChainStateReader) CodeAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) ([]byte, error) {
	return []byte{}, nil
}

func (r *MockChainStateReader) NonceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) (uint64, error) {
	return 0, nil
}

func (r *MockChainStateReader) PendingBalanceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (r *MockChainStateReader) PendingStorageAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, key ethcommon.Hash) ([]byte, error) {
	return []byte{}, nil
}

func (r *MockChainStateReader) PendingCodeAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) ([]byte, error) {
	return []byte{}, nil
}

func (r *MockChainStateReader) PendingNonceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) (uint64, error) {
	if chainID.Text(10) == "0" {
		// Simulate error
		return 0, fmt.Errorf("unknown chain")
	}

	return 10, nil
}

type MockNonceManager struct {
	mock.NonceManager
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
	chainID int64,
	key string,
	nonce, expectedNonceInMetadata uint64,
	expectedErrorCount int,
) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	txctx.Envelope.From = &ethereum.Account{Raw: []byte{}}
	txctx.Envelope.Chain = chain.FromInt(chainID)
	txctx.Envelope.Tx = &ethereum.Transaction{TxData: &ethereum.TxData{
		Nonce: nonce,
	}}
	txctx.In = mockMsg(key)

	txctx.Set("expectedErrorCount", expectedErrorCount)
	txctx.Set("expectedNonceInMetadata", expectedNonceInMetadata)

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
}

func TestChecker(t *testing.T) {
	m := mock.NewNonceManager()
	nm := &MockNonceManager{*m}
	h := engine.CombineHandlers(
		RecoveryStatusSetter(nm),
		Checker(nm, &MockChainStateReader{}),
	)

	testKey1 := "key1"
	// On 1st execution envelope with nonce 10 should be valid (as the mock client returns always return pending nonce 10)
	txctx := makeContext(1, testKey1, 10, 0, 0)
	h(txctx)
	assertTxContext(t, txctx)

	// On 2nd execution envelope with nonce 11 should be valid nonce manager should have incremented last sent
	txctx = makeContext(1, testKey1, 11, 0, 0)
	h(txctx)
	assertTxContext(t, txctx)

	// On 3rd execution envelope with nonce 10 should be too low
	txctx = makeContext(1, testKey1, 10, 0, 0)
	h(txctx)
	assertTxContext(t, txctx)

	// On 4th execution envelope with nonce 14 should be too high
	// Checker should signal in metadata
	txctx = makeContext(1, testKey1, 14, 12, 0)
	h(txctx)
	assertTxContext(t, txctx)
	recovering, _ := nm.IsRecovering(testKey1)
	assert.True(t, recovering, "NonceManager should be recovering")

	// On 5th execution envelope with nonce 15 should be too high
	// Checker should not signal in metadata
	txctx = makeContext(1, testKey1, 15, 0, 0)
	h(txctx)
	assertTxContext(t, txctx)

	// On 6th execution envelope with nonce 12 be valid
	txctx = makeContext(1, testKey1, 12, 0, 0)
	h(txctx)
	assertTxContext(t, txctx)
	recovering, _ = nm.IsRecovering(testKey1)
	assert.False(t, recovering, "NonceManager should have stopped recovering")

	// Execution with invalid chain
	txctx = makeContext(0, "key2", 10, 0, 1)
	h(txctx)
	assertTxContext(t, txctx)

	// Execution with error on nonce manager
	txctx = makeContext(1, "key-error-on-get", 10, 0, 1)
	h(txctx)
	assertTxContext(t, txctx)

	// Execution with error on nonce manager
	txctx = makeContext(1, "key-error-on-set", 10, 0, 0)
	h(txctx)
	assertTxContext(t, txctx)

	// Execution with error on nonce manager
	txctx = makeContext(1, "key-error-on-recovering-is", 10, 0, 0)
	h(txctx)
	assertTxContext(t, txctx)

	// Execution with error on nonce manager
	txctx = makeContext(1, "key-error-on-recovering-set", 10, 0, 0)
	h(txctx)
	assertTxContext(t, txctx)
}
