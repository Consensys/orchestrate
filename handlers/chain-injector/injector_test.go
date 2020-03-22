// +build unit

package chaininjector

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	ethclientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
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
	ListenerCurrentBlock:    &(&struct{ x uint64 }{2}).x,
	ListenerStartingBlock:   &(&struct{ x uint64 }{3}).x,
	ListenerBackOffDuration: &(&struct{ x string }{"4s"}).x,
}

var testChainDefaultTenant = &types.Chain{
	UUID:                    testChainUUID,
	Name:                    testChainName,
	TenantID:                multitenancy.DefaultTenantIDName,
	URLs:                    []string{"test"},
	ListenerDepth:           &(&struct{ x uint64 }{1}).x,
	ListenerCurrentBlock:    &(&struct{ x uint64 }{2}).x,
	ListenerStartingBlock:   &(&struct{ x uint64 }{3}).x,
	ListenerBackOffDuration: &(&struct{ x string }{"4s"}).x,
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
		input         func(txctx *engine.TxContext) *engine.TxContext
		expectedTxctx func(txctx *engine.TxContext) *engine.TxContext
	}{
		{
			"chainUUID filled",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetChainName(testChainName)
				txctx.WithContext(multitenancy.WithTenantID(txctx.Context(), testTenantID))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetChainUUID(MockChainsByName[testTenantID][testChainName].UUID)
				url := fmt.Sprintf("%s/%s", testChainProxyURL, testChainUUID)
				txctx.WithContext(proxy.With(txctx.Context(), url))
				return txctx
			},
		},
		{
			"Without chainUUID and chainName filled",
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.WithContext(multitenancy.WithTenantID(txctx.Context(), multitenancy.DefaultTenantIDName))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.AppendError(errors.DataError("no chain name found").ExtendComponent(component))
				return txctx
			},
		},
		{
			"error when calling the chain registry",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetChainName(testChainError)
				txctx.WithContext(multitenancy.WithTenantID(txctx.Context(), testChainError))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.AppendError(errors.FromError(fmt.Errorf("error")).ExtendComponent(component))
				return txctx
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockChainRegistryClient(mockCtrl)
	mockClient.EXPECT().GetChainByName(gomock.Any(), gomock.Eq(testChainName)).Return(MockChainsByName[testTenantID][testChainName], nil).AnyTimes()
	mockClient.EXPECT().GetChainByName(gomock.Any(), gomock.Eq(testChainError)).Return(nil, fmt.Errorf("error")).AnyTimes()

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())

			h := ChainUUIDHandler(mockClient, testChainProxyURL)
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
				txctx.WithContext(proxy.With(txctx.Context(), noErrorURL))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetChainID(networkID)
				return txctx
			},
		},
		{
			"Set chain UUID in Envelope with error",
			func(txctx *engine.TxContext) *engine.TxContext {
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

	h := ChainIDInjectorHandler(ec)

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
