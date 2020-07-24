// +build unit

package storer

import (
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client/mock"
	"testing"
)

func makeContext(hash, id, endpoint string, expectedErrors int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	txctx.WithContext(proxy.With(txctx.Context(), endpoint))
	_ = txctx.Envelope.SetID(id).SetTxHashString(hash)
	txctx.Set("expectedErrors", expectedErrors)
	return txctx
}

type mockHandler struct {
	callCount int
}

func (h *mockHandler) Handle(txctx *engine.TxContext) {
	h.callCount++
}

func TestTxAlreadySent(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockChainLedgerReader := mock2.NewMockChainLedgerReader(ctrl)
	txSchedulerClient := mock.NewMockTransactionSchedulerClient(ctrl)
	mh := mockHandler{}

	// Prepare a test handler combined with a mock handler to
	// control abort are occurring as expected
	handler := engine.CombineHandlers(TxAlreadySent(mockChainLedgerReader, txSchedulerClient), mh.Handle)

	t.Run("Envelope should be sent if status is CREATED", func(t *testing.T) {
		txctx := makeContext("0x7a34cbb73c02aa3309c343e9e9b35f2a992aaa623c2ec2524816f476c63d2efa", "1", "8", 0)
		_ = txctx.Envelope.SetContextLabelsValue("jobUUID", txctx.Envelope.GetID())
		jobResponse := testutils.FakeJobResponse()
		jobResponse.Status = utils.StatusCreated

		txSchedulerClient.EXPECT().GetJob(txctx.Context(), txctx.Envelope.GetID()).Return(jobResponse, nil)

		handler(txctx)

		assert.Empty(t, txctx.Envelope.Errors)
	})

	t.Run("Envelope should be sent if status is STARTED", func(t *testing.T) {
		txctx := makeContext("0x7a34cbb73c02aa3309c343e9e9b35f2a992aaa623c2ec2524816f476c63d2efa", "1", "8", 0)
		_ = txctx.Envelope.SetContextLabelsValue("jobUUID", txctx.Envelope.GetID())
		jobResponse := testutils.FakeJobResponse()
		jobResponse.Status = utils.StatusStarted

		txSchedulerClient.EXPECT().GetJob(txctx.Context(), txctx.Envelope.GetID()).Return(jobResponse, nil)

		handler(txctx)

		assert.Empty(t, txctx.Envelope.Errors)
	})

	t.Run("Envelope should be sent if status is RECOVERING", func(t *testing.T) {
		txctx := makeContext("0x7a34cbb73c02aa3309c343e9e9b35f2a992aaa623c2ec2524816f476c63d2efa", "1", "8", 0)
		_ = txctx.Envelope.SetContextLabelsValue("jobUUID", txctx.Envelope.GetID())
		jobResponse := testutils.FakeJobResponse()
		jobResponse.Status = utils.StatusRecovering

		txSchedulerClient.EXPECT().GetJob(txctx.Context(), txctx.Envelope.GetID()).Return(jobResponse, nil)

		handler(txctx)

		assert.Empty(t, txctx.Envelope.Errors)
	})

	t.Run("should abort if status is FAILED", func(t *testing.T) {
		txctx := makeContext("0x7a34cbb73c02aa3309c343e9e9b35f2a992aaa623c2ec2524816f476c63d2efa", "1", "8", 0)
		_ = txctx.Envelope.SetContextLabelsValue("jobUUID", txctx.Envelope.GetID())
		jobResponse := testutils.FakeJobResponse()
		jobResponse.Status = utils.StatusFailed

		txSchedulerClient.EXPECT().GetJob(txctx.Context(), txctx.Envelope.GetID()).Return(jobResponse, nil)

		handler(txctx)

		assert.True(t, txctx.Envelope.OnlyWarnings())
	})

	t.Run("should abort if status is MINED", func(t *testing.T) {
		txctx := makeContext("0x7a34cbb73c02aa3309c343e9e9b35f2a992aaa623c2ec2524816f476c63d2efa", "1", "8", 0)
		_ = txctx.Envelope.SetContextLabelsValue("jobUUID", txctx.Envelope.GetID())
		jobResponse := testutils.FakeJobResponse()
		jobResponse.Status = utils.StatusFailed

		txSchedulerClient.EXPECT().GetJob(txctx.Context(), txctx.Envelope.GetID()).Return(jobResponse, nil)

		handler(txctx)

		assert.True(t, txctx.Envelope.OnlyWarnings())
	})

	t.Run("should abort if status is MINED", func(t *testing.T) {
		txctx := makeContext("0x7a34cbb73c02aa3309c343e9e9b35f2a992aaa623c2ec2524816f476c63d2efa", "1", "8", 0)
		_ = txctx.Envelope.SetContextLabelsValue("jobUUID", txctx.Envelope.GetID())
		jobResponse := testutils.FakeJobResponse()
		jobResponse.Status = utils.StatusFailed

		txSchedulerClient.EXPECT().GetJob(txctx.Context(), txctx.Envelope.GetID()).Return(jobResponse, nil)

		handler(txctx)

		assert.True(t, txctx.Envelope.OnlyWarnings())
	})

	t.Run("should get status from node if status is PENDING and succeed if not found", func(t *testing.T) {
		txctx := makeContext("0x7a34cbb73c02aa3309c343e9e9b35f2a992aaa623c2ec2524816f476c63d2efa", "1", "8", 0)
		_ = txctx.Envelope.SetContextLabelsValue("jobUUID", txctx.Envelope.GetID())
		jobResponse := testutils.FakeJobResponse()
		jobResponse.Status = utils.StatusPending

		txSchedulerClient.EXPECT().GetJob(txctx.Context(), txctx.Envelope.GetID()).Return(jobResponse, nil)
		mockChainLedgerReader.EXPECT().
			TransactionByHash(txctx.Context(), gomock.Any(), jobResponse.Transaction.GetHash()).
			Return(nil, false, nil)

		handler(txctx)

		assert.Empty(t, txctx.Envelope.Errors)
	})

	t.Run("should abort if error returned from node", func(t *testing.T) {
		txctx := makeContext("0x7a34cbb73c02aa3309c343e9e9b35f2a992aaa623c2ec2524816f476c63d2efa", "1", "8", 0)
		_ = txctx.Envelope.SetContextLabelsValue("jobUUID", txctx.Envelope.GetID())
		jobResponse := testutils.FakeJobResponse()
		jobResponse.Status = utils.StatusPending

		txSchedulerClient.EXPECT().GetJob(txctx.Context(), txctx.Envelope.GetID()).Return(jobResponse, nil)
		mockChainLedgerReader.EXPECT().
			TransactionByHash(txctx.Context(), gomock.Any(), jobResponse.Transaction.GetHash()).
			Return(nil, false, fmt.Errorf(""))

		handler(txctx)

		assert.NotEmpty(t, txctx.Envelope.Errors)
	})

	t.Run("should abort if tx found", func(t *testing.T) {
		txctx := makeContext("0x7a34cbb73c02aa3309c343e9e9b35f2a992aaa623c2ec2524816f476c63d2efa", "1", "8", 0)
		_ = txctx.Envelope.SetContextLabelsValue("jobUUID", txctx.Envelope.GetID())
		jobResponse := testutils.FakeJobResponse()
		jobResponse.Status = utils.StatusPending

		txSchedulerClient.EXPECT().GetJob(txctx.Context(), txctx.Envelope.GetID()).Return(jobResponse, nil)
		mockChainLedgerReader.EXPECT().
			TransactionByHash(txctx.Context(), gomock.Any(), jobResponse.Transaction.GetHash()).
			Return(&ethtypes.Transaction{}, false, nil)

		handler(txctx)

		assert.True(t, txctx.Envelope.OnlyWarnings())
	})
}
