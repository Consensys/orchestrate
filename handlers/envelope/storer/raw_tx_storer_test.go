// +build unit

package storer

import (
	"context"
	"fmt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"

	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/txscheduler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client/mock"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

func TestRawTxStore(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	schedulerClient := mock.NewMockTransactionSchedulerClient(mockCtrl)

	t.Run("should update the status successfully to PENDING", func(t *testing.T) {
		txctx := engine.NewTxContext()
		_ = txctx.Envelope.SetJobUUID("jobUUID")
		txctx.Logger = log.NewEntry(log.New())

		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), txctx.Envelope.GetJobUUID(), &txschedulertypes.UpdateJobRequest{
				Transaction: &entities.ETHTransaction{
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
				Status: utils.StatusPending,
			}).
			Return(&txschedulertypes.JobResponse{}, nil)

		RawTxStore(schedulerClient)(txctx)

		assert.Empty(t, txctx.Envelope.Error())
	})

	t.Run("should override txHash if hash retrieved is different", func(t *testing.T) {
		txctx := engine.NewTxContext()
		_ = txctx.Envelope.SetJobUUID("jobUUID")
		_ = txctx.Envelope.SetTxHashString("0xd41551c714c8ec769d2edad9adc250ae955d263da161bf59142b7500eea6715e")
		txctx.Logger = log.NewEntry(log.New())

		expectedJobUpdate := &txschedulertypes.UpdateJobRequest{
			Transaction: &entities.ETHTransaction{
				Hash: "0xe41551c714c8ec769d2edad9adc250ae955d263da161bf59142b7500eea6715e",
			},
			Status:  utils.StatusWarning,
			Message: "expected transaction hash 0xd41551c714c8ec769d2edad9adc250ae955d263da161bf59142b7500eea6715e, but got 0xe41551c714c8ec769d2edad9adc250ae955d263da161bf59142b7500eea6715e. Overriding",
		}

		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, jobUUID string, request *txschedulertypes.UpdateJobRequest) (*txschedulertypes.JobResponse, error) {
				_ = txctx.Envelope.SetTxHashString("0xe41551c714c8ec769d2edad9adc250ae955d263da161bf59142b7500eea6715e")
				return nil, nil
			})
		schedulerClient.EXPECT().UpdateJob(txctx.Context(), txctx.Envelope.GetJobUUID(), expectedJobUpdate).
			Return(&txschedulertypes.JobResponse{}, nil)

		RawTxStore(schedulerClient)(txctx)

		assert.Empty(t, txctx.Envelope.Error())
	})

	t.Run("should abort if update fails on PENDING", func(t *testing.T) {
		txctx := engine.NewTxContext()
		_ = txctx.Envelope.SetJobUUID("jobUUID")
		txctx.Logger = log.NewEntry(log.New())

		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), txctx.Envelope.GetJobUUID(), gomock.AssignableToTypeOf(&txschedulertypes.UpdateJobRequest{})).
			Return(nil, fmt.Errorf("error"))

		RawTxStore(schedulerClient)(txctx)

		assert.Len(t, txctx.Envelope.GetErrors(), 1)
	})

	t.Run("should set status to RECOVERING if txctx contains errors", func(t *testing.T) {
		txctx := engine.NewTxContext()
		_ = txctx.Envelope.SetJobUUID("jobUUID")
		txctx.Logger = log.NewEntry(log.New())
		_ = txctx.AbortWithError(fmt.Errorf("error"))

		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), txctx.Envelope.GetJobUUID(), gomock.Any()).
			Return(&txschedulertypes.JobResponse{}, nil)
		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), txctx.Envelope.GetJobUUID(), &txschedulertypes.UpdateJobRequest{
				Status: utils.StatusRecovering,
				Message: fmt.Sprintf(
					"transaction attempt with nonce %v and sender %v failed with error: %v",
					txctx.Envelope.GetNonceString(),
					txctx.Envelope.GetFromString(),
					txctx.Envelope.Error(),
				),
			}).
			Return(&txschedulertypes.JobResponse{}, nil)

		RawTxStore(schedulerClient)(txctx)
	})

	t.Run("should return if update fails on RECOVERING", func(t *testing.T) {
		txctx := engine.NewTxContext()
		_ = txctx.Envelope.SetJobUUID("jobUUID")
		txctx.Logger = log.NewEntry(log.New())
		_ = txctx.AbortWithError(fmt.Errorf("error"))

		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), txctx.Envelope.GetJobUUID(), gomock.Any()).
			Return(&txschedulertypes.JobResponse{}, nil)

		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), txctx.Envelope.GetJobUUID(), &txschedulertypes.UpdateJobRequest{
				Status: utils.StatusRecovering,
				Message: fmt.Sprintf(
					"transaction attempt with nonce %v and sender %v failed with error: %v",
					txctx.Envelope.GetNonceString(),
					txctx.Envelope.GetFromString(),
					txctx.Envelope.Error(),
				),
			}).
			Return(nil, fmt.Errorf("error"))

		RawTxStore(schedulerClient)(txctx)
	})
}
