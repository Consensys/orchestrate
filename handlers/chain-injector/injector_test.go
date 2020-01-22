package chaininjector

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	ethclientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/chain"
)

const (
	testChainProxyURL = "test"
	testNodeID        = "testNodeID"
	testNodeName      = "testNodeName"
	testNodeError     = "error"
	testTenantID      = "testTenantID"
)

var testNode = &types.Node{
	ID:                      testNodeID,
	Name:                    testNodeName,
	TenantID:                testTenantID,
	URLs:                    []string{"test"},
	ListenerDepth:           1,
	ListenerBlockPosition:   2,
	ListenerFromBlock:       3,
	ListenerBackOffDuration: "4s",
}

var testNodeDefaultTenant = &types.Node{
	ID:                      testNodeID,
	Name:                    testNodeName,
	TenantID:                multitenancy.DefaultTenantIDName,
	URLs:                    []string{"test"},
	ListenerDepth:           1,
	ListenerBlockPosition:   2,
	ListenerFromBlock:       3,
	ListenerBackOffDuration: "4s",
}
var MockNodesByID = map[string]map[string]*types.Node{
	testTenantID: {
		testNodeID: testNode,
	},
	multitenancy.DefaultTenantIDName: {
		testNodeID: testNodeDefaultTenant,
	},
}
var MockNodesByName = map[string]map[string]*types.Node{
	testTenantID: {
		testNodeName: testNode,
	},
	multitenancy.DefaultTenantIDName: {
		testNodeName: testNodeDefaultTenant,
	},
}

func TestNodeInjector(t *testing.T) {
	testSet := []struct {
		name          string
		multitenancy  bool
		input         func(txctx *engine.TxContext) *engine.TxContext
		expectedTxctx func(txctx *engine.TxContext) *engine.TxContext
	}{
		{
			"With multitenancy and nodeID filled",
			true,
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain = (&chain.Chain{}).SetNodeID(testNodeID)
				txctx.WithContext(multitenancy.WithTenantID(txctx.Context(), testTenantID))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain.SetNodeName(MockNodesByID[testTenantID][testNodeID].Name)
				url := fmt.Sprintf("%s/%s/%s", testChainProxyURL, testTenantID, MockNodesByID[testTenantID][testNodeID].Name)
				txctx.WithContext(proxy.With(txctx.Context(), url))
				return txctx
			},
		},
		{
			"Without multitenancy and nodeID filled",
			false,
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain = (&chain.Chain{}).SetNodeID(testNodeID)
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain.SetNodeName(MockNodesByID[multitenancy.DefaultTenantIDName][testNodeID].Name)
				url := fmt.Sprintf("%s/%s/%s", testChainProxyURL, multitenancy.DefaultTenantIDName, MockNodesByID[multitenancy.DefaultTenantIDName][testNodeID].Name)
				txctx.WithContext(proxy.With(txctx.Context(), url))
				return txctx
			},
		},
		{
			"Without multitenancy and wrong nodeID filled",
			false,
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain = (&chain.Chain{}).SetNodeID(testNodeError)
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				err := errors.FromError(fmt.Errorf("error")).ExtendComponent(component)
				txctx.Envelope.Errors = append(txctx.Envelope.Errors, err)
				return txctx
			},
		},
		{
			"With multitenancy and nodeName filled",
			true,
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain = (&chain.Chain{}).SetNodeName(testNodeName)
				txctx.WithContext(multitenancy.WithTenantID(txctx.Context(), testTenantID))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain.SetNodeID(MockNodesByName[testTenantID][testNodeName].ID)
				url := fmt.Sprintf("%s/%s/%s", testChainProxyURL, testTenantID, testNodeName)
				txctx.WithContext(proxy.With(txctx.Context(), url))
				return txctx
			},
		},
		{
			"With multitenancy and no tenantID found",
			true,
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain = (&chain.Chain{}).SetNodeName(testNodeName)
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				err := errors.InternalError("invalid tenantID not found").ExtendComponent(component)
				txctx.Envelope.Errors = append(txctx.Envelope.Errors, err)
				return txctx
			},
		},
		{
			"Without nodeID and nodeName filled",
			false,
			func(txctx *engine.TxContext) *engine.TxContext {
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				err := errors.InternalError("invalid envelope - no node id or node name are filled - cannot retrieve chain id").ExtendComponent(component)
				txctx.Envelope.Errors = append(txctx.Envelope.Errors, err)
				return txctx
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mocks.NewMockClient(mockCtrl)
	mockClient.EXPECT().GetNodeByTenantAndNodeID(gomock.Any(), gomock.Eq(testTenantID), gomock.Eq(testNodeID)).Return(MockNodesByID[testTenantID][testNodeID], nil).AnyTimes()
	mockClient.EXPECT().GetNodeByTenantAndNodeID(gomock.Any(), gomock.Eq(multitenancy.DefaultTenantIDName), gomock.Eq(testNodeID)).Return(MockNodesByID[multitenancy.DefaultTenantIDName][testNodeID], nil).AnyTimes()
	mockClient.EXPECT().GetNodeByTenantAndNodeID(gomock.Any(), gomock.Eq(multitenancy.DefaultTenantIDName), gomock.Eq(testNodeError)).Return(nil, fmt.Errorf("error")).AnyTimes()
	mockClient.EXPECT().GetNodeByTenantAndNodeName(gomock.Any(), gomock.Eq(testTenantID), gomock.Eq(testNodeName)).Return(MockNodesByName[testTenantID][testNodeName], nil).AnyTimes()
	mockClient.EXPECT().GetNodeByTenantAndNodeName(gomock.Any(), gomock.Eq(multitenancy.DefaultTenantIDName), gomock.Eq(testNodeName)).Return(MockNodesByName[multitenancy.DefaultTenantIDName][testNodeName], nil).AnyTimes()
	mockClient.EXPECT().GetNodeByTenantAndNodeName(gomock.Any(), gomock.Eq(testNodeError), gomock.Eq(testNodeError)).Return(nil, fmt.Errorf("error")).AnyTimes()

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())

			h := NodeInjector(test.multitenancy, mockClient, testChainProxyURL)
			h(test.input(txctx))

			expectedTxctx := engine.NewTxContext()
			expectedTxctx.Logger = txctx.Logger
			expectedTxctx = test.expectedTxctx(test.input(expectedTxctx))

			assert.True(t, reflect.DeepEqual(txctx, expectedTxctx), "Expected same input")
		})
	}
}

const (
	noErrorURL = "test"
	errorURL   = "error"
)

var (
	networkID    = big.NewInt(123)
	errorNetwork = fmt.Errorf("error")
)

func TestChainIDInjector(t *testing.T) {
	testSet := []struct {
		name          string
		txctx         func(txctx *engine.TxContext) *engine.TxContext
		expectedTxctx func(txctx *engine.TxContext) *engine.TxContext
	}{
		{
			"Set chain ID in Envelope",
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain = &chain.Chain{}
				txctx.WithContext(proxy.With(txctx.Context(), noErrorURL))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain.SetID(networkID)
				return txctx
			},
		},
		{
			"Set chain ID in Envelope with error",
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain = &chain.Chain{}
				txctx.WithContext(proxy.With(txctx.Context(), errorURL))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				err := errors.FromError(errorNetwork).ExtendComponent(component)
				txctx.Envelope.Errors = append(txctx.Envelope.Errors, err)
				return txctx
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ec := ethclientmock.NewMockChainSyncReader(mockCtrl)
	ec.EXPECT().Network(gomock.Any(), gomock.Eq(noErrorURL)).Return(networkID, nil).AnyTimes()
	ec.EXPECT().Network(gomock.Any(), gomock.Eq(errorURL)).Return(nil, errorNetwork).AnyTimes()

	h := ChainIDInjector(ec)

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())

			h(test.txctx(txctx))

			expectedTxctx := engine.NewTxContext()
			expectedTxctx.Logger = txctx.Logger
			expectedTxctx = test.expectedTxctx(test.txctx(expectedTxctx))

			assert.True(t, reflect.DeepEqual(txctx, expectedTxctx), "Expected same input")
		})
	}

}
