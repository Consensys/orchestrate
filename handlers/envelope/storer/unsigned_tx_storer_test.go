// +build unit

package storer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	clientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/client/mock"
	"golang.org/x/net/context"
)

type store string
type setStatus string
type noError string

func TestUnsignedTxStore(t *testing.T) {
	ctxStoreError := context.WithValue(context.Background(), store("store"), "error")
	ctxSetStatusError := context.WithValue(context.Background(), setStatus("setStatus"), "error")
	ctxNoError := context.WithValue(context.Background(), noError("noError"), "noError")

	storeError := fmt.Errorf("error")

	testSuite := []struct {
		name          string
		txctx         func(txctx *engine.TxContext) *engine.TxContext
		expectedTxctx func(txctx *engine.TxContext) *engine.TxContext
	}{
		{
			name: "Errors in envelope",
			txctx: func(txctx *engine.TxContext) *engine.TxContext {
				err := errors.FromError(fmt.Errorf("error")).ExtendComponent(component)
				txctx.Envelope.Errors = append(txctx.Envelope.Errors, err)
				return txctx
			},
			expectedTxctx: func(txctx *engine.TxContext) *engine.TxContext {
				return txctx
			},
		},
		{
			name: "UnsignedTxStore without error",
			txctx: func(txctx *engine.TxContext) *engine.TxContext {
				txctx.WithContext(ctxNoError)
				return txctx
			},
			expectedTxctx: func(txctx *engine.TxContext) *engine.TxContext {
				return txctx
			},
		},
		{
			name: "error when store envelope",
			txctx: func(txctx *engine.TxContext) *engine.TxContext {
				txctx.WithContext(ctxStoreError)
				return txctx
			},
			expectedTxctx: func(txctx *engine.TxContext) *engine.TxContext {
				err := errors.FromError(storeError).ExtendComponent(component)
				txctx.Envelope.Errors = append(txctx.Envelope.Errors, err)
				return txctx
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	client := clientmock.NewMockEnvelopeStoreClient(mockCtrl)
	client.EXPECT().Store(gomock.Eq(ctxStoreError), gomock.Any()).Return(nil, storeError).AnyTimes()
	client.EXPECT().Store(gomock.Not(gomock.Eq(ctxStoreError)), gomock.Any()).Return(nil, nil).AnyTimes()
	client.EXPECT().SetStatus(gomock.Eq(ctxSetStatusError), gomock.Any()).Return(nil, storeError).AnyTimes()
	client.EXPECT().SetStatus(gomock.Not(gomock.Eq(ctxSetStatusError)), gomock.Any()).Return(nil, nil).AnyTimes()

	for _, test := range testSuite {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())

			h := UnsignedTxStore(client)
			h(test.txctx(txctx))

			expectedTxctx := engine.NewTxContext()
			expectedTxctx.Logger = txctx.Logger
			expectedTxctx = test.expectedTxctx(test.txctx(expectedTxctx))

			t.Log(txctx.Envelope.InternalLabels)
			t.Log(expectedTxctx.Envelope.InternalLabels)

			assert.True(t, reflect.DeepEqual(txctx, expectedTxctx), "Expected same input")
		})
	}
}
