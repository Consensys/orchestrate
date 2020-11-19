// +build unit

package noncechecker

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/nonce/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/nonce/memory"
)

const endpointError = "error"

type MockNonceManager struct {
	memory.NonceManager
}

func (nm *MockNonceManager) GetLastSent(key string) (value uint64, ok bool, err error) {
	if strings.Contains(key, "@400") {
		// Simulate error
		return 0, false, fmt.Errorf("could not get nonce")
	}
	return nm.NonceManager.GetLastSent(key)
}

func (nm *MockNonceManager) SetLastSent(key string, value uint64) error {
	if strings.Contains(key, "@404") {
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

func (h *header) Add(_, _ string)     {}
func (h *header) Del(_ string)        {}
func (h *header) Get(_ string) string { return "" }
func (h *header) Set(_, _ string)     {}

func makeContext(
	endpoint,
	chainID string,
	invalid bool,
	nonce, expectedNonceInMetadata uint64,
	expectedRecoveryCount, expectedErrorCount int,
	errorOnSend string,
) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	_ = txctx.Envelope.
		SetFrom(ethcommon.HexToAddress("0x1")).
		SetNonce(nonce).
		SetChainIDString(chainID)

	txctx.WithContext(proxy.With(txctx.Context(), endpoint))

	txctx.Set("expectedErrorCount", expectedErrorCount)
	txctx.Set("expectedNonceInMetadata", expectedNonceInMetadata)
	txctx.Set("expectedInvalid", invalid)
	txctx.Set("expectedRecoveryCount", expectedRecoveryCount)
	txctx.Set("errorOnSend", errorOnSend)

	return txctx
}

func  assertTxContext(t *testing.T, txctx *engine.TxContext) {
 	assert.Len(t, txctx.Envelope.GetErrors(), txctx.Get("expectedErrorCount").(int), "Error count should be correct")

	expectedNonceInMetadata := txctx.Get("expectedNonceInMetadata").(uint64)
	if expectedNonceInMetadata > 0 {
		v := txctx.Envelope.GetInternalLabelsValue("nonce.recovering.expected")
		assert.NotNil(t, v, "Signal for nonce recovery in envelope metadata should have been set")
		lastSent, _ := strconv.ParseUint(v, 10, 64)
		assert.Equal(t, expectedNonceInMetadata, lastSent, "Nonce in metadata should be correct")
	} else {
		v := txctx.Envelope.GetInternalLabelsValue("nonce.recovering.expected")
		assert.Empty(t, v, "Signal for nonce recovery in envelope metadata should not have been set")
	}

	if txctx.Get("expectedInvalid").(bool) {
		assert.True(t, txctx.HasInvalidNonceErr(), "Nonce invalidity should be correct")
	} else {
		assert.False(t, txctx.HasInvalidNonceErr(), "Nonce invalidity should be correct")
	}

	var recoveryCount int
	v := txctx.Envelope.GetInternalLabelsValue("nonce.recovering.count")
	if v != "" {
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
	ctrl := gomock.NewController(t)
	ec := mocks.NewMockEthClient(ctrl)
	ec.EXPECT().PendingNonceAt(gomock.Any(), gomock.Eq(endpointError), gomock.Any()).Return(uint64(0), fmt.Errorf("unknown chain")).AnyTimes()
	ec.EXPECT().PendingNonceAt(gomock.Any(), gomock.Not(gomock.Eq(endpointError)), gomock.Any()).Return(uint64(10), nil).AnyTimes()

	h := engine.CombineHandlers(
		RecoveryStatusSetter(nm, tracker),
		Checker(conf, nm, ec, tracker),
		MockSenderHandler,
	)

	chainID1 := "42"
	// On 1st execution envelope with nonce 10 should be valid (as the mock client returns always return pending nonce 10)
	txctx := makeContext("testURL", chainID1, false, 10, 0, 0, 0, "")
	h(txctx)
	assertTxContext(t, txctx)

	// On 2nd execution envelope with nonce 11 should be valid nonce manager should have incremented last sent
	txctx = makeContext("testURL", chainID1, false, 11, 0, 0, 0, "")
	h(txctx)
	assertTxContext(t, txctx)

	// On 3rd execution envelope with nonce 10 should be too low
	txctx = makeContext("testURL", chainID1, true, 10, 0, 1, 1, "")
	h(txctx)
	assertTxContext(t, txctx)

	// On 4th execution envelope with nonce 14 should be too high
	// Checker should signal in metadata
	txctx = makeContext("testURL", chainID1, true, 14, 12, 1, 1, "")
	h(txctx)
	assertTxContext(t, txctx)
	recovering := tracker.Recovering(txctx.Envelope.PartitionKey()) > 0
	assert.True(t, recovering, "NonceManager should be recovering")

	// On 5th execution envelope with nonce 15 should be too high
	// Checker should not signal in metadata
	txctx = makeContext("testURL", chainID1, true, 15, 0, 3, 1, "")
	_ = txctx.Envelope.SetInternalLabelsValue("nonce.recovering.count", fmt.Sprintf("%v", 2))
	h(txctx)
	assertTxContext(t, txctx)

	// On 6th execution envelope with nonce 12 be valid
	txctx = makeContext("testURL", chainID1, false, 12, 0, 0, 0, "")
	h(txctx)
	assertTxContext(t, txctx)
	recovering = tracker.Recovering(txctx.Envelope.PartitionKey()) > 0
	assert.False(t, recovering, "NonceManager should have stopped recovering")

	// On 7th execution envelope with nonce 14 but raw mode should be valid
	txctx = makeContext("testURL", chainID1, false, 15, 0, 0, 0, "")
	_ = txctx.Envelope.SetJobType(tx.JobType_ETH_RAW_TX)
	h(txctx)
	assertTxContext(t, txctx)
	recovering = tracker.Recovering(txctx.Envelope.PartitionKey()) > 0
	assert.False(t, recovering, "NonceManager should not be recovering")

	// Execution with invalid chain
	txctx = makeContext(endpointError, "12", false, 10, 0, 0, 1, "")
	h(txctx)
	assertTxContext(t, txctx)

	// Execution with error on nonce manager
	txctx = makeContext("testURL", "400", false, 10, 0, 0, 1, "")
	h(txctx)
	assertTxContext(t, txctx)

	// Execution with error on nonce manager
	txctx = makeContext("testURL", "404", false, 10, 0, 0, 0, "")
	h(txctx)
	assertTxContext(t, txctx)

	// Execution with nonce too low on send
	txctx = makeContext("testURL", chainID1, true, 13, 0, 1, 0, "json-rpc: nonce too low")
	h(txctx)
	assertTxContext(t, txctx)
	v, _, _ := nm.GetLastSent(txctx.Envelope.PartitionKey())
	assert.Equal(t, uint64(9), v, "Nonce should have been re-initialized")

	// Execution with recovery count exceeded
	txctx = makeContext("testURL", chainID1, false, 13, 0, 10, 1, "")
	_ = txctx.Envelope.SetInternalLabelsValue("nonce.recovering.count", fmt.Sprintf("%v", 10))
	h(txctx)
	assertTxContext(t, txctx)
}

func TestChecker_EEA(t *testing.T) {
	m := memory.NewNonceManager()
	nm := &MockNonceManager{*m}
	tracker := NewRecoveryTracker()
	conf := &Configuration{
		MaxRecovery: 5,
	}
	ctrl := gomock.NewController(t)
	ec := mocks.NewMockEthClient(ctrl)

	h := engine.CombineHandlers(
		RecoveryStatusSetter(nm, tracker),
		Checker(conf, nm, ec, tracker),
		MockSenderHandler,
	)

	t.Run("envelope for EEA private tx job and privateGroupID should be valid", func(t *testing.T) {
		// On 1nd execution envelope with nonce 0 should be fetch from the chain
		txctx := makeContext("testURL", "111", false, 0, 0, 0, 0, "")
		_ = txctx.Envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX).SetPrivacyGroupID("privateGroupID")
		ec.EXPECT().PrivNonce(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(uint64(0), nil)
		h(txctx)
		assertTxContext(t, txctx)
		
		// On 2nd execution envelope with nonce 1 should be valid from  incremented last sent
		txctx = makeContext("testURL", "111", false, 1, 0, 0, 0, "")
		_ = txctx.Envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX).SetPrivacyGroupID("privateGroupID")
		assertTxContext(t, txctx)
	})
	
	t.Run("envelope for EEA private tx job and privateGroupID should be invalid", func(t *testing.T) {
		// On 1nd execution envelope with nonce 0 should be fetch from the chain
		txctx := makeContext("testURL", "112", true, 0, 0, 1, 1, "")
		_ = txctx.Envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX).SetPrivacyGroupID("privateGroupID")
		ec.EXPECT().PrivNonce(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(uint64(1), nil)
		h(txctx)
		assertTxContext(t, txctx)
	})
	
	t.Run("envelope for EEA private tx job and privateGroupID should be valid", func(t *testing.T) {
		// On 1nd execution envelope with nonce 0 should be fetch from the chain
		txctx := makeContext("testURL", "111", false, 0, 0, 0, 0, "")
		_ = txctx.Envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX).SetPrivateFor([]string{"PrivateFor"})
		ec.EXPECT().PrivEEANonce(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(uint64(0), nil)
		h(txctx)
		assertTxContext(t, txctx)
		
		// On 2nd execution envelope with nonce 1 should be valid from  incremented last sent
		txctx = makeContext("testURL", "111", false, 1, 0, 0, 0, "")
		_ = txctx.Envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX).SetPrivateFor([]string{"PrivateFor"})
		assertTxContext(t, txctx)
	})
}
