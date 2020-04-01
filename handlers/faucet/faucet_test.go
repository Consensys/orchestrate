// +build unit

package faucet

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	registryMock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
	faucetMock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet/mocks"
	faucettypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"

	"github.com/golang/mock/gomock"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

const (
	testChainUUID            = "testChain"
	testChainUUIDError       = "chainError"
	testURL                  = "testURL"
	testIDNoError            = "noError"
	testIDSelfCreditError    = "selfCreditError"
	testIDFaucetWarningError = "warningError"
	testIDFaucetError        = "error"
)

var testFaucets = []*types.Faucet{
	{
		UUID:            "testUUID",
		Name:            "testName",
		TenantID:        "testTenantID",
		ChainRule:       testChainUUID,
		CreditorAccount: "0x7e654d251da770a068413677967f6d3ea2fea9e4",
		MaxBalance:      "10",
		Amount:          "1",
		Cooldown:        "1s",
	},
}

var errGetFaucetsByChainRule = fmt.Errorf("error")
var errSelfCredit = errors.FaucetSelfCreditWarning("error")
var errFaucetWarning = errors.FaucetWarning("error")
var errFaucet = fmt.Errorf("error")

func TestFaucet(t *testing.T) {
	testSet := []struct {
		name          string
		input         func(txctx *engine.TxContext) *engine.TxContext
		expectedTxctx func(txctx *engine.TxContext) *engine.TxContext
	}{
		{
			"credit without error",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetChainUUID(testChainUUID).SetID(testIDNoError)
				txctx.WithContext(proxy.With(txctx.Context(), testURL))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				return txctx
			},
		},
		{
			"credit without chainUUID",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetID(testIDNoError)
				txctx.WithContext(proxy.With(txctx.Context(), testURL))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				return txctx
			},
		},
		{
			"txctx with parentTxID context label",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetChainUUID(testChainUUID).SetContextLabelsValue("faucet.parentTxID", "test")
				txctx.WithContext(proxy.With(txctx.Context(), testURL))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				return txctx
			},
		},
		{
			"txctx with error when get faucets",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetChainUUID(testChainUUIDError)
				txctx.WithContext(proxy.With(txctx.Context(), testURL))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.AppendError(errors.FaucetWarning("could not get faucets for chain rule '%s' - got %v", txctx.Envelope.GetChainUUID(), errGetFaucetsByChainRule).ExtendComponent(component))
				return txctx
			},
		},
		{
			"credit with self credit error",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetChainUUID(testChainUUID).SetID(testIDSelfCreditError)
				txctx.WithContext(proxy.With(txctx.Context(), testURL))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				return txctx
			},
		},
		{
			"credit with faucet warning error",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetChainUUID(testChainUUID).SetID(testIDFaucetWarningError)
				txctx.WithContext(proxy.With(txctx.Context(), testURL))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				return txctx
			},
		},
		{
			"credit with error",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetChainUUID(testChainUUID).SetID(testIDFaucetError)
				txctx.WithContext(proxy.With(txctx.Context(), testURL))
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				return txctx
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRegistry := registryMock.NewMockFaucetClient(mockCtrl)
	mockRegistry.EXPECT().GetFaucetsByChainRule(gomock.Any(), gomock.Eq(testChainUUID)).Return(testFaucets, nil).AnyTimes()
	mockRegistry.EXPECT().GetFaucetsByChainRule(gomock.Any(), gomock.Eq(testChainUUIDError)).Return(nil, errGetFaucetsByChainRule).AnyTimes()

	mockFaucet := faucetMock.NewMockFaucet(mockCtrl)
	mockFaucet.EXPECT().Credit(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, req *faucettypes.Request) (*big.Int, error) {
		switch req.ParentTxID {
		case testIDSelfCreditError:
			return nil, errSelfCredit
		case testIDFaucetWarningError:
			return nil, errFaucetWarning
		case testIDFaucetError:
			return nil, errFaucet
		}
		return big.NewInt(10), nil
	}).AnyTimes()

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())

			h := Faucet(mockFaucet, mockRegistry)
			h(test.input(txctx))

			expectedTxctx := engine.NewTxContext()
			expectedTxctx.Logger = txctx.Logger
			expectedTxctx = test.expectedTxctx(test.input(expectedTxctx))

			assert.True(t, reflect.DeepEqual(txctx, expectedTxctx), "Expected same input")
		})
	}
}
