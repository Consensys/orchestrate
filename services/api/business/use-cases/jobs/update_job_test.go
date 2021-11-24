// +build unit

package jobs

import (
	"context"
	"testing"

	"github.com/consensys/orchestrate/pkg/errors"
	mock2 "github.com/consensys/orchestrate/pkg/toolkit/app/metrics/mock"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	mocks2 "github.com/consensys/orchestrate/services/api/business/use-cases/mocks"
	"github.com/consensys/orchestrate/services/api/metrics/mock"
	"github.com/consensys/orchestrate/services/api/store/models"

	testutils3 "github.com/consensys/orchestrate/pkg/types/testutils"

	"github.com/consensys/orchestrate/services/api/store/mocks"
	testutils2 "github.com/consensys/orchestrate/services/api/store/models/testutils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
	mockDB.EXPECT().Transaction().Return(mockTransactionDA).AnyTimes()
	mockDBTX.EXPECT().Job().Return(mockJobDA).AnyTimes()
	mockDBTX.EXPECT().Log().Return(mockLogDA).AnyTimes()
	mockDBTX.EXPECT().Transaction().Return(mockTransactionDA).AnyTimes()
	mockDBTX.EXPECT().Commit().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Rollback().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Close().Return(nil).AnyTimes()
	mockUpdateChilrenUC.EXPECT().WithDBTransaction(mockDBTX).Return(mockUpdateChilrenUC).AnyTimes()

	userInfo := multitenancy.NewUserInfo("tenantOne", "username")
	usecase := NewUpdateJobUseCase(mockDB, mockUpdateChilrenUC, mockStartNextJobUC, mockMetrics)

	nextStatus := entities.StatusStarted
	logMessage := "message"
	ctx := context.Background()

	t.Run("should execute use case successfully", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = userInfo.TenantID

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.UUID, userInfo.AllowedTenants, userInfo.Username, true).
			Return(jobModel, nil)
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

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, userInfo)

		assert.NoError(t, err)
	})

	t.Run("should execute use case successfully if transaction is empty", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = userInfo.TenantID

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.UUID, userInfo.AllowedTenants, userInfo.Username, true).
			Return(jobModel, nil)
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

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, userInfo)

		assert.NoError(t, err)
	})

	t.Run("should execute use case successfully if status is empty", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = userInfo.TenantID

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.UUID, userInfo.AllowedTenants, userInfo.Username, true).
			Return(jobModel, nil)
		mockTransactionDA.EXPECT().Update(gomock.Any(), jobModel.Transaction).Return(nil)
		mockJobDA.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, jobModelUpdate *models.Job) error {
			assert.Equal(t, jobModelUpdate.InternalData, jobEntity.InternalData)
			assert.Equal(t, jobModelUpdate.Labels, jobEntity.Labels)
			jobModel.ID = 1
			return nil
		})

		_, err := usecase.Execute(ctx, jobEntity, "", "", userInfo)

		assert.NoError(t, err)
	})

	t.Run("should execute use case successfully if status is PENDING", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Status = entities.StatusStarted
		jobEntity.Transaction = nil
		status := entities.StatusPending
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = userInfo.TenantID
		jobModel.Logs = append(jobModel.Logs, &models.Log{Status: entities.StatusStarted})
		jobModel.Status = entities.StatusStarted

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.UUID, userInfo.AllowedTenants, userInfo.Username, true).
			Return(jobModel, nil)
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

		_, err := usecase.Execute(ctx, jobEntity, status, logMessage, userInfo)
		assert.NoError(t, err)
	})

	t.Run("should execute use case successfully if status is MINED", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobEntity.Status = entities.StatusPending
		status := entities.StatusMined
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = userInfo.TenantID
		jobModel.Status = entities.StatusPending
		jobModel.Logs = append(jobModel.Logs, &models.Log{Status: entities.StatusStarted})
		jobModel.Logs = append(jobModel.Logs, &models.Log{Status: entities.StatusPending})

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.UUID, userInfo.AllowedTenants, userInfo.Username, true).
			Return(jobModel, nil)
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

		_, err := usecase.Execute(ctx, jobEntity, status, logMessage, userInfo)
		assert.NoError(t, err)
	})
	
	t.Run("should execute use case successfully if status is MINED and update all the children jobs", func(t *testing.T) {
		jobParentEntity := testutils3.FakeJob()
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobEntity.Status = entities.StatusPending
		jobModel := testutils2.FakeJobModel(0)
		jobModel.UUID = jobEntity.UUID
		jobModel.Schedule.TenantID = userInfo.TenantID
		jobModel.Status = entities.StatusPending
		jobModel.InternalData.ParentJobUUID = jobParentEntity.UUID
		jobEntity.InternalData = jobModel.InternalData 

		nextStatus := entities.StatusMined

		mockJobDA.EXPECT().LockOneByUUID(gomock.Any(), jobParentEntity.UUID).Return(nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.UUID, userInfo.AllowedTenants, userInfo.Username, true).
			Return(jobModel, nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.UUID, userInfo.AllowedTenants, userInfo.Username, false).
			Return(jobModel, nil)
		mockJobDA.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, jobModelUpdate *models.Job) error {
			assert.Equal(t, jobModelUpdate.InternalData, jobEntity.InternalData)
			assert.Equal(t, jobModelUpdate.Labels, jobEntity.Labels)
			jobModel.ID = 1
			jobModel.Logs[0].Status = nextStatus
			return nil
		})
		mockLogDA.EXPECT().Insert(gomock.Any(), &models.Log{
			JobID:   &jobModel.ID,
			Status:  nextStatus,
			Message: logMessage,
		}).Return(nil)
		mockUpdateChilrenUC.EXPECT().
			Execute(gomock.Any(), jobModel.UUID, jobParentEntity.UUID, entities.StatusNeverMined, userInfo).
			Return(nil)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, userInfo)
		assert.NoError(t, err)
	})

	t.Run("should fail with InvalidParameterError if status is MINED", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Status = entities.StatusMined
		jobModel.Logs[0].Status = entities.StatusMined
		jobModel.Schedule.TenantID = userInfo.TenantID

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobEntity.UUID, userInfo.AllowedTenants, userInfo.Username, true).
			Return(jobModel, nil)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, userInfo)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with the same error if update transaction fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils3.FakeJob()
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = userInfo.TenantID

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), userInfo.AllowedTenants, userInfo.Username, true).
			Return(jobModel, nil)
		mockTransactionDA.EXPECT().Update(gomock.Any(), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, userInfo)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with the same error if update job fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils3.FakeJob()
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = userInfo.TenantID

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), userInfo.AllowedTenants, userInfo.Username, true).Return(jobModel, nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockTransactionDA.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockJobDA.EXPECT().Update(gomock.Any(), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, userInfo)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with InvalidStateError if status is invalid for CREATED", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = userInfo.TenantID
		jobModel.Logs[0].Status = entities.StatusCreated

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), userInfo.AllowedTenants, userInfo.Username, true).Return(jobModel, nil)

		_, err := usecase.Execute(ctx, jobEntity, entities.StatusPending, logMessage, userInfo)
		assert.True(t, errors.IsInvalidStateError(err))
	})

	t.Run("should fail with InvalidStateError if status is invalid for STARTED", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = userInfo.TenantID
		jobModel.Status = entities.StatusStarted

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), userInfo.AllowedTenants, userInfo.Username, true).Return(jobModel, nil)

		_, err := usecase.Execute(ctx, jobEntity, entities.StatusMined, logMessage, userInfo)
		assert.True(t, errors.IsInvalidStateError(err))
	})

	t.Run("should fail with InvalidStateError if status is invalid for PENDING", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = userInfo.TenantID
		jobModel.Status = entities.StatusPending

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), userInfo.AllowedTenants, userInfo.Username, true).Return(jobModel, nil)

		_, err := usecase.Execute(ctx, jobEntity, entities.StatusStarted, logMessage, userInfo)
		assert.True(t, errors.IsInvalidStateError(err))
	})

	t.Run("should fail with InvalidStateError if status is invalid for FAILED", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = userInfo.TenantID
		jobModel.Status = entities.StatusCreated
		jobModel.Logs[0].Status = entities.StatusCreated

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), userInfo.AllowedTenants, userInfo.Username, true).Return(jobModel, nil)
		_, err := usecase.Execute(ctx, jobEntity, entities.StatusFailed, logMessage, userInfo)
		assert.True(t, errors.IsInvalidStateError(err))
	})
	
	t.Run("should allow transition of job from RESENDING to RESENDING", func(t *testing.T) {
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = userInfo.TenantID
		jobModel.Status = entities.StatusResending
		jobModel.Logs[0].Status = entities.StatusResending

		status := entities.StatusResending
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), userInfo.AllowedTenants, userInfo.Username, true).Return(jobModel, nil)
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
		_, err := usecase.Execute(ctx, jobEntity, status, logMessage, userInfo)
		assert.NoError(t, err)
	})

	t.Run("should fail with the same error if insert log fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils3.FakeJob()
		jobEntity.Transaction = nil
		jobModel := testutils2.FakeJobModel(0)
		jobModel.Schedule.TenantID = userInfo.TenantID

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), gomock.Any(), userInfo.AllowedTenants, userInfo.Username, true).Return(jobModel, nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(ctx, jobEntity, nextStatus, logMessage, userInfo)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should trigger next job start if nextStatus is STORED", func(t *testing.T) {
		jobModel := testutils2.FakeJobModel(0)
		nextJobModel := testutils2.FakeJobModel(0)
		jobModel.Status = entities.StatusStarted
		jobModel.Logs[0].Status = entities.StatusStarted
		jobModel.Schedule.TenantID = userInfo.TenantID
		jobModel.NextJobUUID = nextJobModel.UUID

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobModel.UUID, userInfo.AllowedTenants, userInfo.Username, true).Return(jobModel, nil)
		mockJobDA.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockTransactionDA.EXPECT().Update(gomock.Any(), jobModel.Transaction).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockStartNextJobUC.EXPECT().Execute(gomock.Any(), jobModel.UUID, userInfo).Return(nil)

		jobEntity := parsers.NewJobEntityFromModels(jobModel)
		_, err := usecase.Execute(ctx, jobEntity, entities.StatusStored, "", userInfo)

		assert.NoError(t, err)
	})
}
