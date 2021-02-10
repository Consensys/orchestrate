// +build unit

package jobs

import (
	"context"
	"fmt"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	testutils3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/mocks"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models/testutils"
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
	mockDBTX.EXPECT().Transaction().Return(mockTransactionDA).AnyTimes()
	mockDB.EXPECT().Schedule().Return(mockScheduleDA).AnyTimes()
	mockDB.EXPECT().Account().Return(mockAccountDA).AnyTimes()
	mockDBTX.EXPECT().Job().Return(mockJobDA).AnyTimes()
	mockDBTX.EXPECT().Log().Return(mockLogDA).AnyTimes()
	mockDBTX.EXPECT().Commit().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Rollback().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Close().Return(nil).AnyTimes()

	usecase := NewCreateJobUseCase(mockDB, mockGetChainUC)

	tenantID := "tenantID"
	tenants := []string{tenantID}
	fakeChain := testutils3.FakeChain()

	t.Run("should execute use case successfully", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, tenants).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From, tenants).Return(nil, nil)
		mockScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.ScheduleUUID, tenants).Return(fakeSchedule, nil)
		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		_, err := usecase.Execute(context.Background(), jobEntity, tenants)

		assert.NoError(t, err)
	})

	t.Run("should execute use case successfully for child job", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.InternalData.ParentJobUUID = "myParentJobUUID"
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)
		parentJobModel := testutils2.FakeJobModel(fakeSchedule.ID)
		parentJobModel.Logs[0].Status = entities.StatusPending

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, tenants).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From, tenants).Return(nil, nil)
		mockScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.ScheduleUUID, tenants).Return(fakeSchedule, nil)
		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.InternalData.ParentJobUUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.InternalData.ParentJobUUID, tenants).Return(parentJobModel, nil)
		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		_, err := usecase.Execute(context.Background(), jobEntity, tenants)

		assert.NoError(t, err)
	})

	t.Run("should fail with InvalidParameterError if chain is not found", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, tenants).Return(nil, errors.NotFoundError("error"))

		_, err := usecase.Execute(context.Background(), jobEntity, tenants)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with same error if chain is invalid", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")
		jobEntity := testutils3.FakeJob()

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, tenants).Return(nil, expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, tenants)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with InvalidParameterError if account does not exist", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, tenants).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From, tenants).Return(nil, errors.NotFoundError("error"))

		_, err := usecase.Execute(context.Background(), jobEntity, tenants)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with InvalidParameterError if schedule is not found", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, tenants).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From, tenants).Return(nil, nil)
		mockScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.ScheduleUUID, tenants).Return(nil, errors.NotFoundError("error"))

		_, err := usecase.Execute(context.Background(), jobEntity, tenants)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with same error if cannot fetch selected ScheduleUUID", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, tenants).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From, tenants).Return(nil, nil)
		mockScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.ScheduleUUID, tenants).Return(nil, expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, tenants)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if cannot insert a Transaction fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, tenants).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From, tenants).Return(nil, nil)
		mockScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.ScheduleUUID, tenants).Return(fakeSchedule, nil)
		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, tenants)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with InvalidState if parentJob is not PENDING when creating a child job", func(t *testing.T) {
		expectedErr := errors.InvalidStateError("cannot create a child job in a finalized schedule")
		jobEntity := testutils3.FakeJob()
		jobEntity.InternalData.ParentJobUUID = "myParentJobUUID"
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		parentJobModel := testutils2.FakeJobModel(fakeSchedule.ID)
		parentJobModel.Logs[0].Status = entities.StatusMined

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, tenants).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From, tenants).Return(nil, nil)
		mockScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.ScheduleUUID, tenants).Return(fakeSchedule, nil)
		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.InternalData.ParentJobUUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.InternalData.ParentJobUUID, tenants).Return(parentJobModel, nil)

		_, err := usecase.Execute(context.Background(), jobEntity, tenants)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if cannot insert a Job fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, tenants).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From, tenants).Return(nil, nil)
		mockScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.ScheduleUUID, tenants).Return(fakeSchedule, nil)
		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, tenants)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})

	t.Run("should fail with same error if cannot insert a Log fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		jobEntity := testutils3.FakeJob()
		fakeSchedule := testutils2.FakeSchedule(tenantID)
		fakeSchedule.ID = 1
		fakeSchedule.UUID = jobEntity.ScheduleUUID
		jobModel := parsers.NewJobModelFromEntities(jobEntity, &fakeSchedule.ID)

		mockGetChainUC.EXPECT().Execute(gomock.Any(), jobEntity.ChainUUID, tenants).Return(fakeChain, nil)
		mockAccountDA.EXPECT().FindOneByAddress(gomock.Any(), jobEntity.Transaction.From, tenants).Return(nil, nil)
		mockScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.ScheduleUUID, tenants).Return(fakeSchedule, nil)
		mockTransactionDA.EXPECT().Insert(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, tenants)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})
}
