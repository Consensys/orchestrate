// +build unit

package jobs

import (
	"context"
	"fmt"
	"testing"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	mocks2 "github.com/consensys/orchestrate/services/api/business/use-cases/mocks"

	"github.com/consensys/orchestrate/pkg/errors"
	testutils3 "github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	"github.com/consensys/orchestrate/services/api/store/mocks"
	testutils2 "github.com/consensys/orchestrate/services/api/store/models/testutils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateJob_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockDBTX := mocks.NewMockTx(ctrl)
	mockScheduleDA := mocks.NewMockScheduleAgent(ctrl)
	mockAccountDA := mocks.NewMockAccountAgent(ctrl)
	mockTransactionDA := mocks.NewMockTransactionAgent(ctrl)
	mockJobDA := mocks.NewMockJobAgent(ctrl)
	mockLogDA := mocks.NewMockLogAgent(ctrl)
	mockGetChainUC := mocks2.NewMockGetChainUseCase(ctrl)

	mockDB.EXPECT().Begin().Return(mockDBTX, nil).AnyTimes()
	mockDB.EXPECT().Schedule().Return(mockScheduleDA).AnyTimes()
	mockDB.EXPECT().Account().Return(mockAccountDA).AnyTimes()
	mockDB.EXPECT().Transaction().Return(mockTransactionDA).AnyTimes()
	mockDBTX.EXPECT().Job().Return(mockJobDA).AnyTimes()
	mockDBTX.EXPECT().Log().Return(mockLogDA).AnyTimes()
	mockDBTX.EXPECT().Transaction().Return(mockTransactionDA).AnyTimes()
	mockDBTX.EXPECT().Commit().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Rollback().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Close().Return(nil).AnyTimes()

	userInfo := multitenancy.NewUserInfo("tenantOne", "username")
	usecase := NewCreateJobUseCase(mockDB, mockGetChainUC)
	fakeChain := testutils3.FakeChain()

	t.Run("should execute use case successfully", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(userInfo.TenantID, userInfo.Username)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, userInfo).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From.String(), userInfo.AllowedTenants, userInfo.Username).
			Return(nil, nil)
		mockScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.ScheduleUUID, userInfo.AllowedTenants, userInfo.Username).
			Return(fakeSchedule, nil)
		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		_, err := usecase.Execute(context.Background(), jobEntity, userInfo)

		assert.NoError(t, err)
	})

	t.Run("should execute use case successfully for child job", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.InternalData.ParentJobUUID = "myParentJobUUID"
		fakeSchedule := testutils2.FakeSchedule(userInfo.TenantID, userInfo.Username)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)
		parentJobModel := testutils2.FakeJobModel(fakeSchedule.ID)
		parentJobModel.Status = entities.StatusPending
		parentJobModel.Logs[0].Status = entities.StatusPending

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, userInfo).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From.String(), userInfo.AllowedTenants, userInfo.Username).
			Return(nil, nil)
		mockScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.ScheduleUUID, userInfo.AllowedTenants, userInfo.Username).
			Return(fakeSchedule, nil)
		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.InternalData.ParentJobUUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.InternalData.ParentJobUUID, userInfo.AllowedTenants, 
			userInfo.Username, false).Return(parentJobModel, nil)
		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		_, err := usecase.Execute(context.Background(), jobEntity, userInfo)

		assert.NoError(t, err)
	})

	t.Run("should fail with InvalidParameterError if chain is not found", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, userInfo).
			Return(nil, errors.NotFoundError("error"))

		_, err := usecase.Execute(context.Background(), jobEntity, userInfo)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with same error if chain is invalid", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")
		jobEntity := testutils3.FakeJob()

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, userInfo).Return(nil, expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, userInfo)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with InvalidParameterError if account does not exist", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, userInfo).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From.String(), userInfo.AllowedTenants, userInfo.Username).Return(nil, errors.NotFoundError("error"))

		_, err := usecase.Execute(context.Background(), jobEntity, userInfo)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with InvalidParameterError if schedule is not found", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(userInfo.TenantID, userInfo.Username)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, userInfo).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From.String(), userInfo.AllowedTenants, userInfo.Username).Return(nil, nil)
		mockScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.ScheduleUUID, userInfo.AllowedTenants, userInfo.Username).Return(nil, errors.NotFoundError("error"))

		_, err := usecase.Execute(context.Background(), jobEntity, userInfo)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with same error if cannot fetch selected ScheduleUUID", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(userInfo.TenantID, userInfo.Username)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, userInfo).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From.String(), userInfo.AllowedTenants, userInfo.Username).Return(nil, nil)
		mockScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.ScheduleUUID, userInfo.AllowedTenants, userInfo.Username).Return(nil, expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, userInfo)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if cannot insert a Transaction fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(userInfo.TenantID, userInfo.Username)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, userInfo).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From.String(), userInfo.AllowedTenants, userInfo.Username).Return(nil, nil)
		mockScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.ScheduleUUID, userInfo.AllowedTenants, userInfo.Username).Return(fakeSchedule, nil)
		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, userInfo)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with InvalidState if parentJob is not PENDING when creating a child job", func(t *testing.T) {
		expectedErr := errors.InvalidStateError("cannot create a child job in a finalized schedule")
		jobEntity := testutils3.FakeJob()
		jobEntity.InternalData.ParentJobUUID = "myParentJobUUID"
		fakeSchedule := testutils2.FakeSchedule(userInfo.TenantID, userInfo.Username)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		parentJobModel := testutils2.FakeJobModel(fakeSchedule.ID)
		parentJobModel.Status = entities.StatusMined

		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, userInfo).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From.String(), userInfo.AllowedTenants, userInfo.Username).Return(nil, nil)
		mockScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.ScheduleUUID, userInfo.AllowedTenants, userInfo.Username).Return(fakeSchedule, nil)
		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.InternalData.ParentJobUUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.InternalData.ParentJobUUID, userInfo.AllowedTenants, userInfo.Username, false).Return(parentJobModel, nil)

		_, err := usecase.Execute(context.Background(), jobEntity, userInfo)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if cannot insert a Job fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(userInfo.TenantID, userInfo.Username)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, userInfo).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From.String(), userInfo.AllowedTenants, userInfo.Username).Return(nil, nil)
		mockScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.ScheduleUUID, userInfo.AllowedTenants, userInfo.Username).Return(fakeSchedule, nil)
		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, userInfo)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if cannot insert a Log fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(userInfo.TenantID, userInfo.Username)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, userInfo).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From.String(), userInfo.AllowedTenants, userInfo.Username).Return(nil, nil)
		mockScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.ScheduleUUID, userInfo.AllowedTenants, userInfo.Username).Return(fakeSchedule, nil)
		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, userInfo)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})
}
