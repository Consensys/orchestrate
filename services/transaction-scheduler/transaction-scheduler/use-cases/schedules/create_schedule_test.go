// +build unit

package schedules

import (
	"context"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/store/models"
)

func TestCreateSchedule_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScheduleDA := mocks.NewMockScheduleAgent(ctrl)
	mockDB := mocks.NewMockDB(ctrl)

	mockDB.EXPECT().Schedule().Return(mockScheduleDA).AnyTimes()

	usecase := NewCreateScheduleUseCase(mockDB)
	ctx := context.Background()

	t.Run("should execute use case successfully", func(t *testing.T) {
		scheduleEntity := testutils.FakeSchedule()

		mockScheduleDA.EXPECT().
			Insert(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, schedule *models.Schedule) error {
				schedule.UUID = scheduleEntity.UUID
				schedule.ID = 1
				return nil
			})

		scheduleResponse, err := usecase.Execute(ctx, scheduleEntity)

		assert.NoError(t, err)
		assert.Equal(t, scheduleEntity.UUID, scheduleResponse.UUID)
	})

	t.Run("should fail with same error if Insert fails", func(t *testing.T) {
		scheduleEntity := testutils.FakeSchedule()
		expectedErr := errors.PostgresConnectionError("error")

		mockScheduleDA.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)

		scheduleResponse, err := usecase.Execute(ctx, scheduleEntity)
		assert.Nil(t, scheduleResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createScheduleComponent), err)
	})
}
