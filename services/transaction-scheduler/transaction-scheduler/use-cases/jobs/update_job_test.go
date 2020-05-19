// +build unit

package jobs

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types/testutils"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/orm/mocks"
)

func TestUpdateJob_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockJobDA := mocks.NewMockJobAgent(ctrl)
	mockORM := mocks2.NewMockORM(ctrl)
	

	mockDB.EXPECT().Job().Return(mockJobDA).AnyTimes()

	usecase := NewUpdateJobUseCase(mockDB, mockORM)
	tenantID := "tenantID"

	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils2.FakeJob(0)
		jobRequest := testutils.FakeJobUpdateRequest()
		expectedResponse := &types.JobResponse{
			UUID:        job.UUID,
			Transaction: jobRequest.Transaction,
			Status:      types.JobStatusCreated,
		}

		mockJobDA.EXPECT().FindOneByUUID(ctx, job.UUID, tenantID).Return(job, nil)
		mockORM.EXPECT().InsertOrUpdateJob(gomock.Any(), gomock.Eq(mockDB), gomock.Any()).Return(nil)

		jobResponse, err := usecase.Execute(ctx, job.UUID, jobRequest, tenantID)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse.Transaction, jobResponse.Transaction)
		assert.Equal(t, expectedResponse.Status, jobResponse.Status)
	})
	
	t.Run("should fail with the same error if find one fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		job := testutils2.FakeJob(0)
		jobRequest := testutils.FakeJobUpdateRequest()

		mockJobDA.EXPECT().FindOneByUUID(ctx, job.UUID, tenantID).Return(job, expectedErr)

		_, err := usecase.Execute(ctx, job.UUID, jobRequest, tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	 t.Run("should fail with the same error if insert or update job fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		job := testutils2.FakeJob(0)
		jobRequest := testutils.FakeJobUpdateRequest()

		mockJobDA.EXPECT().FindOneByUUID(ctx, job.UUID, tenantID).Return(job, nil)
		mockORM.EXPECT().InsertOrUpdateJob(gomock.Any(), gomock.Eq(mockDB), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(ctx, job.UUID, jobRequest, tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})
}
