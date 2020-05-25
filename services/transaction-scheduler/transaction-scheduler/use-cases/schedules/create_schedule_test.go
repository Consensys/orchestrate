// +build unit

package schedules

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators/mocks"
)

func TestCreateSchedule_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScheduleDA := mocks.NewMockScheduleAgent(ctrl)
	mockTxValidator := mocks2.NewMockTransactionValidator(ctrl)
	mockDB := mocks.NewMockDB(ctrl)
	tenantID := "tenantID"
	chainUUID := "ChainUUID"

	mockDB.EXPECT().Schedule().Return(mockScheduleDA).AnyTimes()

	usecase := NewCreateScheduleUseCase(mockTxValidator, mockDB)
	ctx := context.Background()

	t.Run("should execute use case successfully", func(t *testing.T) {
		scheduleEntity := testutils.FakeScheduleEntity(chainUUID)

		mockTxValidator.EXPECT().
			ValidateChainExists(ctx, scheduleEntity.ChainUUID).
			Return(nil)

		mockScheduleDA.EXPECT().
			Insert(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, schedule *models.Schedule) error {
				schedule.UUID = scheduleEntity.UUID
				schedule.ChainUUID = scheduleEntity.ChainUUID
				schedule.ID = 1
				return nil
			})

		scheduleResponse, err := usecase.Execute(ctx, scheduleEntity, tenantID)

		assert.Nil(t, err)
		assert.Equal(t, scheduleEntity.UUID, scheduleResponse.UUID)
		assert.Equal(t, scheduleEntity.ChainUUID, scheduleResponse.ChainUUID)
	})

	t.Run("should fail with same error if validator fails", func(t *testing.T) {
		scheduleEntity := testutils.FakeScheduleEntity(chainUUID)
		expectedErr := errors.DataError("error")
	
		mockTxValidator.EXPECT().ValidateChainExists(ctx, scheduleEntity.ChainUUID).Return(expectedErr)
	
		scheduleResponse, err := usecase.Execute(ctx, scheduleEntity, tenantID)
		assert.Nil(t, scheduleResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createScheduleComponent), err)
	})
	
	t.Run("should fail with same error if Insert fails", func(t *testing.T) {
		scheduleEntity := testutils.FakeScheduleEntity(chainUUID)
		expectedErr := errors.PostgresConnectionError("error")
	
		mockTxValidator.EXPECT().ValidateChainExists(ctx, scheduleEntity.ChainUUID).Return(nil)
		mockScheduleDA.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)
	
		scheduleResponse, err := usecase.Execute(ctx, scheduleEntity, tenantID)
		assert.Nil(t, scheduleResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createScheduleComponent), err)
	})
}
