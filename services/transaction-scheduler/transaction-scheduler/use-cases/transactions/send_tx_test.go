// +build unit

package transactions

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/interfaces/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs/mocks"
	mocks3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators/mocks"
	"testing"
	"time"
)

func TestSendTx_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := mocks3.NewMockTransactionValidator(ctrl)
	mockTxRequestDA := mocks2.NewMockTransactionRequestAgent(ctrl)
	mockScheduleDA := mocks2.NewMockScheduleAgent(ctrl)
	mockJobDA := mocks2.NewMockJobAgent(ctrl)
	mockLogDA := mocks2.NewMockLogAgent(ctrl)
	mockTx := mocks2.NewMockTx(ctrl)
	mockDB := mocks2.NewMockDB(ctrl)

	mockTx.EXPECT().Schedule().Return(mockScheduleDA).AnyTimes()
	mockTx.EXPECT().Job().Return(mockJobDA).AnyTimes()
	mockTx.EXPECT().Log().Return(mockLogDA).AnyTimes()
	mockTx.EXPECT().TransactionRequest().Return(mockTxRequestDA).AnyTimes()

	mockStartJob := mocks.NewMockStartJobUseCase(ctrl)
	tenantID := "tenantID"

	usecase := NewSendTxUseCase(mockValidator, mockDB, mockStartJob)
	ctx := context.Background()

	t.Run("should execute use case successfully", func(t *testing.T) {
		timeNow := time.Now()
		txRequest := testutils.FakeTransactionRequest()

		mockValidator.EXPECT().ValidateRequestHash(ctx, txRequest.Params, txRequest.IdempotencyKey).Return("requestHash", nil)
		mockValidator.EXPECT().ValidateChainExists(ctx, gomock.Any()).Return(nil)
		mockDB.EXPECT().Begin().Return(mockTx, nil)
		mockScheduleDA.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, schedule *models.Schedule) error {
			schedule.ID = 1
			schedule.UUID = "scheduleUUID"
			return nil
		})
		mockJobDA.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, job *models.Job) error {
			job.ID = 1
			job.UUID = "jobUUID"
			return nil
		})
		mockLogDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		mockTxRequestDA.EXPECT().SelectOrInsert(ctx, gomock.Any()).Return(nil)
		mockTx.EXPECT().Commit().Return(nil)
		mockStartJob.EXPECT().Execute(ctx, "jobUUID", tenantID).Return(nil)

		txResponse, err := usecase.Execute(ctx, txRequest, tenantID)

		expectedResponse := &types.TransactionResponse{
			IdempotencyKey: txRequest.IdempotencyKey,
			Schedule: &types.ScheduleResponse{
				UUID:      "scheduleUUID",
				ChainUUID: txRequest.ChainUUID,
				CreatedAt: timeNow,
			},
			CreatedAt: timeNow,
		}
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse.IdempotencyKey, txResponse.IdempotencyKey)
		assert.Equal(t, expectedResponse.Schedule.UUID, txResponse.Schedule.UUID)
		assert.Equal(t, expectedResponse.Schedule.ChainUUID, txResponse.Schedule.ChainUUID)
		assert.Equal(t, expectedResponse.Schedule.CreatedAt, timeNow)
		assert.Equal(t, expectedResponse.CreatedAt, timeNow)
	})

	t.Run("should fail with same error if validator fails to validate request hash", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.InvalidParameterError("error")

		mockValidator.EXPECT().ValidateRequestHash(ctx, gomock.Any(), gomock.Any()).Return("", expectedErr)

		txResponse, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	t.Run("should fail with same error if validator fails to validate chain", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.InvalidParameterError("error")

		mockValidator.EXPECT().ValidateRequestHash(ctx, gomock.Any(), gomock.Any()).Return("requestHash", nil)
		mockValidator.EXPECT().ValidateChainExists(ctx, gomock.Any()).Return(expectedErr)

		txResponse, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	t.Run("should fail with same error if Begin fails", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.PostgresConnectionError("error")

		mockValidator.EXPECT().ValidateRequestHash(ctx, txRequest.Params, txRequest.IdempotencyKey).Return("requestHash", nil)
		mockValidator.EXPECT().ValidateChainExists(ctx, gomock.Any()).Return(nil)
		mockDB.EXPECT().Begin().Return(nil, expectedErr)

		txResponse, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	t.Run("should fail with same error if Insert fails for schedules", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.PostgresConnectionError("error")

		mockValidator.EXPECT().ValidateRequestHash(ctx, txRequest.Params, txRequest.IdempotencyKey).Return("requestHash", nil)
		mockValidator.EXPECT().ValidateChainExists(ctx, gomock.Any()).Return(nil)
		mockDB.EXPECT().Begin().Return(mockTx, nil)
		mockScheduleDA.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)

		txResponse, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	t.Run("should fail with same error if Insert fails for jobs", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.PostgresConnectionError("error")

		mockValidator.EXPECT().ValidateRequestHash(ctx, txRequest.Params, txRequest.IdempotencyKey).Return("requestHash", nil)
		mockValidator.EXPECT().ValidateChainExists(ctx, gomock.Any()).Return(nil)
		mockDB.EXPECT().Begin().Return(mockTx, nil)
		mockScheduleDA.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, schedule *models.Schedule) error {
			schedule.ID = 1
			schedule.UUID = "scheduleUUID"
			return nil
		})
		mockJobDA.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)

		txResponse, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	t.Run("should fail with same error if Insert fails for logs", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.PostgresConnectionError("error")

		mockValidator.EXPECT().ValidateRequestHash(ctx, txRequest.Params, txRequest.IdempotencyKey).Return("requestHash", nil)
		mockValidator.EXPECT().ValidateChainExists(ctx, gomock.Any()).Return(nil)
		mockDB.EXPECT().Begin().Return(mockTx, nil)
		mockScheduleDA.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, schedule *models.Schedule) error {
			schedule.ID = 1
			schedule.UUID = "scheduleUUID"
			return nil
		})
		mockJobDA.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, job *models.Job) error {
			job.ID = 1
			job.UUID = "jobUUID"
			return nil
		})
		mockLogDA.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)

		txResponse, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	t.Run("should fail with same error if SelectOrInsert fails for tx request", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.PostgresConnectionError("error")

		mockValidator.EXPECT().ValidateRequestHash(ctx, txRequest.Params, txRequest.IdempotencyKey).Return("requestHash", nil)
		mockValidator.EXPECT().ValidateChainExists(ctx, gomock.Any()).Return(nil)
		mockDB.EXPECT().Begin().Return(mockTx, nil)
		mockScheduleDA.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, schedule *models.Schedule) error {
			schedule.ID = 1
			schedule.UUID = "scheduleUUID"
			return nil
		})
		mockJobDA.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, job *models.Job) error {
			job.ID = 1
			job.UUID = "jobUUID"
			return nil
		})
		mockLogDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		mockTxRequestDA.EXPECT().SelectOrInsert(ctx, gomock.Any()).Return(expectedErr)

		txResponse, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	t.Run("should fail with same error if Commit fails", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.KafkaConnectionError("error")

		mockValidator.EXPECT().ValidateRequestHash(ctx, txRequest.Params, txRequest.IdempotencyKey).Return("requestHash", nil)
		mockValidator.EXPECT().ValidateChainExists(ctx, gomock.Any()).Return(nil)
		mockDB.EXPECT().Begin().Return(mockTx, nil)
		mockScheduleDA.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, schedule *models.Schedule) error {
			schedule.ID = 1
			schedule.UUID = "scheduleUUID"
			return nil
		})
		mockJobDA.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, job *models.Job) error {
			job.ID = 1
			job.UUID = "jobUUID"
			return nil
		})
		mockLogDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		mockTxRequestDA.EXPECT().SelectOrInsert(ctx, gomock.Any()).Return(nil)
		mockTx.EXPECT().Commit().Return(expectedErr)

		txResponse, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	t.Run("should fail with same error if start job fails", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.KafkaConnectionError("error")

		mockValidator.EXPECT().ValidateRequestHash(ctx, txRequest.Params, txRequest.IdempotencyKey).Return("requestHash", nil)
		mockValidator.EXPECT().ValidateChainExists(ctx, gomock.Any()).Return(nil)
		mockDB.EXPECT().Begin().Return(mockTx, nil)
		mockScheduleDA.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, schedule *models.Schedule) error {
			schedule.ID = 1
			schedule.UUID = "scheduleUUID"
			return nil
		})
		mockJobDA.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, job *models.Job) error {
			job.ID = 1
			job.UUID = "jobUUID"
			return nil
		})
		mockLogDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		mockTxRequestDA.EXPECT().SelectOrInsert(ctx, gomock.Any()).Return(nil)
		mockTx.EXPECT().Commit().Return(nil)
		mockStartJob.EXPECT().Execute(ctx, "jobUUID", tenantID).Return(expectedErr)

		txResponse, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})
}
