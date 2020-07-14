// +build unit

package storer

import (
	"context"
	"fmt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client/mock"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	clientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/client/mock"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

func TestRawTxStore(t *testing.T) {
	testSet := []struct {
		name           string
		input          func(txctx *engine.TxContext) *engine.TxContext
		expectedStatus svc.Status
	}{
		{
			"Store",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetChainID(big.NewInt(1)).SetID("test")
				return txctx
			},
			svc.Status_PENDING,
		},
		{
			"Store envelope without Metadata UUID",
			func(txctx *engine.TxContext) *engine.TxContext {
				_ = txctx.Envelope.SetChainID(big.NewInt(1)).SetID("test")
				err := errors.InternalError("error").ExtendComponent(component)
				txctx.Envelope.Errors = append(txctx.Envelope.Errors, err)
				return txctx
			},
			svc.Status_ERROR,
		},
	}

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txctx := engine.NewTxContext()
			txctx.Logger = log.NewEntry(log.New())
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			registry := clientmock.NewMockEnvelopeStoreClient(mockCtrl)
			schedulerClient := mock.NewMockTransactionSchedulerClient(mockCtrl)
			registry.EXPECT().Store(gomock.Any(), gomock.AssignableToTypeOf(&svc.StoreRequest{}))
			registry.EXPECT().SetStatus(gomock.Any(), &svc.SetStatusRequest{
				Id:     "test",
				Status: test.expectedStatus,
			})
			registry.EXPECT().LoadByID(gomock.Any(), &svc.LoadByIDRequest{
				Id: "test",
			}).Return(&svc.StoreResponse{
				StatusInfo: &svc.StatusInfo{Status: test.expectedStatus},
			}, nil)

			h := RawTxStore(registry, schedulerClient)
			h(test.input(txctx))
			e, _ := registry.LoadByID(txctx.Context(), &svc.LoadByIDRequest{Id: txctx.Envelope.GetID()})
			assert.Equal(t, test.expectedStatus, e.StatusInfo.Status, "Expected same status")
		})
	}
}

func TestRawTxStore_TxScheduler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	registry := clientmock.NewMockEnvelopeStoreClient(mockCtrl)
	schedulerClient := mock.NewMockTransactionSchedulerClient(mockCtrl)

	t.Run("should update the status successfully to PENDING", func(t *testing.T) {
		txctx := engine.NewTxContext()
		_ = txctx.Envelope.SetID("test").SetContextLabelsValue("jobUUID", "test")
		txctx.Logger = log.NewEntry(log.New())

		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), txctx.Envelope.GetID(), &types.UpdateJobRequest{
				Transaction: &types.ETHTransaction{
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
				Status: types.StatusPending,
			}).
			Return(&types.JobResponse{}, nil)

		h := RawTxStore(registry, schedulerClient)
		h(txctx)

		assert.Empty(t, txctx.Envelope.Error())
	})

	t.Run("should abort if update fails on PENDING", func(t *testing.T) {
		txctx := engine.NewTxContext()
		_ = txctx.Envelope.SetID("test").SetContextLabelsValue("jobUUID", "test")
		txctx.Logger = log.NewEntry(log.New())

		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), txctx.Envelope.GetID(), gomock.AssignableToTypeOf(&types.UpdateJobRequest{})).
			Return(nil, fmt.Errorf("error"))

		h := RawTxStore(registry, schedulerClient)
		h(txctx)

		assert.Len(t, txctx.Envelope.GetErrors(), 1)
	})

	t.Run("should set status to RECOVERING if txctx contains errors", func(t *testing.T) {
		txctx := engine.NewTxContext()
		_ = txctx.Envelope.SetID("test").SetContextLabelsValue("jobUUID", "test")
		txctx.Logger = log.NewEntry(log.New())
		_ = txctx.AbortWithError(fmt.Errorf("error"))

		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), txctx.Envelope.GetID(), gomock.Any()).
			Return(&types.JobResponse{}, nil)
		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), txctx.Envelope.GetID(), &types.UpdateJobRequest{
				Status: types.StatusRecovering,
				Message: fmt.Sprintf(
					"transaction attempt with nonce %v and sender %v failed with error: %v",
					txctx.Envelope.GetNonceString(),
					txctx.Envelope.GetFromString(),
					txctx.Envelope.Error(),
				),
			}).
			Return(&types.JobResponse{}, nil)

		h := RawTxStore(registry, schedulerClient)
		h(txctx)
	})

	t.Run("should return if update fails on RECOVERING", func(t *testing.T) {
		txctx := engine.NewTxContext()
		_ = txctx.Envelope.SetID("test").SetContextLabelsValue("jobUUID", "test")
		txctx.Logger = log.NewEntry(log.New())
		_ = txctx.AbortWithError(fmt.Errorf("error"))

		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), txctx.Envelope.GetID(), gomock.Any()).
			Return(&types.JobResponse{}, nil)
		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), txctx.Envelope.GetID(), &types.UpdateJobRequest{
				Status: types.StatusRecovering,
				Message: fmt.Sprintf(
					"transaction attempt with nonce %v and sender %v failed with error: %v",
					txctx.Envelope.GetNonceString(),
					txctx.Envelope.GetFromString(),
					txctx.Envelope.Error(),
				),
			}).
			Return(nil, fmt.Errorf("error"))

		h := RawTxStore(registry, schedulerClient)
		h(txctx)
	})

	t.Run("should fail if computed tx hash is different than retrieved hash", func(t *testing.T) {
		txctx := engine.NewTxContext()
		_ = txctx.Envelope.SetID("test").SetContextLabelsValue("jobUUID", "test").MustSetTxHashString("0x1")
		txctx.Logger = log.NewEntry(log.New())

		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), txctx.Envelope.GetID(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, jobUUID string, request *types.UpdateJobRequest) (*types.JobResponse, error) {
				_ = txctx.Envelope.MustSetTxHashString("0x2")
				return &types.JobResponse{}, nil
			})

		h := RawTxStore(registry, schedulerClient)
		h(txctx)

		assert.Len(t, txctx.Envelope.Errors, 1)
	})
}
