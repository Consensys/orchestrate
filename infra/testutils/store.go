package testutils

import (
	"context"
	"fmt"
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra"
	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	ethereum "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/ethereum"
	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

// TraceStoreTestSuite is a test suit for TraceStore
type TraceStoreTestSuite struct {
	suite.Suite
	Store infra.TraceStore
}

// TestTraceStore test trace store
func (suite *TraceStoreTestSuite) TestTraceStore() {
	txData := (&ethereum.TxData{}).
		SetNonce(10).
		SetTo(ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C")).
		SetValue(big.NewInt(100000)).
		SetGas(2000).
		SetGasPrice(big.NewInt(200000)).
		SetData(hexutil.MustDecode("0xabcd"))

	tr := &trace.Trace{
		Chain:    &common.Chain{Id: "0x3"},
		Metadata: &trace.Metadata{Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11"},
		Tx: &ethereum.Transaction{
			TxData: txData,
			Raw:    "0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
			Hash:   "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210",
		},
	}

	// Store Trace
	status, storedAt, err := suite.Store.Store(context.Background(), tr)
	assert.Nil(suite.T(), err, "Should properly store trace")
	assert.Equal(suite.T(), "stored", status, "Default status should be correct")
	assert.True(suite.T(), time.Now().Sub(storedAt) < 5*time.Second, "Stored date should be close")

	// Load Trace
	tr = &trace.Trace{}
	status, _, err = suite.Store.LoadByTxHash(context.Background(), "0x3", "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210", tr)
	assert.Nil(suite.T(), err, "Should properly store trace")
	assert.Equal(suite.T(), "stored", status, "Status should be correct")
	assert.Equal(suite.T(), "0x3", tr.GetChain().GetId(), "ChainID should be correct")
	assert.Equal(suite.T(), "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11", tr.GetMetadata().GetId(), "MetadataID should be correct")

	// Set Status
	err = suite.Store.SetStatus(context.Background(), "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11", "pending")
	assert.Nil(suite.T(), err, "Setting status to %q", "pending")

	status, sentAt, err := suite.Store.GetStatus(context.Background(), "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11")
	assert.Equal(suite.T(), "pending", status, "Status should be correct")
	assert.True(suite.T(), sentAt.Sub(storedAt) > 0, "Stored should be older than sent date")

	// Stores an already existing
	tr = &trace.Trace{
		Chain:    &common.Chain{Id: "0x3"},
		Metadata: &trace.Metadata{Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11"},
		Tx: &ethereum.Transaction{
			TxData: txData,
			Raw:    "0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
			Hash:   "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210",
		},
	}

	status, _, err = suite.Store.Store(context.Background(), tr)
	assert.Nil(suite.T(), err, "Should update")
	assert.Equal(suite.T(), "pending", status, "Status should be correct")

	// Set status to error
	err = suite.Store.SetStatus(context.Background(), "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11", "error")
	assert.Nil(suite.T(), err, "Setting status to %q", "error")

	status, errorAt, err := suite.Store.GetStatus(context.Background(), "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11")
	assert.Equal(suite.T(), "error", status, "Status should be correct")
	assert.True(suite.T(), errorAt.Sub(sentAt) > 0, "Stored date should be close")

	// Test to Load By ID
	status, _, err = suite.Store.LoadByTraceID(context.Background(), "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11", tr)
	assert.Equal(suite.T(), "error", status, "Status should be correct")
}

// TestLoadPendingTraces test load pending traces
func (suite *TraceStoreTestSuite) TestLoadPendingTraces() {
	txData := (&ethereum.TxData{}).
		SetNonce(10).
		SetTo(ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C")).
		SetValue(big.NewInt(100000)).
		SetGas(2000).
		SetGasPrice(big.NewInt(200000)).
		SetData(hexutil.MustDecode("0xabcd"))

	for i, chain := range []string{"0x1", "0x2", "0x3", "0xa2", "0x42", "0xab"} {
		tr := &trace.Trace{
			Chain:    &common.Chain{Id: chain},
			Metadata: &trace.Metadata{Id: fmt.Sprintf("a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a1%v", i)},
			Tx: &ethereum.Transaction{
				TxData: txData,
				Raw:    "0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
				Hash:   "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210",
			},
		}

		_, _, err := suite.Store.Store(context.Background(), tr)
		assert.Nil(suite.T(), err, "No error expected")
		time.Sleep(100 * time.Millisecond)

		if i%2 == 0 {
			err = suite.Store.SetStatus(context.Background(), fmt.Sprintf("a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a1%v", i), "pending")
			assert.Nil(suite.T(), err, "No error expected")
		}
	}

	traces, err := suite.Store.LoadPendingTraces(context.Background(), 0)
	assert.Nil(suite.T(), err, "No error expected on LoadPendingTraces")
	assert.Len(suite.T(), traces, 3, "Count of trace pending incorrect")

	traces, err = suite.Store.LoadPendingTraces(context.Background(), 300*time.Millisecond)
	assert.Nil(suite.T(), err, "No error expected on LoadPendingTraces")
	assert.Len(suite.T(), traces, 2, "Count of trace pending incorrect")

	traces, err = suite.Store.LoadPendingTraces(context.Background(), 500*time.Millisecond)
	assert.Nil(suite.T(), err, "No error expected on LoadPendingTraces")
	assert.Len(suite.T(), traces, 1, "Count of trace pending incorrect")

	traces, err = suite.Store.LoadPendingTraces(context.Background(), 700*time.Millisecond)
	assert.Nil(suite.T(), err, "No error expected on LoadPendingTraces")
	assert.Len(suite.T(), traces, 0, "Count of trace pending incorrect")
}
