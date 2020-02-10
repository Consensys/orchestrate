package enricher

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/mock/gomock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/common"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	ethclientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	registrymock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry/client/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

var testCodeNoError = []byte{1, 2}
var testCodeError = []byte{1, 2, 3}

const (
	noErrorURL                 = "noError"
	codeAtErrorURL             = "codeAtError"
	setAccountCodeHashErrorURL = "setAccountCodeHashError"
	testAccount                = "0xdbb881a51cd4023e4400cef3ef73046743f08da3"
)

func TestEnricher(t *testing.T) {
	testSet := []struct {
		name          string
		input         func(txctx *engine.TxContext) *engine.TxContext
		expectedTxctx func(txctx *engine.TxContext) *engine.TxContext
	}{
		{
			"Enrich without error",
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.WithContext(proxy.With(txctx.Context(), noErrorURL))
				txctx.Envelope.Receipt = &ethereum.Receipt{ContractAddress: testAccount}
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				return txctx
			},
		},
		{
			"Enrich with error at CodeAt",
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.WithContext(proxy.With(txctx.Context(), codeAtErrorURL))
				txctx.Envelope.Receipt = &ethereum.Receipt{ContractAddress: testAccount}
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				err := errors.InternalError(
					"could not read account code for chain %s and account %s",
					codeAtErrorURL,
					txctx.Envelope.GetReceipt().GetContractAddr().Hex(),
				).SetComponent(component)
				txctx.Envelope.Errors = append(txctx.Envelope.Errors, err)
				return txctx
			},
		},
		{
			"Enrich with error at SetAccountCodeHash",
			func(txctx *engine.TxContext) *engine.TxContext {
				txctx.WithContext(proxy.With(txctx.Context(), setAccountCodeHashErrorURL))
				txctx.Envelope.Receipt = &ethereum.Receipt{ContractAddress: testAccount}
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				err := errors.InternalError("invalid input message format").SetComponent(component)
				txctx.Envelope.Errors = append(txctx.Envelope.Errors, err)
				return txctx
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := ethclientmock.NewMockChainStateReader(mockCtrl)
	mockRegistry := registrymock.NewMockContractRegistryClient(mockCtrl)

	mockClient.EXPECT().
		CodeAt(gomock.Any(), gomock.Eq(noErrorURL), gomock.Any(), gomock.Any()).
		Return(testCodeNoError, nil).
		AnyTimes()
	mockClient.EXPECT().
		CodeAt(gomock.Any(), gomock.Eq(codeAtErrorURL), gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("error")).
		AnyTimes()
	mockClient.EXPECT().
		CodeAt(gomock.Any(), gomock.Eq(setAccountCodeHashErrorURL), gomock.Any(), gomock.Any()).
		Return(testCodeError, nil).
		AnyTimes()
	mockRegistry.EXPECT().
		SetAccountCodeHash(gomock.Any(), &svc.SetAccountCodeHashRequest{
			AccountInstance: &common.AccountInstance{},
			CodeHash:        crypto.Keccak256Hash(testCodeNoError).String(),
		}).
		Return(nil, nil).
		AnyTimes()
	mockRegistry.EXPECT().
		SetAccountCodeHash(gomock.Any(), &svc.SetAccountCodeHashRequest{
			AccountInstance: &common.AccountInstance{},
			CodeHash:        crypto.Keccak256Hash(testCodeError).String(),
		}).
		Return(nil, fmt.Errorf("error")).
		AnyTimes()

	h := Enricher(mockRegistry, mockClient)

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())

			h(test.input(txctx))

			expectedTxctx := engine.NewTxContext()
			expectedTxctx.Logger = txctx.Logger
			expectedTxctx = test.expectedTxctx(test.input(expectedTxctx))

			assert.True(t, reflect.DeepEqual(txctx, expectedTxctx), "Expected same input")
		})
	}
}
