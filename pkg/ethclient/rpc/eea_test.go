// +build unit

package rpc

import (
	"context"
	"net/http"
	"testing"

	"github.com/cenkalti/backoff/v4"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient/testutils"
	proto "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/ethereum"
	pkgUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
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
}
