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
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/orm/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types/testutils"
)

func TestCreateJob_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockScheduleDA := mocks.NewMockScheduleAgent(ctrl)
	mockORM := mocks2.NewMockORM(ctrl)
	

	mockDB.EXPECT().Schedule().Return(mockScheduleDA).AnyTimes()

	usecase := NewCreateJobUseCase(mockDB, mockORM)
	fakeSchedule := testutils2.FakeSchedule("")
	fakeSchedule.ID = 1
	tenantID := "tenantID"

	t.Run("should execute use case successfully", func(t *testing.T) {
		jobRequest := testutils.FakeJobRequest()
		expectedResponse := &types.JobResponse{
			UUID:        "testJobUUID",
			Transaction: jobRequest.Transaction,
			Status:      types.JobStatusCreated,
		}

		mockScheduleDA.EXPECT().FindOneByUUID(ctx, jobRequest.ScheduleUUID, tenantID).Return(fakeSchedule, nil)
		mockORM.EXPECT().InsertOrUpdateJob(gomock.Any(), gomock.Eq(mockDB), gomock.Any()).Return(nil)

		jobResponse, err := usecase.Execute(context.Background(), jobRequest, tenantID)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse.Transaction, jobResponse.Transaction)
		assert.Equal(t, expectedResponse.Status, jobResponse.Status)
	})

	t.Run("should fail with InvalidParameterError error if it fails to validate request", func(t *testing.T) {
		jobRequest := testutils.FakeJobRequest()
		jobRequest.Type = ""

		jobResponse, err := usecase.Execute(ctx, jobRequest, tenantID)
		assert.True(t, errors.IsInvalidParameterError(err))
		assert.Nil(t, jobResponse)
	})

	t.Run("should fail with same error when schedule in not found", func(t *testing.T) {
		jobRequest := testutils.FakeJobRequest()
		expectedErr := errors.NotFoundError("scheduleNotFount")

		mockScheduleDA.EXPECT().FindOneByUUID(ctx, jobRequest.ScheduleUUID, tenantID).Return(nil, expectedErr)

		jobResponse, err := usecase.Execute(context.Background(), jobRequest, tenantID)

		assert.Nil(t, jobResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if InsertOrUpdate fails", func(t *testing.T) {
		jobRequest := testutils.FakeJobRequest()
		expectedErr := errors.PostgresConnectionError("error")

		mockScheduleDA.EXPECT().FindOneByUUID(ctx, jobRequest.ScheduleUUID, tenantID).Return(fakeSchedule, nil)
		mockORM.EXPECT().InsertOrUpdateJob(gomock.Any(), gomock.Eq(mockDB), gomock.Any()).Return(expectedErr)

		jobResponse, err := usecase.Execute(context.Background(), jobRequest, tenantID)

		assert.Nil(t, jobResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	// @TODO Ensure tenantID corresponds to schedule tenantID otherwise it throws NotAuth ERR 
}
