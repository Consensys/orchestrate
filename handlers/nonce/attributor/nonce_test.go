// +build unit

package nonceattributor

import (
	"fmt"
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

func (nm *MockNonceManager) GetLastAttributed(key string) (value uint64, ok bool, err error) {
	if strings.Contains(key, "400") {
		// Simulate error
		return 0, false, fmt.Errorf("could not get nonce")
	}
	return nm.NonceManager.GetLastAttributed(key)
}

func (nm *MockNonceManager) SetLastAttributed(key string, value uint64) error {
	if strings.Contains(key, "404") {
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

func makeNonceContext(endpoint, chainID string, jobType tx.JobType, expectedNonce uint64, expectedErrorCount int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	_ = txctx.Envelope.
		SetFrom(ethcommon.HexToAddress("0x1")).
		SetJobType(jobType).
		SetChainIDString(chainID)
	txctx.WithContext(proxy.With(txctx.Context(), endpoint))

	txctx.Set("expectedErrorCount", expectedErrorCount)
	txctx.Set("expectedNonce", expectedNonce)

	return txctx
}

func assertTxContext(t *testing.T, txctx *engine.TxContext) {
	assert.Len(t, txctx.Envelope.GetErrors(), txctx.Get("expectedErrorCount").(int), "Error count should be correct")
	assert.Equal(t, txctx.Get("expectedNonce").(uint64), txctx.Envelope.MustGetNonceUint64(), "Nonce should be correct")
}

func TestNonceHandler(t *testing.T) {
	m := memory.NewNonceManager()
	nm := &MockNonceManager{*m}
	ctrl := gomock.NewController(t)
	ec := mocks.NewMockEthClient(ctrl)
	ec.EXPECT().PendingNonceAt(gomock.Any(), gomock.Eq(endpointError), gomock.Any()).
		Return(uint64(0), fmt.Errorf("unknown chain")).
		AnyTimes()
	ec.EXPECT().PendingNonceAt(gomock.Any(), gomock.Not(gomock.Eq(endpointError)), gomock.Any()).
		Return(uint64(10), nil).
		AnyTimes()

	h := Nonce(nm, ec)

	testKey1 := "42"
	// On 1st execution nonce should be 10 (as the mock client returns always return pending nonce 10)
	txctx := makeNonceContext("1", testKey1, tx.JobType_ETH_TX, 10, 0)
	h(txctx)
	assertTxContext(t, txctx)

	// On 2nd execution nonce should be 11 (as nonce should be retrieved from cache)
	txctx = makeNonceContext("1", testKey1, tx.JobType_ETH_TX, 11, 0)
	h(txctx)
	assertTxContext(t, txctx)

	// On 3rd execution we signal a recovery from 5 so expected nonce should be 5
	txctx = makeNonceContext("1", testKey1, tx.JobType_ETH_TX, 5, 0)
	_ = txctx.Envelope.SetInternalLabelsValue("nonce.recovering.expected", "5")
	h(txctx)
	assertTxContext(t, txctx)

	// NonceManager should trigger an error get
	txctx = makeNonceContext("1", "400", tx.JobType_ETH_TX, 0, 1)
	h(txctx)
	assertTxContext(t, txctx)

	// NonceManager should trigger an error on set
	txctx = makeNonceContext("1", "404", tx.JobType_ETH_TX, 10, 0)
	h(txctx)
	assertTxContext(t, txctx)

	// NonceManager should error when unknown chain
	txctx = makeNonceContext(endpointError, "key", tx.JobType_ETH_TX, 0, 1)
	h(txctx)
	assertTxContext(t, txctx)
}

func TestEEANonceHandler(t *testing.T) {
	m := memory.NewNonceManager()
	nm := &MockNonceManager{*m}
	ctrl := gomock.NewController(t)
	ec := mocks.NewMockEthClient(ctrl)
	ec.EXPECT().PrivNonce(gomock.Any(), gomock.Not(gomock.Eq(endpointError)), gomock.Any(), gomock.Any()).
		Return(uint64(10), nil).
		AnyTimes()

	handler := Nonce(nm, ec)

	t.Run("should execute PrivEEANonce without errors", func(t *testing.T) {
		expectedNonce := uint64(10)
		txctx := makeNonceContext("endpoint", "1000", tx.JobType_ETH_ORION_EEA_TX, expectedNonce, 0)
		txctx.Envelope.PrivateFor = []string{"PrivateFor"}
		ec.EXPECT().PrivEEANonce(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(expectedNonce, nil)
	
		handler(txctx)
	
		assert.Len(t, txctx.Envelope.GetErrors(),
			txctx.Get("expectedErrorCount").(int), "Error count should be correct")
		assert.Equal(t, txctx.Get("expectedNonce").(uint64), txctx.Envelope.MustGetNonceUint64(),
			"Nonce should be correct")
	})
	
	t.Run("should execute PrivEEANonce without errors", func(t *testing.T) {
		expectedNonce := uint64(10)
		txctx := makeNonceContext("endpoint", "1000", tx.JobType_ETH_ORION_EEA_TX, expectedNonce, 0)
		txctx.Envelope.PrivacyGroupID = "PrivacyGroupID"
		ec.EXPECT().PrivNonce(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(expectedNonce, nil)

		handler(txctx)

		assert.Len(t, txctx.Envelope.GetErrors(),
			txctx.Get("expectedErrorCount").(int), "Error count should be correct")
		assert.Equal(t, txctx.Get("expectedNonce").(uint64), txctx.Envelope.MustGetNonceUint64(),
			"Nonce should be correct")
	})
}
