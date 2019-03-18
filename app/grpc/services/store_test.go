package services

import (
	"context"
	"math/big"
	"testing"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	store "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/context-store"
	ethereum "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/ethereum"
	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

type StoreServiceTestSuite struct {
	suite.Suite
	store *StoreService
}

func (suite *StoreServiceTestSuite) SetupTest() {
	suite.store = &StoreService{store: mock.NewTraceStore()}
}

func (suite *StoreServiceTestSuite) TestStore() {
	req := &store.StoreRequest{Trace: testTrace()}
	_, err := suite.store.Store(context.Background(), req)
	assert.Nil(suite.T(), err, "Store should not error")
}

func (suite *StoreServiceTestSuite) TestLoadByTxHash() {
	// Stores a trace
	suite.store.Store(context.Background(), &store.StoreRequest{Trace: testTrace()})

	req := &store.TxHashRequest{
		ChainId: "0x3",
		TxHash:  "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210",
	}
	resp, err := suite.store.LoadByTxHash(context.Background(), req)
	assert.Nil(suite.T(), err, "LoadByTxHash should not error")
	assert.Equal(suite.T(), "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11", resp.GetTrace().GetMetadata().GetId(), "Trace should be correctly loaded")
}

func (suite *StoreServiceTestSuite) TestLoadByTraceID() {
	// Stores a trace
	suite.store.Store(context.Background(), &store.StoreRequest{Trace: testTrace()})

	req := &store.TraceIDRequest{
		TraceId: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
	}
	resp, err := suite.store.LoadByTraceID(context.Background(), req)
	assert.Nil(suite.T(), err, "LoadByTraceID should not error")
	assert.Equal(suite.T(), "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210", resp.GetTrace().GetTx().GetHash(), "Trace should be correctly loaded")
}

func (suite *StoreServiceTestSuite) TestGetStatus() {
	// Stores a trace
	suite.store.Store(context.Background(), &store.StoreRequest{Trace: testTrace()})

	req := &store.TraceIDRequest{
		TraceId: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
	}
	resp, err := suite.store.GetStatus(context.Background(), req)
	assert.Nil(suite.T(), err, "GetStatus should not error")
	assert.Equal(suite.T(), "stored", resp.GetStatus(), "Status should be correct")
}

func (suite *StoreServiceTestSuite) TestSetStatus() {
	// Stores a trace
	suite.store.Store(context.Background(), &store.StoreRequest{Trace: testTrace()})

	req := &store.SetStatusRequest{
		TraceId: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
		Status:  "pending",
	}

	_, err := suite.store.SetStatus(context.Background(), req)
	assert.Nil(suite.T(), err, "SetStatus should not error")

	resp, err := suite.store.GetStatus(context.Background(), &store.TraceIDRequest{
		TraceId: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
	})
	assert.Nil(suite.T(), err, "GetStatus should not error")
	assert.Equal(suite.T(), "pending", resp.GetStatus(), "Status should be correct")
}

func (suite *StoreServiceTestSuite) TestLoadPendingTraces() {
	// Stores a trace and set its status to pending
	suite.store.Store(context.Background(), &store.StoreRequest{Trace: testTrace()})
	suite.store.SetStatus(context.Background(), &store.SetStatusRequest{
		TraceId: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
		Status:  "pending",
	})
	time.Sleep(100 * time.Millisecond)

	req := &store.PendingTracesRequest{
		Duration: (50 * time.Millisecond).Nanoseconds(),
	}

	resp, err := suite.store.LoadPendingTraces(context.Background(), req)
	assert.Nil(suite.T(), err, "LoadPendingTraces should not error")
	assert.Len(suite.T(), resp.Traces, 1, "Expect pending traces")

	req = &store.PendingTracesRequest{
		Duration: (200 * time.Millisecond).Nanoseconds(),
	}
	resp, err = suite.store.LoadPendingTraces(context.Background(), req)
	assert.Nil(suite.T(), err, "LoadPendingTraces should not error")
	assert.Len(suite.T(), resp.Traces, 0, "Expect no pending traces")
}

func testTrace() *trace.Trace {
	txData := (&ethereum.TxData{}).
		SetNonce(10).
		SetTo(ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C")).
		SetValue(big.NewInt(100000)).
		SetGas(2000).
		SetGasPrice(big.NewInt(200000)).
		SetData(hexutil.MustDecode("0xabcd"))

	return &trace.Trace{
		Chain:    &common.Chain{Id: "0x3"},
		Metadata: &trace.Metadata{Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11"},
		Tx: &ethereum.Transaction{
			TxData: txData,
			Raw:    "0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
			Hash:   "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210",
		},
	}
}

func TestTodoTestSuite(t *testing.T) {
	suite.Run(t, new(StoreServiceTestSuite))
}
