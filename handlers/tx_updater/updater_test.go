// +build unit

package txupdater

import (
	"fmt"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client/mock"
	"testing"
)

func TestTransactionUpdater(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	schedulerClient := mock.NewMockTransactionSchedulerClient(mockCtrl)

	t.Run("should do nothing if the tx does not contain errors", func(t *testing.T) {
		txctx := engine.NewTxContext()
		_ = txctx.Envelope.SetID("test")
		txctx.Logger = log.NewEntry(log.New())

		h := TransactionUpdater(schedulerClient)
		h(txctx)
	})

	t.Run("should update the status successfully to FAILED if envelope contains errors", func(t *testing.T) {
		txctx := engine.NewTxContext()
		_ = txctx.Envelope.SetID("test")
		_ = txctx.AbortWithError(fmt.Errorf("error"))
		txctx.Logger = log.NewEntry(log.New())

		schedulerClient.EXPECT().
			UpdateJob(txctx.Context(), txctx.Envelope.GetID(), &types.UpdateJobRequest{
				Status:  utils.StatusFailed,
				Message: txctx.Envelope.Error(),
			}).
			Return(&types.JobResponse{}, nil)

		h := TransactionUpdater(schedulerClient)
		h(txctx)
	})
}
