package rpc

import (
	"context"
	"errors"
	"math/big"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/rpc/geth"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
)

const chainID = 888

var errTest = errors.New("test error")
var ctx, _ = context.WithCancel(context.Background())
var gethClient *Client
var ctrl *gomock.Controller
var mockRPCClient *mocks.MockClient

func setupTest(t *testing.T) {
	ctrl = gomock.NewController(t)
	gethClient = NewClient(&geth.Config{})
	mockRPCClient = mocks.NewMockClient(ctrl)
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
		AnyTimes()

	clientType, err := gethClient.GetClientType(ctx, big.NewInt(int64(chainID)))
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
		AnyTimes()

	clientType, err := gethClient.GetClientType(ctx, big.NewInt(int64(chainID)))
	assert.EqualError(t, err, errTest.Error())
	assert.Equal(t, types.UnknownClient, clientType)
}

func TestReturnErrorIfCannotGetRPC(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	clientType, err := gethClient.GetClientType(ctx, big.NewInt(int64(chainID)))

	assert.EqualError(t, err, "no RPC connection registered for chain \"888\"")
	assert.Equal(t, types.UnknownClient, clientType)
}

func setMockClient(mockClient *mocks.MockClient) {
	gethClient.rpcs[strconv.Itoa(chainID)] = mockClient
}
