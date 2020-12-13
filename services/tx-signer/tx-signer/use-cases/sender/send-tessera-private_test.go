// +build unit

package sender

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/txscheduler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	mock3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client/mock"
)

func TestSendTesseraPrivate_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ec := mock2.NewMockQuorumTransactionSender(ctrl)
	txSchedulerClient := mock3.NewMockTransactionSchedulerClient(ctrl)
	chainRegistryURL := "chainRegistryURL:8081"
	ctx := context.Background()

	usecase := NewSendTesseraPrivateTxUseCase(ec, txSchedulerClient, chainRegistryURL)

	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.PrivateFrom = "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=" 
		job.Transaction.Data = "0xfe378324abcde723"
		enclaveKey := "0xenclaveKey"

		proxyURL := fmt.Sprintf("%s/tessera/%s", chainRegistryURL, job.ChainUUID)
		
		data, _ := hexutil.Decode(job.Transaction.Data)
		ec.EXPECT().StoreRaw(ctx, proxyURL, data, job.Transaction.PrivateFrom).Return(enclaveKey, nil)

		txSchedulerClient.EXPECT().UpdateJob(ctx, job.UUID, &txschedulertypes.UpdateJobRequest{
			Status: utils.StatusStored,
			Transaction: job.Transaction,
		})

		err := usecase.Execute(ctx, job)
		assert.NoError(t, err)
		assert.Equal(t, job.Transaction.EnclaveKey, enclaveKey)
	})
	
	t.Run("should fail with same error executing use case if storeRaw fails", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.PrivateFrom = "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=" 
		job.Transaction.Data = "0xfe378324abcde723"
		enclaveKey := "0xenclaveKey"

		proxyURL := fmt.Sprintf("%s/tessera/%s", chainRegistryURL, job.ChainUUID)
		
		expectedErr := errors.InternalError("internal_err")
		data, _ := hexutil.Decode(job.Transaction.Data)
		ec.EXPECT().StoreRaw(ctx, proxyURL, data, job.Transaction.PrivateFrom).Return(enclaveKey, expectedErr)

		err := usecase.Execute(ctx, job)
		assert.Equal(t, err, expectedErr)
	})
}
