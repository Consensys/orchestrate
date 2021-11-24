// +build unit

package rpc

import (
	"context"
	"net/http"
	"testing"

	"github.com/consensys/orchestrate/pkg/toolkit/ethclient/testutils"
	proto "github.com/consensys/orchestrate/pkg/types/ethereum"
	pkgUtils "github.com/consensys/orchestrate/pkg/utils"
	"github.com/cenkalti/backoff/v4"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func newEEAClient() *Client {
	newBackOff := func() backoff.BackOff { return pkgUtils.NewBackOff(testutils.TestConfig) }
	return NewClient(newBackOff, &http.Client{
		Transport: testutils.MockRoundTripper{},
	})
}

func TestPrivateTransactionReceipt(t *testing.T) {
	ec := newEEAClient()

	ethReceipt := &ethtypes.Receipt{
		Status:            0,
		CumulativeGasUsed: 1000,
		Logs: []*ethtypes.Log{
			{
				Address: ethcommon.BytesToAddress([]byte{0x11}),
				Topics:  []ethcommon.Hash{ethcommon.HexToHash("dead"), ethcommon.HexToHash("beef")},
				Data:    []byte{0x01, 0x00, 0xff},
			},
			{
				Address: ethcommon.BytesToAddress([]byte{0x01, 0x11}),
				Topics:  []ethcommon.Hash{ethcommon.HexToHash("dead"), ethcommon.HexToHash("beef")},
				Data:    []byte{0x01, 0x00, 0xff},
			},
		},
		ContractAddress: ethcommon.BytesToAddress([]byte{0x01, 0x11, 0x11}),
		GasUsed:         111111,
	}
	ethReceipt.Bloom = ethtypes.CreateBloom(ethtypes.Receipts{ethReceipt})

	privReceipt := &privateReceipt{
		Status:          "0x1",
		Output:          "0x12123",
		ContractAddress: "0xContractAddr",
		Logs: []*proto.Log{
			{
				Address: ethcommon.BytesToAddress([]byte{0x11}).String(),
				Topics:  []string{ethcommon.HexToHash("0x12123").String(), ethcommon.HexToHash("0x12123").String()},
				Data:    string([]byte{0x01, 0x00, 0xff}),
			},
			{
				Address: ethcommon.BytesToAddress([]byte{0x01, 0x11}).String(),
				Topics:  []string{ethcommon.HexToHash("0x12123").String(), ethcommon.HexToHash("0x12123").String()},
				Data:    string([]byte{0x01, 0x00, 0xff}),
			},
		},
		PrivateFor:  []string{"PrivateFor"},
		PrivateFrom: "PrivateFrom",
	}

	t.Run("should fetch receipts successfully including private receipt", func(t *testing.T) {
		ctx := testutils.NewContext(nil, 200, testutils.MakeRespBody(privReceipt, ""))
		// First tx receipt to fetch is the Public receipt
		ctx = context.WithValue(ctx, testutils.TestCtxKey("pre_call"),
			testutils.NewContext(nil, 200, testutils.MakeRespBody(testutils.NewReceiptResp(ethReceipt), "")))

		receipt, err := ec.PrivateTransactionReceipt(ctx, "test-endpoint", ethcommon.HexToHash(""))
		assert.NoError(t, err, "TransactionReceipt should not  error")
		assert.Equal(t, ethReceipt.CumulativeGasUsed, receipt.CumulativeGasUsed, "TransactionReceipt receipt should have correct cumulative gas used")
		assert.Equal(t, privReceipt.Output, receipt.Output, "TransactionReceipt receipt should have correct priv tx output")
		assert.Equal(t, privReceipt.ContractAddress, receipt.ContractAddress, "TransactionReceipt receipt should have correct priv contract addr")
		assert.Equal(t, uint64(0x1), receipt.Status, "TransactionReceipt receipt should have correct priv tx status")
	})

	// @TODO Following test only work running individually as there is context data leak on the way this test were implemented
	// t.Run("should fetch receipts successfully ingoring private receipt", func(t *testing.T) {
	// 	ctx := testutils.NewContext(nil, 200, testutils.MakeRespBody(nil, ""))
	// 	// First tx receipt to fetch is the Public receipt
	// 	ctx = context.WithValue(ctx, testutils.TestCtxKey("pre_call"),
	// 		testutils.NewContext(nil, 200, testutils.MakeRespBody(testutils.NewReceiptResp(ethReceipt), "")))
	// 
	// 	receipt, err := ec.PrivateTransactionReceipt(ctx, "test-endpoint", ethcommon.HexToHash(""))
	// 	assert.Error(t, err, "TransactionReceipt should error")
	// 	assert.True(t, errors.IsInvalidParameterError(err), "error should come as InvalidParameter for not found private receipts")
	// 	assert.NotNil(t, receipt, "Public receipt should be there")
	// })
	// 
	// t.Run("should fail to fetch receipt", func(t *testing.T) {
	// 	ctx := testutils.NewContext(fmt.Errorf("failed to fetch"), 500, nil)
	// 	// First tx receipt to fetch is the Public receipt
	// 	ctx = context.WithValue(ctx, testutils.TestCtxKey("pre_call"),
	// 		testutils.NewContext(nil, 200, testutils.MakeRespBody(testutils.NewReceiptResp(ethReceipt), "")))
	// 
	// 	receipt, err := ec.PrivateTransactionReceipt(ctx, "test-endpoint", ethcommon.HexToHash(""))
	// 	assert.Error(t, err, "TransactionReceipt should error")
	// 	assert.False(t, errors.IsInvalidParameterError(err), "error should not come as InvalidParameter for not found private receipts")
	// 	assert.NotNil(t, receipt, "Public receipt should be there")
	// })
}
