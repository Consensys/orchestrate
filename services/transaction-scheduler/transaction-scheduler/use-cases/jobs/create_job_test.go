// +build unit

package jobs

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types/testutils"
	"testing"
)

func TestCreateJob_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJobDA := mocks.NewMockJobAgent(ctrl)
	usecase := NewCreateJobUseCase(mockJobDA)
	ctx := context.Background()

	t.Run("should execute use case successfully", func(t *testing.T) {
		jobRequest := testutils.FakeJobRequest()
		expectedResponse := &types.JobResponse{
			UUID:        "testUUID",
			Transaction: jobRequest.Transaction,
			Status:      types.LogStatusCreated,
		}

		mockJobDA.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, job *models.Job) error {
			job.UUID = "testUUID"
			return nil
		})
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

	t.Run("should fail with same error if Insert fails", func(t *testing.T) {
		jobRequest := testutils.FakeJobRequest()
		expectedErr := errors.PostgresConnectionError("error")

		mockJobDA.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)
		jobResponse, err := usecase.Execute(context.Background(), jobRequest)

		assert.Nil(t, jobResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})
}
