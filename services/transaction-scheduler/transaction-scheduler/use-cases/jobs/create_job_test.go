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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"
)

func TestCreateJob_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockDBTX := mocks.NewMockTx(ctrl)
	mockScheduleDA := mocks.NewMockScheduleAgent(ctrl)
	mockTransactionDA := mocks.NewMockTransactionAgent(ctrl)
	mockJobDA := mocks.NewMockJobAgent(ctrl)
	mockLogDA := mocks.NewMockLogAgent(ctrl)

	usecase := NewCreateJobUseCase(mockDB)

	tenantID := "tenantID"

	t.Run("should execute use case successfully", func(t *testing.T) {
		jobEntity := testutils.FakeJobEntity()
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockDB.EXPECT().Schedule().Return(mockScheduleDA).Times(1)
		mockScheduleDA.EXPECT().FindOneByUUID(ctx, jobEntity.ScheduleUUID, tenantID).Return(fakeSchedule, nil)

		mockDB.EXPECT().Begin().Return(mockDBTX, nil).Times(1)
		mockDBTX.EXPECT().Transaction().Return(mockTransactionDA).Times(1)
		mockDBTX.EXPECT().Job().Return(mockJobDA).Times(1)
		mockDBTX.EXPECT().Log().Return(mockLogDA).Times(1)
		mockDBTX.EXPECT().Commit().Return(nil).Times(1)
		mockDBTX.EXPECT().Close().Return(nil).Times(1)

		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(nil).Times(1)
		mockJobDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)

		assert.Nil(t, err)
	})

	t.Run("should fail with same error if cannot fetch selected ScheduleUUID", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		jobEntity := testutils.FakeJobEntity()
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID

		mockDB.EXPECT().Schedule().Return(mockScheduleDA).Times(1)
		mockScheduleDA.EXPECT().FindOneByUUID(ctx, jobEntity.ScheduleUUID, tenantID).Return(nil, expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if cannot insert a Transaction fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		jobEntity := testutils.FakeJobEntity()
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockDB.EXPECT().Schedule().Return(mockScheduleDA).Times(1)
		mockScheduleDA.EXPECT().FindOneByUUID(ctx, jobEntity.ScheduleUUID, tenantID).Return(fakeSchedule, nil)

		mockDB.EXPECT().Begin().Return(mockDBTX, nil).Times(1)
		mockDBTX.EXPECT().Transaction().Return(mockTransactionDA).Times(1)
		mockDBTX.EXPECT().Rollback().Return(nil).Times(1)
		mockDBTX.EXPECT().Close().Return(nil).Times(1)

		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(expectedErr).Times(1)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if cannot insert a Job fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		jobEntity := testutils.FakeJobEntity()
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockDB.EXPECT().Schedule().Return(mockScheduleDA).Times(1)
		mockScheduleDA.EXPECT().FindOneByUUID(ctx, jobEntity.ScheduleUUID, tenantID).Return(fakeSchedule, nil)

		mockDB.EXPECT().Begin().Return(mockDBTX, nil).Times(1)
		mockDBTX.EXPECT().Transaction().Return(mockTransactionDA).Times(1)
		mockDBTX.EXPECT().Job().Return(mockJobDA).Times(1)
		mockDBTX.EXPECT().Rollback().Return(nil).Times(1)
		mockDBTX.EXPECT().Close().Return(nil).Times(1)

		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(nil).Times(1)
		mockJobDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr).Times(1)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if cannot insert a Log fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		jobEntity := testutils.FakeJobEntity()
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockDB.EXPECT().Schedule().Return(mockScheduleDA).Times(1)
		mockScheduleDA.EXPECT().FindOneByUUID(ctx, jobEntity.ScheduleUUID, tenantID).Return(fakeSchedule, nil)

		mockDB.EXPECT().Begin().Return(mockDBTX, nil).Times(1)
		mockDBTX.EXPECT().Transaction().Return(mockTransactionDA).Times(1)
		mockDBTX.EXPECT().Job().Return(mockJobDA).Times(1)
		mockDBTX.EXPECT().Log().Return(mockLogDA).Times(1)
		mockDBTX.EXPECT().Rollback().Return(nil).Times(1)
		mockDBTX.EXPECT().Close().Return(nil).Times(1)

		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(nil).Times(1)
		mockJobDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr).Times(1)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})
}
