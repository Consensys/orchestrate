// +build unit

package schedules

import (
	"context"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/services/api/store/mocks"
	"github.com/consensys/orchestrate/services/api/store/models"
)

func TestCreateSchedule_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScheduleDA := mocks.NewMockScheduleAgent(ctrl)
	mockDB := mocks.NewMockDB(ctrl)

	mockDB.EXPECT().Schedule().Return(mockScheduleDA).AnyTimes()

	userInfo := multitenancy.NewUserInfo("tenantOne", "username")
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

		scheduleResponse, err := usecase.Execute(ctx, scheduleEntity, userInfo)

		assert.NoError(t, err)
		assert.Equal(t, scheduleEntity.UUID, scheduleResponse.UUID)
	})

	t.Run("should fail with same error if Insert fails", func(t *testing.T) {
		scheduleEntity := testutils.FakeSchedule()
		expectedErr := errors.PostgresConnectionError("error")

		mockScheduleDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		scheduleResponse, err := usecase.Execute(ctx, scheduleEntity, userInfo)
		assert.Nil(t, scheduleResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createScheduleComponent), err)
	})
}
