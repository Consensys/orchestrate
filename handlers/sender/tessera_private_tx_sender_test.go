// +build unit

package sender

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/proxy"
)

const (
	testPrivPrivateFrom      = "test"
	testPrivPrivateFromError = "error"
	testPrivChainProxyURL    = "chainURL"
	testPrivEnclaveKey       = "0xABC0123"
)

var (
	testPrivError = errors.HTTPConnectionError("error")
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
					SetJobType(tx.JobType_ETH_TESSERA_PRIVATE_TX).
					SetData([]byte{11, 22}).
					SetPrivateFrom(testPrivPrivateFrom)

				return txctx.WithContext(proxy.With(txctx.Context(), testPrivChainProxyURL))
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetEnclaveKey(testPrivEnclaveKey)
				return txctx
			},
		},
		{
			"error without data filled",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.
					SetJobType(tx.JobType_ETH_TESSERA_PRIVATE_TX).
					SetPrivateFrom(testPrivPrivateFrom)
				return txctx.WithContext(proxy.With(txctx.Context(), testPrivChainProxyURL))
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				err := errors.DataError("cannot send transaction without data to Tessera").SetComponent(component)
				_ = txctx.Error(err)
				return txctx
			},
		},
		{
			"error in storeraw client",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.
					SetJobType(tx.JobType_ETH_TESSERA_PRIVATE_TX).
					SetData([]byte{11, 22}).
					SetPrivateFrom(testPrivPrivateFromError)
				return txctx.WithContext(proxy.With(txctx.Context(), testPrivChainProxyURL))
			},
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Error(testPrivError)
				return txctx
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockQuorumClient(mockCtrl)
	mockClient.EXPECT().StoreRaw(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(testPrivPrivateFrom)).
		Return(testPrivEnclaveKey, nil).AnyTimes()
	mockClient.EXPECT().StoreRaw(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(testPrivPrivateFromError)).
		Return("", testPrivError).AnyTimes()

	for _, test := range testSet {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())

			h := TesseraPrivateTxSender(mockClient)
			h(test.input(txctx))

			expectedTxctx := engine.NewTxContext()
			expectedTxctx.Logger = txctx.Logger
			expectedTxctx = test.expectedTxctx(test.input(expectedTxctx))

			assert.True(t, reflect.DeepEqual(txctx, expectedTxctx), fmt.Sprintf("%s: expected same input", test.name))
		})
	}
}
