package rpc

import (
	"context"
	"errors"
	"math/big"
	"strconv"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/rpc/geth"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
)

const chainID = 888

var chainIDBigInt = big.NewInt(int64(chainID))

var errTest = errors.New("test error")
var privateArgs = &types.PrivateArgs{
	PrivateFrom:   "0x01",
	PrivateFor:    []string{"0x02"},
	PrivateTxType: "abc",
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
		CallContext(ctx, gomock.Any(), "eth_sendRawPrivateTransaction", "0x010203", []string{"0x02"}).
		Return(nil).
		SetArg(1, "0x1234").
		Times(1)

	hash, err := gethClient.SendQuorumRawPrivateTransaction(ctx, chainIDBigInt, []byte{1, 2, 3}, privateArgs)
	assert.NoError(t, err)
	assert.Equal(t, ethcommon.HexToHash("0x1234"), hash)
}

func TestQuorumRawPrivateTransactionWhenRPCCallFails(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	setMockClient(mockRPCClient)
	mockRPCClient.
		EXPECT().
		CallContext(ctx, gomock.Any(), "eth_sendRawPrivateTransaction", "0x010203", []string{"0x02"}).
		Return(errTest).
		Times(1)

	hash, err := gethClient.SendQuorumRawPrivateTransaction(ctx, chainIDBigInt, []byte{1, 2, 3}, privateArgs)
	assert.Error(t, err, errTest.Error())
	assert.Equal(t, ethcommon.HexToHash("0x0"), hash)
}

func TestReturnErrorIfCannotGetRPCWhenSendingQuorumPrivateTransaction(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	hash, err := gethClient.SendQuorumRawPrivateTransaction(ctx, chainIDBigInt, []byte{1, 2, 3}, privateArgs)

	assert.EqualError(t, err, "no RPC connection registered for chain \"888\"")
	assert.Equal(t, ethcommon.HexToHash("0x0"), hash)
}

func TestSendRawPrivateTransaction(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	setMockClient(mockRPCClient)
	mockRPCClient.
		EXPECT().
		CallContext(ctx, gomock.Any(), "eea_sendRawTransaction", "0x0102038430783031c5843078303283616263").
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
		CallContext(ctx, gomock.Any(), "eea_sendRawTransaction", "0x0102038430783031c5843078303283616263").
		Return(errTest).
		Times(1)

	hash, err := gethClient.SendRawPrivateTransaction(ctx, chainIDBigInt, []byte{1, 2, 3}, privateArgs)
	assert.Error(t, err, errTest.Error())
	assert.Equal(t, ethcommon.HexToHash("0x0"), hash)
}

func TestReturnErrorIfCannotGetRPCWhenSendingRawPrivateTransaction(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	hash, err := gethClient.SendRawPrivateTransaction(ctx, chainIDBigInt, []byte{1, 2, 3}, privateArgs)

	assert.EqualError(t, err, "no RPC connection registered for chain \"888\"")
	assert.Equal(t, ethcommon.HexToHash("0x0"), hash)
}

func TestDetectClientVersion(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	setMockClient(mockRPCClient)
	mockRPCClient.
		EXPECT().
		CallContext(ctx, gomock.Any(), "web3_clientVersion").
		Return(nil).
		SetArg(1, "pantheon/1.1.1").
		Times(1)

	clientType, err := gethClient.GetClientType(ctx, chainIDBigInt)
	assert.NoError(t, err)
	assert.Equal(t, types.PantheonClient, clientType)
}

func TestReturnErrorIfClientVersionMethodFails(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	setMockClient(mockRPCClient)
	mockRPCClient.
		EXPECT().
		CallContext(ctx, gomock.Any(), "web3_clientVersion").
		Return(errTest).
		Times(1)

	clientType, err := gethClient.GetClientType(ctx, chainIDBigInt)
	assert.EqualError(t, err, errTest.Error())
	assert.Equal(t, types.UnknownClient, clientType)
}

func TestReturnErrorIfCannotGetRPCWhenDetectingClient(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	clientType, err := gethClient.GetClientType(ctx, chainIDBigInt)

	assert.EqualError(t, err, "no RPC connection registered for chain \"888\"")
	assert.Equal(t, types.UnknownClient, clientType)
}

func setMockClient(mockClient *mocks.MockClient) {
	gethClient.rpcs[strconv.Itoa(chainID)] = mockClient
}
