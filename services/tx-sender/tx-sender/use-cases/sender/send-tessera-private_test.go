// +build unit

package sender

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender/tx-sender/use-cases/mocks"
)

func TestSendTesseraPrivate_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ec := mock2.NewMockQuorumTransactionSender(ctrl)
	crafter :=  mocks.NewMockCraftTransactionUseCase(ctrl)
	jobClient := mock.NewMockJobClient(ctrl)
	chainRegistryURL := "chainRegistryURL:8081"
	ctx := context.Background()

	usecase := NewSendTesseraPrivateTxUseCase(ec, crafter, jobClient, chainRegistryURL)

	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.PrivateFrom = "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=" 
		job.Transaction.Data = "0xfe378324abcde723"
		enclaveKey := "0xenclaveKey"

		proxyURL := utils.GetProxyTesseraURL(chainRegistryURL, job.ChainUUID)
		
		data, _ := hexutil.Decode(job.Transaction.Data)
		crafter.EXPECT().Execute(gomock.Any(), job)
		ec.EXPECT().StoreRaw(gomock.Any(), proxyURL, data, job.Transaction.PrivateFrom).Return(enclaveKey, nil)

		jobClient.EXPECT().UpdateJob(gomock.Any(), job.UUID, &api.UpdateJobRequest{
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

		proxyURL := utils.GetProxyTesseraURL(chainRegistryURL, job.ChainUUID)
		crafter.EXPECT().Execute(gomock.Any(), job)
		expectedErr := errors.InternalError("internal_err")
		data, _ := hexutil.Decode(job.Transaction.Data)
		ec.EXPECT().StoreRaw(gomock.Any(), proxyURL, data, job.Transaction.PrivateFrom).Return(enclaveKey, expectedErr)

		err := usecase.Execute(ctx, job)
		assert.Equal(t, err, expectedErr)
	})
}
