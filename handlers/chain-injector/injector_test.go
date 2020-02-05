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
	testChainUUID     = "testChainUUID"
	testChainName     = "testChainName"
	testChainError    = "error"
	testTenantID      = "testTenantID"
)

var testChain = &types.Chain{
	UUID:                    testChainUUID,
	Name:                    testChainName,
	TenantID:                testTenantID,
	URLs:                    []string{"test"},
	ListenerDepth:           &(&struct{ x uint64 }{1}).x,
	ListenerBlockPosition:   &(&struct{ x int64 }{2}).x,
	ListenerFromBlock:       &(&struct{ x int64 }{3}).x,
	ListenerBackOffDuration: &(&struct{ x string }{"4s"}).x,
}

var testChainDefaultTenant = &types.Chain{
	UUID:                    testChainUUID,
	Name:                    testChainName,
	TenantID:                multitenancy.DefaultTenantIDName,
	URLs:                    []string{"test"},
	ListenerDepth:           &(&struct{ x uint64 }{1}).x,
	ListenerBlockPosition:   &(&struct{ x int64 }{2}).x,
	ListenerFromBlock:       &(&struct{ x int64 }{3}).x,
	ListenerBackOffDuration: &(&struct{ x string }{"4s"}).x,
}
var MockChainsByUUID = map[string]map[string]*types.Chain{
	testTenantID: {
		testChainUUID: testChain,
	},
	multitenancy.DefaultTenantIDName: {
		testChainUUID: testChainDefaultTenant,
	},
}
var MockChainsByName = map[string]map[string]*types.Chain{
	testTenantID: {
		testChainName: testChain,
	},
	multitenancy.DefaultTenantIDName: {
		testChainName: testChainDefaultTenant,
	},
}

func TestChainInjector(t *testing.T) {
	testSet := []struct {
		name          string
		multitenancy  bool
		input         func(txctx *engine.TxContext) *engine.TxContext
		expectedTxctx func(txctx *engine.TxContext) *engine.TxContext
	}{
		{
			"With multitenancy and chainUUID filled",
			true,
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain = (&chain.Chain{}).SetUUID(testChainUUID)
				txctx.WithContext(multitenancy.WithTenantID(txctx.Context(), testTenantID))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.GetChain().SetName(MockChainsByUUID[testTenantID][testChainUUID].Name)
				url := fmt.Sprintf("%s/%s/%s", testChainProxyURL, testTenantID, MockChainsByUUID[testTenantID][testChainUUID].Name)
				txctx.WithContext(proxy.With(txctx.Context(), url))
				return txctx
			},
		},
		{
			"Without multitenancy and chainUUID filled",
			false,
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain = (&chain.Chain{}).SetUUID(testChainUUID)
				txctx.WithContext(multitenancy.WithTenantID(txctx.Context(), multitenancy.DefaultTenantIDName))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.GetChain().SetName(MockChainsByUUID[multitenancy.DefaultTenantIDName][testChainUUID].Name)
				url := fmt.Sprintf("%s/%s/%s", testChainProxyURL, multitenancy.DefaultTenantIDName, MockChainsByUUID[multitenancy.DefaultTenantIDName][testChainUUID].Name)
				txctx.WithContext(proxy.With(txctx.Context(), url))
				return txctx
			},
		},
		{
			"Without multitenancy and wrong chainUUID filled",
			false,
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain = (&chain.Chain{}).SetUUID(testChainError)
				txctx.WithContext(multitenancy.WithTenantID(txctx.Context(), multitenancy.DefaultTenantIDName))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				err := errors.FromError(fmt.Errorf("error")).ExtendComponent(component)
				txctx.Envelope.Errors = append(txctx.Envelope.Errors, err)
				return txctx
			},
		},
		{
			"With multitenancy and chainName filled",
			true,
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain = (&chain.Chain{}).SetName(testChainName)
				txctx.WithContext(multitenancy.WithTenantID(txctx.Context(), testTenantID))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain.SetUUID(MockChainsByName[testTenantID][testChainName].UUID)
				url := fmt.Sprintf("%s/%s/%s", testChainProxyURL, testTenantID, testChainName)
				txctx.WithContext(proxy.With(txctx.Context(), url))
				return txctx
			},
		},
		{
			"With multitenancy and no tenantID found",
			true,
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain = (&chain.Chain{}).SetName(testChainName)
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				err := errors.InternalError("invalid tenantID not found").ExtendComponent(component)
				txctx.Envelope.Errors = append(txctx.Envelope.Errors, err)
				return txctx
			},
		},
		{
			"Without chainUUID and chainName filled",
			false,
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.WithContext(multitenancy.WithTenantID(txctx.Context(), multitenancy.DefaultTenantIDName))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				err := errors.InternalError("invalid envelope - no chain uuid or chain name are filled - cannot retrieve chain data").ExtendComponent(component)
				txctx.Envelope.Errors = append(txctx.Envelope.Errors, err)
				return txctx
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mocks.NewMockClient(mockCtrl)
	mockClient.EXPECT().GetChainByTenantAndUUID(gomock.Any(), gomock.Eq(testTenantID), gomock.Eq(testChainUUID)).Return(MockChainsByUUID[testTenantID][testChainUUID], nil).AnyTimes()
	mockClient.EXPECT().GetChainByTenantAndUUID(gomock.Any(), gomock.Eq(multitenancy.DefaultTenantIDName), gomock.Eq(testChainUUID)).Return(MockChainsByUUID[multitenancy.DefaultTenantIDName][testChainUUID], nil).AnyTimes()
	mockClient.EXPECT().GetChainByTenantAndUUID(gomock.Any(), gomock.Eq(multitenancy.DefaultTenantIDName), gomock.Eq(testChainError)).Return(nil, fmt.Errorf("error")).AnyTimes()
	mockClient.EXPECT().GetChainByTenantAndName(gomock.Any(), gomock.Eq(testTenantID), gomock.Eq(testChainName)).Return(MockChainsByName[testTenantID][testChainName], nil).AnyTimes()
	mockClient.EXPECT().GetChainByTenantAndName(gomock.Any(), gomock.Eq(multitenancy.DefaultTenantIDName), gomock.Eq(testChainName)).Return(MockChainsByName[multitenancy.DefaultTenantIDName][testChainName], nil).AnyTimes()
	mockClient.EXPECT().GetChainByTenantAndName(gomock.Any(), gomock.Eq(testChainError), gomock.Eq(testChainError)).Return(nil, fmt.Errorf("error")).AnyTimes()

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())

			h := ChainInjector(mockClient, testChainProxyURL)
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
			"Set chain UUID in Envelope",
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain = &chain.Chain{}
				txctx.WithContext(proxy.With(txctx.Context(), noErrorURL))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.Envelope.Chain.SetChainID(networkID)
				return txctx
			},
		},
		{
			"Set chain UUID in Envelope with error",
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
