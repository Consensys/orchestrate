// +build unit

package schedules

import (
	"context"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/pkg/types/testutils"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/services/api/business/parsers"
	"github.com/ConsenSys/orchestrate/services/api/store/mocks"
	"github.com/ConsenSys/orchestrate/services/api/store/models"
)

func TestSearchSchedules_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockScheduleDA := mocks.NewMockScheduleAgent(ctrl)
	mockJobDA := mocks.NewMockJobAgent(ctrl)

	usecase := NewSearchSchedulesUseCase(mockDB)
	tenantID := "tenantID"
	ctx := context.Background()

	t.Run("should execute use case successfully", func(t *testing.T) {
		scheduleEntity := testutils.FakeSchedule()
		scheduleModel := parsers.NewScheduleModelFromEntities(scheduleEntity)
		expectedResponse := []*entities.Schedule{parsers.NewScheduleEntityFromModels(scheduleModel)}

		mockDB.EXPECT().Schedule().Return(mockScheduleDA).Times(1)
		mockDB.EXPECT().Job().Return(mockJobDA).Times(1)

		mockScheduleDA.EXPECT().
			FindAll(gomock.Any(), []string{tenantID}).
			Return([]*models.Schedule{scheduleModel}, nil)

		mockJobDA.EXPECT().
			FindOneByUUID(gomock.Any(), scheduleModel.Jobs[0].UUID, []string{tenantID}, false).
			Return(scheduleModel.Jobs[0], nil)

		schedulesResponse, err := usecase.Execute(ctx, []string{tenantID})

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, schedulesResponse)
	})

	t.Run("should fail with same error if FindAll fails for schedules", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		mockDB.EXPECT().Schedule().Return(mockScheduleDA).Times(1)

		mockScheduleDA.EXPECT().
			FindAll(gomock.Any(), []string{tenantID}).
			Return(nil, expectedErr)

		scheduleResponse, err := usecase.Execute(ctx, []string{tenantID})

		assert.Nil(t, scheduleResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createScheduleComponent), err)
	})

	t.Run("should fail with same error if FindOne fails for jobs", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		scheduleEntity := testutils.FakeSchedule()
		scheduleModel := parsers.NewScheduleModelFromEntities(scheduleEntity)

		mockDB.EXPECT().Schedule().Return(mockScheduleDA).Times(1)
		mockDB.EXPECT().Job().Return(mockJobDA).Times(1)

		mockScheduleDA.EXPECT().
			FindAll(gomock.Any(), []string{tenantID}).
			Return([]*models.Schedule{scheduleModel}, nil)

		mockJobDA.EXPECT().
			FindOneByUUID(gomock.Any(), scheduleModel.Jobs[0].UUID, []string{tenantID}, false).
			Return(scheduleModel.Jobs[0], expectedErr)

		scheduleResponse, err := usecase.Execute(ctx, []string{tenantID})

		assert.Nil(t, scheduleResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createScheduleComponent), err)
	})
}
