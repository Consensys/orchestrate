package rpc

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/rpc/geth"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
)

const chainID = 888

var chainIDBigInt = big.NewInt(int64(chainID))

var errTest = fmt.Errorf("test error")
var privateArgs = &types.PrivateArgs{
	PrivateFrom:   "0x01",
	PrivateFor:    []string{"0x02"},
	PrivateTxType: "abc",
}

var expectedPrivateForArg = map[string]interface{}{
	"privateFor": privateArgs.PrivateFor,
}

var ctx, _ = context.WithCancel(context.Background())
var gethClient *Client
var ctrl *gomock.Controller
var mockRPCClient *mocks.MockClient

func setupTest(t *testing.T) {
	ctrl = gomock.NewController(t)
	gethClient = NewClient(&geth.Config{})
	mockRPCClient = mocks.NewMockClient(ctrl)
}

func TestQuorumRawPrivateTransaction(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	setMockClient(mockRPCClient)
	mockRPCClient.
		EXPECT().
		CallContext(ctx, gomock.Any(), "eth_sendRawPrivateTransaction", "0x010203", expectedPrivateForArg).
		Return(nil).
		SetArg(1, "0x1234").
		Times(1)

	hash, err := gethClient.SendQuorumRawPrivateTransaction(ctx, chainIDBigInt, []byte{1, 2, 3}, privateArgs.PrivateFor)
	assert.NoError(t, err)
	assert.Equal(t, ethcommon.HexToHash("0x1234"), hash)
}

func TestQuorumRawPrivateTransactionWhenRPCCallFails(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	setMockClient(mockRPCClient)
	mockRPCClient.
		EXPECT().
		CallContext(ctx, gomock.Any(), "eth_sendRawPrivateTransaction", "0x010203", expectedPrivateForArg).
		Return(errTest).
		Times(1)

	hash, err := gethClient.SendQuorumRawPrivateTransaction(ctx, chainIDBigInt, []byte{1, 2, 3}, privateArgs.PrivateFor)
	assert.Error(t, err, errTest.Error())
	assert.Equal(t, ethcommon.HexToHash("0x0"), hash)
}

func TestSendRawPrivateTransaction(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	setMockClient(mockRPCClient)
	mockRPCClient.
		EXPECT().
		CallContext(ctx, gomock.Any(), "eea_sendRawTransaction", "0x010203").
		Return(nil).
		SetArg(1, "0x1234").
		Times(1)

	hash, err := gethClient.SendRawPrivateTransaction(ctx, chainIDBigInt, []byte{1, 2, 3}, privateArgs)
	assert.NoError(t, err)
	assert.Equal(t, ethcommon.HexToHash("0x1234"), hash)
}

func TestSendRawPrivateTransactionWhenRPCCallFails(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	setMockClient(mockRPCClient)
	mockRPCClient.
		EXPECT().
		CallContext(ctx, gomock.Any(), "eea_sendRawTransaction", "0x010203").
		Return(errTest).
		Times(1)

	hash, err := gethClient.SendRawPrivateTransaction(ctx, chainIDBigInt, []byte{1, 2, 3}, privateArgs)
	assert.Error(t, err, errTest.Error())
	assert.Equal(t, ethcommon.HexToHash("0x0"), hash)
}

func TestReturnErrorIfCannotGetRPC(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	hash, err := gethClient.SendRawPrivateTransaction(ctx, chainIDBigInt, []byte{1, 2, 3}, privateArgs)
	e := errors.FromError(err)
	assert.Equal(t, "no RPC connection registered for chain \"888\"", e.GetMessage(), "Error message should be correct")
	assert.Equal(t, "08300", e.Hex(), "Error hex code should be correct")
	assert.Equal(t, ethcommon.HexToHash("0x0"), hash)
}

func setMockClient(mockClient *mocks.MockClient) {
	gethClient.rpcs[strconv.Itoa(chainID)] = mockClient
}
