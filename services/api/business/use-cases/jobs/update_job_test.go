// +build unit

package jobs

import (
	"context"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/metrics/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"

	testutils3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/mocks"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models/testutils"
)

func TestUpdateJob_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockDBTX := mocks.NewMockTx(ctrl)
	mockTransactionDA := mocks.NewMockTransactionAgent(ctrl)
	mockJobDA := mocks.NewMockJobAgent(ctrl)
	mockLogDA := mocks.NewMockLogAgent(ctrl)
	mockUpdateChilrenUC := mocks2.NewMockUpdateChildrenUseCase(ctrl)
	mockStartNextJobUC := mocks2.NewMockStartNextJobUseCase(ctrl)
	mockMetrics := mock.NewMockTransactionSchedulerMetrics(ctrl)

	jobsLatencyHistogram := mock2.NewMockHistogram(ctrl)
	jobsLatencyHistogram.EXPECT().With(gomock.Any()).AnyTimes().Return(jobsLatencyHistogram)
	jobsLatencyHistogram.EXPECT().Observe(gomock.Any()).AnyTimes()
	mockMetrics.EXPECT().JobsLatencyHistogram().AnyTimes().Return(jobsLatencyHistogram)

	minedLatencyHistogram := mock2.NewMockHistogram(ctrl)
	minedLatencyHistogram.EXPECT().With(gomock.Any()).AnyTimes().Return(minedLatencyHistogram)
	minedLatencyHistogram.EXPECT().Observe(gomock.Any()).AnyTimes()
	mockMetrics.EXPECT().MinedLatencyHistogram().AnyTimes().Return(minedLatencyHistogram)

	mockDB.EXPECT().Job().Return(mockJobDA).AnyTimes()
	mockDB.EXPECT().Begin().Return(mockDBTX, nil).AnyTimes()
	mockDBTX.EXPECT().Transaction().Return(mockTransactionDA).AnyTimes()
	mockDBTX.EXPECT().Job().Return(mockJobDA).AnyTimes()
	mockDBTX.EXPECT().Log().Return(mockLogDA).AnyTimes()
	mockDBTX.EXPECT().Commit().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Rollback().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Close().Return(nil).AnyTimes()
	mockUpdateChilrenUC.EXPECT().WithDBTransaction(mockDBTX).Return(mockUpdateChilrenUC).AnyTimes()

	usecase := NewUpdateJobUseCase(mockDB, mockUpdateChilrenUC, mockStartNextJobUC, mockMetrics)

	tenantID := "tenantID"
	nextStatus := utils.StatusStarted
	logMessage := "message"
	ctx := context.Background()

	t.Run("should execute use case successfully", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID

		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.UUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.UUID, []string{tenantID}).Return(jobModel, nil).Times(2)
		mockTransactionDA.EXPECT().Update(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, jobModelUpdate *models.Job) error {
			assert.Equal(t, jobModelUpdate.InternalData, jobEntity.InternalData)
			assert.Equal(t, jobModelUpdate.Labels, jobEntity.Labels)
			jobModel.ID = 1
			return nil
		})
		mockLogDA.EXPECT().Insert(gomock.Any(), &models.Log{
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

		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.UUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.UUID, []string{tenantID}).Return(jobModel, nil).Times(2)
		mockJobDA.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, jobModelUpdate *models.Job) error {
			assert.Equal(t, jobModelUpdate.InternalData, jobEntity.InternalData)
			assert.Equal(t, jobModelUpdate.Labels, jobEntity.Labels)
			jobModel.ID = 1
			return nil
		})
		mockLogDA.EXPECT().Insert(gomock.Any(), &models.Log{
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

		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.UUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.UUID, []string{tenantID}).Return(jobModel, nil).Times(2)
		mockTransactionDA.EXPECT().Update(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, jobModelUpdate *models.Job) error {
			assert.Equal(t, jobModelUpdate.InternalData, jobEntity.InternalData)
			assert.Equal(t, jobModelUpdate.Labels, jobEntity.Labels)
			jobModel.ID = 1
			return nil
		})

		_, err := usecase.Execute(ctx, jobEntity, "", "", []string{tenantID})

		assert.NoError(t, err)
	})

	t.Run("should execute use case successfully if status is PENDING", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		status := utils.StatusPending
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID
		jobModel.Logs[0].Status = utils.StatusStarted

		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.UUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.UUID, []string{tenantID}).Return(jobModel, nil).Times(2)
		mockJobDA.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, jobModelUpdate *models.Job) error {
			assert.Equal(t, jobModelUpdate.InternalData, jobEntity.InternalData)
			assert.Equal(t, jobModelUpdate.Labels, jobEntity.Labels)
			jobModel.ID = 1
			jobModel.Logs[0].Status = status
			return nil
		})
		mockLogDA.EXPECT().Insert(gomock.Any(), &models.Log{
			JobID:   &jobModel.ID,
			Status:  status,
			Message: logMessage,
		}).Return(nil)

		_, err := usecase.Execute(ctx, jobEntity, status, logMessage, []string{tenantID})
		assert.NoError(t, err)
	})

	t.Run("should execute use case successfully if status is MINED and update all the children jobs", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		status := utils.StatusMined
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID
		jobModel.Logs[0].Status = utils.StatusPending

		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.UUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.UUID, []string{tenantID}).Return(jobModel, nil).Times(2)
		mockJobDA.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, jobModelUpdate *models.Job) error {
			assert.Equal(t, jobModelUpdate.InternalData, jobEntity.InternalData)
			assert.Equal(t, jobModelUpdate.Labels, jobEntity.Labels)
			jobModel.ID = 1
			jobModel.Logs[0].Status = status
			return nil
		})
		mockLogDA.EXPECT().Insert(gomock.Any(), &models.Log{
			JobID:   &jobModel.ID,
			Status:  status,
			Message: logMessage,
		}).Return(nil)
		mockUpdateChilrenUC.EXPECT().
			Execute(gomock.Any(), jobModel.UUID, jobModel.InternalData.ParentJobUUID, utils.StatusNeverMined, []string{tenantID}).
			Return(nil)

		_, err := usecase.Execute(ctx, jobEntity, status, logMessage, []string{tenantID})
		assert.NoError(t, err)
	})

	t.Run("should fail with InvalidParameterError if status is MINED", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Logs[0].Status = utils.StatusMined
		jobModel.Schedule.TenantID = tenantID

		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.UUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.UUID, []string{tenantID}).Return(jobModel, nil)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, []string{tenantID})
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with the same error if update transaction fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils3.FakeJob()
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID

		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.UUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(jobModel, nil)
		mockTransactionDA.EXPECT().Update(gomock.Any(), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, []string{tenantID})
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with the same error if update job fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils3.FakeJob()
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID

		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.UUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(jobModel, nil)
		mockTransactionDA.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockJobDA.EXPECT().Update(gomock.Any(), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, []string{tenantID})
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with InvalidStateError if status is invalid for CREATED", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID
		jobModel.Logs[0].Status = utils.StatusCreated

		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.UUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(jobModel, nil)

		_, err := usecase.Execute(ctx, jobEntity, utils.StatusPending, logMessage, []string{tenantID})
		assert.True(t, errors.IsInvalidStateError(err))
	})

	t.Run("should fail with InvalidStateError if status is invalid for STARTED", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID
		jobModel.Logs[0].Status = utils.StatusStarted

		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.UUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(jobModel, nil)

		_, err := usecase.Execute(ctx, jobEntity, utils.StatusMined, logMessage, []string{tenantID})
		assert.True(t, errors.IsInvalidStateError(err))
	})

	t.Run("should fail with InvalidStateError if status is invalid for PENDING", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID
		jobModel.Logs[0].Status = utils.StatusPending

		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.UUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(jobModel, nil)

		_, err := usecase.Execute(ctx, jobEntity, utils.StatusStarted, logMessage, []string{tenantID})
		assert.True(t, errors.IsInvalidStateError(err))
	})

	t.Run("should fail with InvalidStateError if status is invalid for FAILED", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID
		jobModel.Logs[0].Status = utils.StatusCreated

		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.UUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(jobModel, nil)

		_, err := usecase.Execute(ctx, jobEntity, utils.StatusFailed, logMessage, []string{tenantID})
		assert.True(t, errors.IsInvalidStateError(err))
	})

	t.Run("should fail with the same error if insert log fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID

		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.UUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(jobModel, nil)
		mockJobDA.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, []string{tenantID})
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with the same error if find one fails on the second call", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = tenantID

		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobEntity.UUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(jobModel, nil)
		mockJobDA.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, []string{tenantID})
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should trigger next job start if nextStatus is STORED", func(t *testing.T) {
		jobModel := testutils2.FakeJobModel(0)
		nextJobModel := testutils2.FakeJobModel(0)
		jobModel.Logs[0].Status = utils.StatusStarted
		jobModel.Schedule.TenantID = tenantID
		jobModel.NextJobUUID = nextJobModel.UUID

		tenants := []string{tenantID}

		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobModel.UUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobModel.UUID, tenants).Return(jobModel, nil)
		mockJobDA.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockTransactionDA.EXPECT().Update(gomock.Any(), jobModel.Transaction).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(jobModel, nil)
		mockStartNextJobUC.EXPECT().Execute(gomock.Any(), jobModel.UUID, tenants).Return(nil)

		jobEntity := parsers.NewJobEntityFromModels(jobModel)
		_, err := usecase.Execute(ctx, jobEntity, utils.StatusStored, "", []string{tenantID})

		assert.NoError(t, err)
	})
}
