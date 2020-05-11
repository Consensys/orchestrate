// +build unit

package schedules

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/interfaces/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types/testutils"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators/mocks"
	"testing"
	"time"
)

func TestCreateSchedule_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScheduleDA := mocks.NewMockScheduleAgent(ctrl)
	mockTxValidator := mocks2.NewMockTransactionValidator(ctrl)
	mockDB := mocks.NewMockDB(ctrl)
	tenantID := "tenantID"

	mockDB.EXPECT().Schedule().Return(mockScheduleDA).AnyTimes()

	usecase := NewCreateScheduleUseCase(mockTxValidator, mockDB)
	ctx := context.Background()

	t.Run("should execute use case successfully", func(t *testing.T) {
		timeNow := time.Now()
		scheduleRequest := testutils.FakeScheduleRequest()
		expectedResponse := &types.ScheduleResponse{
			UUID:      "testUUID",
			ChainUUID: scheduleRequest.ChainUUID,
			CreatedAt: timeNow,
		}

		mockTxValidator.EXPECT().ValidateChainExists(ctx, scheduleRequest.ChainUUID).Return(nil)
		mockScheduleDA.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, schedule *models.Schedule) error {
			schedule.UUID = "testUUID"
			schedule.ID = 1
			schedule.CreatedAt = timeNow
			return nil
		})
		scheduleResponse, err := usecase.Execute(ctx, scheduleRequest, tenantID)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse.UUID, scheduleResponse.UUID)
		assert.Equal(t, expectedResponse.ChainUUID, scheduleResponse.ChainUUID)
		assert.Equal(t, expectedResponse.CreatedAt, scheduleResponse.CreatedAt)
	})

	t.Run("should fail with InvalidParameterError error if it fails to validate request", func(t *testing.T) {
		scheduleRequest := testutils.FakeScheduleRequest()
		scheduleRequest.ChainUUID = ""

		scheduleResponse, err := usecase.Execute(ctx, scheduleRequest, "tenantID")
		assert.True(t, errors.IsInvalidParameterError(err))
		assert.Nil(t, scheduleResponse)
	})

	t.Run("should fail with same error if validator fails", func(t *testing.T) {
		scheduleRequest := testutils.FakeScheduleRequest()
		expectedErr := errors.DataError("error")

		mockTxValidator.EXPECT().ValidateChainExists(ctx, scheduleRequest.ChainUUID).Return(expectedErr)

		scheduleResponse, err := usecase.Execute(ctx, scheduleRequest, tenantID)
		assert.Nil(t, scheduleResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createScheduleComponent), err)
	})

	t.Run("should fail with same error if Insert fails", func(t *testing.T) {
		scheduleRequest := testutils.FakeScheduleRequest()
		expectedErr := errors.PostgresConnectionError("error")

		mockTxValidator.EXPECT().ValidateChainExists(ctx, scheduleRequest.ChainUUID).Return(nil)
		mockScheduleDA.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)

		scheduleResponse, err := usecase.Execute(ctx, scheduleRequest, tenantID)
		assert.Nil(t, scheduleResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createScheduleComponent), err)
	})
}
