// +build unit

package usecases

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/txscheduler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client/mock"
	"testing"
)

func TestCreateChildJob_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	initialGasPrice := "1000000000"
	ctx := context.Background()

	mockTxSchedulerClient := mock.NewMockTransactionSchedulerClient(ctrl)

	usecase := NewCreateChildJobUseCase(mockTxSchedulerClient)

	t.Run("should do nothing if status of the job is not PENDING", func(t *testing.T) {
		parentJob := testutils.FakeJob()
		parentJobResponse := testutils.FakeJobResponse()
		jobResponses := []*types.JobResponse{parentJobResponse}

		mockTxSchedulerClient.EXPECT().SearchJob(ctx, &entities.JobFilters{
			ChainUUID:     parentJob.ChainUUID,
			ParentJobUUID: parentJob.UUID,
		}).Return(jobResponses, nil)

		childJobUUID, err := usecase.Execute(ctx, parentJob)
		assert.NoError(t, err)
		assert.Empty(t, childJobUUID)
	})

	t.Run("should create a new child job if the parent job status is PENDING", func(t *testing.T) {
		parentJob := testutils.FakeJob()
		childJobResponse := testutils.FakeJobResponse()
		parentJobResponse := testutils.FakeJobResponse()
		parentJobResponse.Status = utils.StatusPending
		parentJobResponse.Transaction.GasPrice = initialGasPrice
		jobResponses := []*types.JobResponse{parentJobResponse}

		mockTxSchedulerClient.EXPECT().SearchJob(ctx, gomock.Any()).Return(jobResponses, nil)
		mockTxSchedulerClient.EXPECT().CreateJob(ctx, gomock.Any()).Return(childJobResponse, nil)
		mockTxSchedulerClient.EXPECT().StartJob(ctx, childJobResponse.UUID).Return(nil)

		childJobUUID, err := usecase.Execute(ctx, parentJob)
		assert.NoError(t, err)
		assert.NotEmpty(t, childJobUUID)
	})

	t.Run("should send the same job if job is a raw transaction", func(t *testing.T) {
		parentJob := testutils.FakeJob()
		childJobResponse := testutils.FakeJobResponse()
		parentJobResponse := testutils.FakeJobResponse()
		parentJobResponse.Transaction.Raw = "0xraw"
		parentJobResponse.Type = utils.EthereumRawTransaction
		parentJobResponse.Status = utils.StatusPending
		jobResponses := []*types.JobResponse{parentJobResponse}

		mockTxSchedulerClient.EXPECT().SearchJob(ctx, gomock.Any()).Return(jobResponses, nil)
		mockTxSchedulerClient.EXPECT().CreateJob(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, req *types.CreateJobRequest) (*types.JobResponse, error) {
				assert.Equal(t, parentJobResponse.Transaction.Raw, req.Transaction.Raw)
				return childJobResponse, nil
			})
		mockTxSchedulerClient.EXPECT().StartJob(ctx, childJobResponse.UUID).Return(nil)

		childJobUUID, err := usecase.Execute(ctx, parentJob)
		assert.NoError(t, err)
		assert.NotEmpty(t, childJobUUID)
	})

	t.Run("should create a new child job by increasing the gasPrice by Increment", func(t *testing.T) {
		parentJob := testutils.FakeJob()
		childJobResponse := testutils.FakeJobResponse()

		parentJobResponse := testutils.FakeJobResponse()
		parentJobResponse.Status = utils.StatusPending
		parentJobResponse.Transaction.GasPrice = initialGasPrice
		parentJobResponse.Transaction.Nonce = "1"
		parentJobResponse.Annotations.GasPricePolicy.RetryPolicy.Increment = 0.06
		parentJobResponse.Annotations.GasPricePolicy.RetryPolicy.Limit = 0.12

		mockTxSchedulerClient.EXPECT().SearchJob(ctx, gomock.Any()).Return([]*types.JobResponse{parentJobResponse, childJobResponse}, nil)
		mockTxSchedulerClient.EXPECT().CreateJob(ctx, gomock.Any()).
			DoAndReturn(func(timeoutCtx context.Context, req *types.CreateJobRequest) (*types.JobResponse, error) {
				assert.Equal(t, "1120000000", req.Transaction.GasPrice)
				assert.Equal(t, parentJobResponse.Transaction.Nonce, req.Transaction.Nonce)
				return childJobResponse, nil
			})
		mockTxSchedulerClient.EXPECT().StartJob(ctx, childJobResponse.UUID).Return(nil)

		childJobUUID, err := usecase.Execute(ctx, parentJob)
		assert.NoError(t, err)
		assert.NotEmpty(t, childJobUUID)
	})

	t.Run("should create a new child job by increasing the gasPrice and not exceed the limit", func(t *testing.T) {
		parentJob := testutils.FakeJob()
		childJobResponse := testutils.FakeJobResponse()

		parentJobResponse := testutils.FakeJobResponse()
		parentJobResponse.Status = utils.StatusPending
		parentJobResponse.Transaction.GasPrice = initialGasPrice
		parentJobResponse.Transaction.Nonce = "1"
		parentJobResponse.Annotations.GasPricePolicy.RetryPolicy.Increment = 0.06
		parentJobResponse.Annotations.GasPricePolicy.RetryPolicy.Limit = 0.05

		mockTxSchedulerClient.EXPECT().SearchJob(ctx, gomock.Any()).Return([]*types.JobResponse{parentJobResponse, childJobResponse}, nil)
		mockTxSchedulerClient.EXPECT().CreateJob(ctx, gomock.Any()).
			DoAndReturn(func(timeoutCtx context.Context, req *types.CreateJobRequest) (*types.JobResponse, error) {
				assert.Equal(t, "1050000000", req.Transaction.GasPrice)
				assert.Equal(t, parentJobResponse.Transaction.Nonce, req.Transaction.Nonce)
				return childJobResponse, nil
			})
		mockTxSchedulerClient.EXPECT().StartJob(ctx, childJobResponse.UUID).Return(nil)

		childJobUUID, err := usecase.Execute(ctx, parentJob)
		assert.NoError(t, err)
		assert.NotEmpty(t, childJobUUID)
	})
}
