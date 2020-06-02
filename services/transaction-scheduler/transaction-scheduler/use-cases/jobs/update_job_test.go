// +build unit

package jobs

import (
	"context"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"
)

func TestUpdateJob_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockDBTX := mocks.NewMockTx(ctrl)
	mockTransactionDA := mocks.NewMockTransactionAgent(ctrl)
	mockJobDA := mocks.NewMockJobAgent(ctrl)
	mockLogDA := mocks.NewMockLogAgent(ctrl)

	mockDB.EXPECT().Job().Return(mockJobDA).AnyTimes()
	mockDB.EXPECT().Begin().Return(mockDBTX, nil).AnyTimes()
	mockDBTX.EXPECT().Transaction().Return(mockTransactionDA).AnyTimes()
	mockDBTX.EXPECT().Job().Return(mockJobDA).AnyTimes()
	mockDBTX.EXPECT().Log().Return(mockLogDA).AnyTimes()
	mockDBTX.EXPECT().Commit().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Rollback().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Close().Return(nil).AnyTimes()

	usecase := NewUpdateJobUseCase(mockDB)

	tenantID := "tenantID"
	newStatus := entities.JobStatusPending

	t.Run("should execute use case successfully", func(t *testing.T) {
		jobEntity := testutils.FakeJobEntity()
		jobEntity.Status = newStatus
		jobModel := testutils2.FakeJob(0)
		parsers.UpdateJobModelFromEntities(jobModel, jobEntity)

		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, tenantID).Return(jobModel, nil).Times(2)
		mockTransactionDA.EXPECT().Update(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Update(gomock.Any(), jobModel).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)

		assert.Nil(t, err)
	})

	t.Run("should execute use case successfully if transaction is empty", func(t *testing.T) {
		jobEntity := testutils.FakeJobEntity()
		jobEntity.Transaction = nil
		jobEntity.Status = newStatus
		jobModel := testutils2.FakeJob(0)

		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, tenantID).Return(jobModel, nil).Times(2)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		job, err := usecase.Execute(context.Background(), jobEntity, tenantID)

		assert.Nil(t, err)
		assert.Equal(t, job.Status, jobModel.GetStatus())
	})

	t.Run("should execute use case successfully if status is empty", func(t *testing.T) {
		jobEntity := testutils.FakeJobEntity()
		jobEntity.Status = ""
		jobModel := testutils2.FakeJob(0)
		parsers.UpdateJobModelFromEntities(jobModel, jobEntity)

		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, tenantID).Return(jobModel, nil).Times(2)
		mockTransactionDA.EXPECT().Update(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Update(gomock.Any(), jobModel).Return(nil)

		job, err := usecase.Execute(context.Background(), jobEntity, tenantID)

		assert.Nil(t, err)
		assert.Equal(t, job.Status, jobModel.GetStatus())
	})

	t.Run("should fail with the same error if find one fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils.FakeJobEntity()
		jobEntity.Status = newStatus
		jobModel := testutils2.FakeJob(0)
		parsers.UpdateJobModelFromEntities(jobModel, jobEntity)

		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, tenantID).Return(jobModel, expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with the same error if update transaction fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils.FakeJobEntity()
		jobEntity.Status = newStatus
		jobModel := testutils2.FakeJob(0)
		parsers.UpdateJobModelFromEntities(jobModel, jobEntity)

		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, tenantID).Return(jobModel, nil)
		mockTransactionDA.EXPECT().Update(gomock.Any(), jobModel.Transaction).Return(expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with the same error if update job fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils.FakeJobEntity()
		jobEntity.Status = newStatus
		jobModel := testutils2.FakeJob(0)
		parsers.UpdateJobModelFromEntities(jobModel, jobEntity)

		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, tenantID).Return(jobModel, nil)
		mockTransactionDA.EXPECT().Update(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Update(gomock.Any(), jobModel).Return(expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with the same error if insert log fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils.FakeJobEntity()
		jobEntity.Status = newStatus
		jobModel := testutils2.FakeJob(0)
		parsers.UpdateJobModelFromEntities(jobModel, jobEntity)

		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, tenantID).Return(jobModel, nil)
		mockTransactionDA.EXPECT().Update(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Update(gomock.Any(), jobModel).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with the same error if find one fails on the second call", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils.FakeJobEntity()
		jobEntity.Status = newStatus
		jobModel := testutils2.FakeJob(0)
		parsers.UpdateJobModelFromEntities(jobModel, jobEntity)

		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, tenantID).Return(jobModel, nil)
		mockTransactionDA.EXPECT().Update(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Update(gomock.Any(), jobModel).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, tenantID).Return(nil, expectedErr)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})
}
