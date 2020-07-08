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
	types2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	clientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/client/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
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
	storeClient := clientmock.NewMockEnvelopeStoreClient(mockCtrl)
	schedulerClient := mock.NewMockTransactionSchedulerClient(mockCtrl)
	storeClient.EXPECT().Store(gomock.Eq(ctxStoreError), gomock.Any()).Return(nil, storeError).AnyTimes()
	storeClient.EXPECT().Store(gomock.Not(gomock.Eq(ctxStoreError)), gomock.Any()).Return(nil, nil).AnyTimes()
	storeClient.EXPECT().SetStatus(gomock.Eq(ctxSetStatusError), gomock.Any()).Return(nil, storeError).AnyTimes()
	storeClient.EXPECT().SetStatus(gomock.Not(gomock.Eq(ctxSetStatusError)), gomock.Any()).Return(nil, nil).AnyTimes()

	for _, test := range testSuite {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())

			h := UnsignedTxStore(storeClient, schedulerClient)
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

func TestUnsignedTxStore_TxScheduler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storeClient := clientmock.NewMockEnvelopeStoreClient(mockCtrl)
	schedulerClient := mock.NewMockTransactionSchedulerClient(mockCtrl)

	t.Run("should update the status successfully to SENT", func(t *testing.T) {
		txctx := engine.NewTxContext()
		_ = txctx.Envelope.SetID("test").SetContextLabelsValue("jobUUID", "test")
		txctx.Logger = log.NewEntry(log.New())

		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), txctx.Envelope.GetID(), &types.UpdateJobRequest{
				Transaction: &types2.ETHTransaction{
					Hash:           txctx.Envelope.GetTxHashString(),
					From:           txctx.Envelope.GetFromString(),
					To:             txctx.Envelope.GetToString(),
					Nonce:          txctx.Envelope.GetNonceString(),
					Value:          txctx.Envelope.GetValueString(),
					GasPrice:       txctx.Envelope.GetGasPriceString(),
					Gas:            txctx.Envelope.GetGasString(),
					Raw:            txctx.Envelope.GetRaw(),
					PrivateFrom:    txctx.Envelope.GetPrivateFrom(),
					PrivateFor:     txctx.Envelope.GetPrivateFor(),
					PrivacyGroupID: txctx.Envelope.GetPrivacyGroupID(),
				},
				Status: types2.StatusSent,
			}).Return(&types.JobResponse{}, nil)

		h := UnsignedTxStore(storeClient, schedulerClient)
		h(txctx)

		assert.Empty(t, txctx.Envelope.Error())
	})

	t.Run("should abort if update fails on SENT", func(t *testing.T) {
		txctx := engine.NewTxContext()
		_ = txctx.Envelope.SetID("test").SetContextLabelsValue("jobUUID", "test")
		txctx.Logger = log.NewEntry(log.New())
	
		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), txctx.Envelope.GetID(), gomock.AssignableToTypeOf(&types.UpdateJobRequest{})).
			Return(nil, fmt.Errorf("error"))
	
		h := UnsignedTxStore(storeClient, schedulerClient)
		h(txctx)
	
		assert.Len(t, txctx.Envelope.GetErrors(), 1)
	})
}
