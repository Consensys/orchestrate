// +build unit

package jobs

import (
	"context"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"testing"

	testutils3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
)

func TestUpdateJob_Execute(t *testing.T) {
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
	nextStatus := utils.StatusStarted
	logMessage := "message"
	ctx := context.Background()

	t.Run("should execute use case successfully", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID

		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, []string{tenantID}).Return(jobModel, nil).Times(2)
		mockTransactionDA.EXPECT().Update(ctx, jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, jobModelUpdate *models.Job) error {
			assert.Equal(t, jobModelUpdate.Annotations, jobEntity.Annotations)
			assert.Equal(t, jobModelUpdate.Labels, jobEntity.Labels)
			jobModel.ID = 1
			return nil
		})
		mockLogDA.EXPECT().Insert(ctx, &models.Log{
			JobID:   &jobModel.ID,
			Status:  nextStatus,
			Message: logMessage,
		}).Return(nil)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, []string{tenantID})

		assert.NoError(t, err)
	})

	t.Run("should execute use case successfully if transaction is empty", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID

		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, []string{tenantID}).Return(jobModel, nil).Times(2)
		mockJobDA.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, jobModelUpdate *models.Job) error {
			assert.Equal(t, jobModelUpdate.Annotations, jobEntity.Annotations)
			assert.Equal(t, jobModelUpdate.Labels, jobEntity.Labels)
			jobModel.ID = 1
			return nil
		})
		mockLogDA.EXPECT().Insert(ctx, &models.Log{
			JobID:   &jobModel.ID,
			Status:  nextStatus,
			Message: logMessage,
		}).Return(nil)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, []string{tenantID})

		assert.NoError(t, err)
	})

	t.Run("should execute use case successfully if status is empty", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID

		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, []string{tenantID}).Return(jobModel, nil).Times(2)
		mockTransactionDA.EXPECT().Update(ctx, jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, jobModelUpdate *models.Job) error {
			assert.Equal(t, jobModelUpdate.Annotations, jobEntity.Annotations)
			assert.Equal(t, jobModelUpdate.Labels, jobEntity.Labels)
			jobModel.ID = 1
			return nil
		})

		_, err := usecase.Execute(ctx, jobEntity, "", "", []string{tenantID})

		assert.NoError(t, err)
	})

	t.Run("should fail with InvalidParameterError if status is MINED", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Logs[0].Status = utils.StatusMined
		jobModel.Schedule.TenantID = tenantID

		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, []string{tenantID}).Return(jobModel, nil)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, []string{tenantID})
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with the same error if update transaction fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils3.FakeJob()
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID

		mockJobDA.EXPECT().FindOneByUUID(ctx, gomock.Any(), gomock.Any()).Return(jobModel, nil)
		mockTransactionDA.EXPECT().Update(ctx, gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, []string{tenantID})
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with the same error if update job fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils3.FakeJob()
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID

		mockJobDA.EXPECT().FindOneByUUID(ctx, gomock.Any(), gomock.Any()).Return(jobModel, nil)
		mockTransactionDA.EXPECT().Update(ctx, gomock.Any()).Return(nil)
		mockJobDA.EXPECT().Update(ctx, gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, []string{tenantID})
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with InvalidStateError if status is invalid for CREATED", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID
		jobModel.Logs[0].Status = utils.StatusCreated

		mockJobDA.EXPECT().FindOneByUUID(ctx, gomock.Any(), gomock.Any()).Return(jobModel, nil)
		mockJobDA.EXPECT().Update(ctx, gomock.Any()).Return(nil)

		_, err := usecase.Execute(ctx, jobEntity, utils.StatusPending, logMessage, []string{tenantID})
		assert.Equal(t, errors.InvalidStateError("invalid status update for the current job state").ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with InvalidStateError if status is invalid for STARTED", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID
		jobModel.Logs[0].Status = utils.StatusStarted

		mockJobDA.EXPECT().FindOneByUUID(ctx, gomock.Any(), gomock.Any()).Return(jobModel, nil)
		mockJobDA.EXPECT().Update(ctx, gomock.Any()).Return(nil)

		_, err := usecase.Execute(ctx, jobEntity, utils.StatusMined, logMessage, []string{tenantID})
		assert.Equal(t, errors.InvalidStateError("invalid status update for the current job state").ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with InvalidStateError if status is invalid for PENDING", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID
		jobModel.Logs[0].Status = utils.StatusPending

		mockJobDA.EXPECT().FindOneByUUID(ctx, gomock.Any(), gomock.Any()).Return(jobModel, nil)
		mockJobDA.EXPECT().Update(ctx, gomock.Any()).Return(nil)

		_, err := usecase.Execute(ctx, jobEntity, utils.StatusStarted, logMessage, []string{tenantID})
		assert.Equal(t, errors.InvalidStateError("invalid status update for the current job state").ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with InvalidStateError if status is invalid for FAILED", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID
		jobModel.Logs[0].Status = utils.StatusCreated

		mockJobDA.EXPECT().FindOneByUUID(ctx, gomock.Any(), gomock.Any()).Return(jobModel, nil)
		mockJobDA.EXPECT().Update(ctx, gomock.Any()).Return(nil)

		_, err := usecase.Execute(ctx, jobEntity, utils.StatusFailed, logMessage, []string{tenantID})
		assert.Equal(t, errors.InvalidStateError("invalid status update for the current job state").ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with the same error if insert log fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID

		mockJobDA.EXPECT().FindOneByUUID(ctx, gomock.Any(), gomock.Any()).Return(jobModel, nil)
		mockJobDA.EXPECT().Update(ctx, gomock.Any()).Return(nil)
		mockLogDA.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, []string{tenantID})
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with the same error if find one fails on the second call", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID

		mockJobDA.EXPECT().FindOneByUUID(ctx, gomock.Any(), gomock.Any()).Return(jobModel, nil)
		mockJobDA.EXPECT().Update(ctx, gomock.Any()).Return(nil)
		mockLogDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(ctx, gomock.Any(), gomock.Any()).Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, []string{tenantID})
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})
}
