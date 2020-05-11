// +build unit

package jobs

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/interfaces/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types/testutils"
	"testing"
)

func TestCreateJob_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJobDA := mocks.NewMockJobAgent(ctrl)
	mockLogDA := mocks.NewMockLogAgent(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockDB := mocks.NewMockDB(ctrl)

	usecase := NewCreateJobUseCase(mockDB)
	ctx := context.Background()

	t.Run("should execute use case successfully", func(t *testing.T) {
		jobRequest := testutils.FakeJobRequest()
		expectedResponse := &types.JobResponse{
			UUID:        "testJobUUID",
			Transaction: jobRequest.Transaction,
			Status:      types.JobStatusCreated,
		}

		mockDB.EXPECT().Begin().Return(mockTx, nil)
		mockTx.EXPECT().Job().Return(mockJobDA)
		mockTx.EXPECT().Log().Return(mockLogDA)
		mockJobDA.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, job *models.Job) error {
			job.UUID = "testJobUUID"
			return nil
		})
		mockLogDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		mockTx.EXPECT().Commit().Return(nil)

		jobResponse, err := usecase.Execute(context.Background(), jobRequest)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse.UUID, jobResponse.UUID)
		assert.Equal(t, expectedResponse.Transaction, jobResponse.Transaction)
		assert.Equal(t, expectedResponse.Status, jobResponse.Status)
	})

	t.Run("should fail with InvalidParameterError error if it fails to validate request", func(t *testing.T) {
		jobRequest := testutils.FakeJobRequest()
		jobRequest.Type = ""

		jobResponse, err := usecase.Execute(ctx, jobRequest)
		assert.True(t, errors.IsInvalidParameterError(err))
		assert.Nil(t, jobResponse)
	})

	t.Run("should fail with same error if Begin fails", func(t *testing.T) {
		jobRequest := testutils.FakeJobRequest()
		expectedErr := errors.PostgresConnectionError("error")

		mockDB.EXPECT().Begin().Return(nil, expectedErr)

		jobResponse, err := usecase.Execute(context.Background(), jobRequest)

		assert.Nil(t, jobResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if Insert of job fails", func(t *testing.T) {
		jobRequest := testutils.FakeJobRequest()
		expectedErr := errors.PostgresConnectionError("error")

		mockDB.EXPECT().Begin().Return(mockTx, nil)
		mockTx.EXPECT().Job().Return(mockJobDA)
		mockJobDA.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)
		jobResponse, err := usecase.Execute(context.Background(), jobRequest)

		assert.Nil(t, jobResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if Insert of log fails", func(t *testing.T) {
		jobRequest := testutils.FakeJobRequest()
		expectedErr := errors.PostgresConnectionError("error")

		mockDB.EXPECT().Begin().Return(mockTx, nil)
		mockTx.EXPECT().Job().Return(mockJobDA)
		mockTx.EXPECT().Log().Return(mockLogDA)
		mockJobDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		mockLogDA.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)

		jobResponse, err := usecase.Execute(context.Background(), jobRequest)

		assert.Nil(t, jobResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if Commit fails", func(t *testing.T) {
		jobRequest := testutils.FakeJobRequest()
		expectedErr := errors.PostgresConnectionError("error")

		mockDB.EXPECT().Begin().Return(mockTx, nil)
		mockTx.EXPECT().Job().Return(mockJobDA)
		mockTx.EXPECT().Log().Return(mockLogDA)
		mockJobDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		mockLogDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		mockTx.EXPECT().Commit().Return(expectedErr)

		jobResponse, err := usecase.Execute(context.Background(), jobRequest)

		assert.Nil(t, jobResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})
}
