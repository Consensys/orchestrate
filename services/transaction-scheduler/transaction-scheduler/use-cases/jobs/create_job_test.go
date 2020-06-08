// +build unit

package jobs

import (
	"context"
	testutils3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
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
	mockTxValidator := mocks2.NewMockTransactionValidator(ctrl)

	mockDB.EXPECT().Begin().Return(mockDBTX, nil).AnyTimes()
	mockDBTX.EXPECT().Transaction().Return(mockTransactionDA).AnyTimes()
	mockDB.EXPECT().Schedule().Return(mockScheduleDA).AnyTimes()
	mockDBTX.EXPECT().Job().Return(mockJobDA).AnyTimes()
	mockDBTX.EXPECT().Log().Return(mockLogDA).AnyTimes()
	mockDBTX.EXPECT().Commit().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Rollback().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Close().Return(nil).AnyTimes()

	usecase := NewCreateJobUseCase(mockDB, mockTxValidator)

	tenantID := "tenantID"

	t.Run("should execute use case successfully", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockTxValidator.EXPECT().ValidateChainExists(ctx, jobEntity.ChainUUID).Return(nil)
		mockScheduleDA.EXPECT().FindOneByUUID(ctx, jobEntity.ScheduleUUID, tenantID).Return(fakeSchedule, nil)
		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)

		assert.Nil(t, err)
	})

	t.Run("should fail with same error if chain is invalid", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils3.FakeJob()

		mockTxValidator.EXPECT().ValidateChainExists(ctx, jobEntity.ChainUUID).Return(expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if cannot fetch selected ScheduleUUID", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID

		mockTxValidator.EXPECT().ValidateChainExists(ctx, jobEntity.ChainUUID).Return(nil)
		mockScheduleDA.EXPECT().FindOneByUUID(ctx, jobEntity.ScheduleUUID, tenantID).Return(nil, expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if cannot insert a Transaction fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockTxValidator.EXPECT().ValidateChainExists(ctx, jobEntity.ChainUUID).Return(nil)
		mockScheduleDA.EXPECT().FindOneByUUID(ctx, jobEntity.ScheduleUUID, tenantID).Return(fakeSchedule, nil)
		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if cannot insert a Job fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockTxValidator.EXPECT().ValidateChainExists(ctx, jobEntity.ChainUUID).Return(nil)
		mockScheduleDA.EXPECT().FindOneByUUID(ctx, jobEntity.ScheduleUUID, tenantID).Return(fakeSchedule, nil)
		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if cannot insert a Log fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockTxValidator.EXPECT().ValidateChainExists(ctx, jobEntity.ChainUUID).Return(nil)
		mockScheduleDA.EXPECT().FindOneByUUID(ctx, jobEntity.ScheduleUUID, tenantID).Return(fakeSchedule, nil)
		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})
}
