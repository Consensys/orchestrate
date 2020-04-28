// +build unit

package storer

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	proto "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	clientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/client/mock"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

type MockChainLedgerReader struct {
	txs map[string]bool
}

func NewMockChainLedgerReader() *MockChainLedgerReader {
	return &MockChainLedgerReader{
		txs: make(map[string]bool),
	}
}

func (ec *MockChainLedgerReader) SendTx(hash string) {
	ec.txs[hash] = true
}

func (ec *MockChainLedgerReader) BlockByHash(ctx context.Context, endpoint string, hash ethcommon.Hash) (*ethtypes.Block, error) {
	return nil, fmt.Errorf("not implemented")
}

func (ec *MockChainLedgerReader) BlockByNumber(ctx context.Context, endpoint string, number *big.Int) (*ethtypes.Block, error) {
	return nil, fmt.Errorf("not implemented")
}

func (ec *MockChainLedgerReader) HeaderByHash(ctx context.Context, endpoint string, hash ethcommon.Hash) (*ethtypes.Header, error) {
	return nil, fmt.Errorf("not implemented")
}

func (ec *MockChainLedgerReader) HeaderByNumber(ctx context.Context, endpoint string, number *big.Int) (*ethtypes.Header, error) {
	return nil, fmt.Errorf("not implemented")
}

func (ec *MockChainLedgerReader) TransactionByHash(ctx context.Context, endpoint string, hash ethcommon.Hash) (*ethtypes.Transaction, bool, error) {
	if endpoint == "0" {
		return nil, false, fmt.Errorf("unknown chain")
	}
	_, ok := ec.txs[hash.Hex()]
	if ok {
		return &ethtypes.Transaction{}, false, nil
	}
	return nil, false, nil
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
func (ec *MockChainLedgerReader) TransactionReceipt(ctx context.Context, endpoint string, txHash ethcommon.Hash) (*proto.Receipt, error) {
	return nil, fmt.Errorf("not implemented")
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
func (ec *MockChainLedgerReader) PrivateTransactionReceipt(ctx context.Context, endpoint string, txHash ethcommon.Hash) (*proto.Receipt, error) {
	return nil, fmt.Errorf("not implemented")
}

func makeContext(hash, id, endpoint string, expectedErrors int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	txctx.WithContext(proxy.With(txctx.Context(), endpoint))
	_ = txctx.Envelope.SetID(id).SetTxHashString(hash)
	txctx.Set("expectedErrors", expectedErrors)
	return txctx
}

func assertCtx(t *testing.T, txctx *engine.TxContext) {
	assert.Len(t, txctx.Envelope.GetErrors(), txctx.Get("expectedErrors").(int), "Error count should be valid")
}

type mockHandler struct {
	callCount int
}

func (h *mockHandler) Handle(txctx *engine.TxContext) {
	h.callCount++
}

func TestTxAlreadySent(t *testing.T) {
	ec := NewMockChainLedgerReader()
	ctrl := gomock.NewController(t)
	client := clientmock.NewMockEnvelopeStoreClient(ctrl)
	client.EXPECT().LoadByID(gomock.Any(), gomock.AssignableToTypeOf(&svc.LoadByIDRequest{})).AnyTimes()
	client.EXPECT().Store(gomock.Any(), gomock.AssignableToTypeOf(&svc.StoreRequest{})).Times(2)
	client.EXPECT().SetStatus(gomock.Any(), gomock.AssignableToTypeOf(&svc.SetStatusRequest{})).Times(2)
	mh := mockHandler{}

	// Prepare a test handler combined with a mock handler to
	// control abort are occurring as expected
	handler := engine.CombineHandlers(
		TxAlreadySent(ec, client),
		mh.Handle,
	)

	// #1: First envelope should be send correctly and mock handler
	txctx := makeContext(
		"0x7a34cbb73c02aa3309c343e9e9b35f2a992aaa623c2ec2524816f476c63d2efa",
		"1",
		"8",
		0,
	)
	handler(txctx)
	assertCtx(t, txctx)
	assert.Equal(t, 1, mh.callCount, "Mock handler should been executed")

	// Store envelope, do not send transaction and set envelope status before handing context
	b := tx.NewEnvelope().SetID("2").SetChainID(big.NewInt(8)).SetTxHash(ethcommon.HexToHash("0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b"))
	_, _ = client.Store(
		context.Background(),
		&svc.StoreRequest{
			Envelope: b.TxEnvelopeAsRequest(),
		},
	)
	ec.SendTx("0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b")
	_, _ = client.SetStatus(
		context.Background(),
		&svc.SetStatusRequest{
			Id:     "2",
			Status: svc.Status_PENDING,
		},
	)
	txctx = makeContext(
		"0xf2beaddb2dc4e4c9055148a808365edbadd5f418c31631dcba9ad99af34ae66b",
		"2",
		"8",
		0,
	)
	handler(txctx)
	assertCtx(t, txctx)
	assert.Equal(t, 2, mh.callCount, "Mock handler should have been executed")

	// Store envelope, does not send transaction and set envelope status before handing context
	b = tx.NewEnvelope().SetID("3").SetChainID(big.NewInt(8)).SetTxHash(ethcommon.HexToHash("0x60a417c21da71cea33821071e99871fa2c23ad8103b889cf8a459b0b5320fd46"))
	_, _ = client.Store(
		context.Background(),
		&svc.StoreRequest{
			Envelope: b.TxEnvelopeAsRequest(),
		},
	)
	_, _ = client.SetStatus(
		context.Background(),
		&svc.SetStatusRequest{
			Id:     "3",
			Status: svc.Status_PENDING,
		},
	)
	txctx = makeContext(
		"0x60a417c21da71cea33821071e99871fa2c23ad8103b889cf8a459b0b5320fd46",
		"3",
		"8",
		0,
	)
	handler(txctx)
	assertCtx(t, txctx)
	assert.Equal(t, 3, mh.callCount, "Mock handler should have been executed")

}
