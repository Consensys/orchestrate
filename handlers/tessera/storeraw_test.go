// +build unit

package tessera

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/tessera/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
)

const (
	testPrivateFrom      = "test"
	testPrivateFromError = "error"
	testChainProxyURL    = "test"
	testEnclaveKey       = "testEnclaveKey"
)

var (
	storerawError = errors.HTTPConnectionError("error")
)

func TestStoreRaw(t *testing.T) {
	testSet := []struct {
		name          string
		input         func(txctx *engine.TxContext) *engine.TxContext
		expectedTxctx func(txctx *engine.TxContext) *engine.TxContext
	}{
		{
			"storeraw without error",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.
					SetMethod(tx.Method_ETH_SENDRAWPRIVATETRANSACTION).
					SetData([]byte{11, 22}).
					SetPrivateFrom(testPrivateFrom)
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetEnclaveKey(testEnclaveKey)
				return txctx
			},
		},
		{
			"error without data filled",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.
					SetMethod(tx.Method_ETH_SENDRAWPRIVATETRANSACTION).
					SetPrivateFrom(testPrivateFrom)
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				err := errors.DataError("can not send transaction without data to Tessera").SetComponent(component)
				_ = txctx.Envelope.AppendError(err)
				return txctx
			},
		},
		{
			"error in storeraw client",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.
					SetMethod(tx.Method_ETH_SENDRAWPRIVATETRANSACTION).
					SetData([]byte{11, 22}).
					SetPrivateFrom(testPrivateFromError)
				return txctx
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.AppendError(storerawError)
				return txctx
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockClient(mockCtrl)
	mockClient.EXPECT().StoreRaw(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(testPrivateFrom)).Return(testEnclaveKey, nil).AnyTimes()
	mockClient.EXPECT().StoreRaw(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(testPrivateFromError)).Return("", storerawError).AnyTimes()

	for _, test := range testSet {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())

			h := StoreRaw(mockClient, testChainProxyURL)
			h(test.input(txctx))

			expectedTxctx := engine.NewTxContext()
			expectedTxctx.Logger = txctx.Logger
			expectedTxctx = test.expectedTxctx(test.input(expectedTxctx))

			assert.True(t, reflect.DeepEqual(txctx, expectedTxctx), "Expected same input")
		})
	}
}
