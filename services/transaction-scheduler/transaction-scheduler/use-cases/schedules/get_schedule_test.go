// +build unit

package schedules

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/interfaces/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
	"testing"
)

func TestGetSchedule_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScheduleDA := mocks.NewMockScheduleAgent(ctrl)
	mockJobDA := mocks.NewMockJobAgent(ctrl)
	mockDB := mocks.NewMockDB(ctrl)

	mockDB.EXPECT().Schedule().Return(mockScheduleDA).AnyTimes()
	mockDB.EXPECT().Job().Return(mockJobDA).AnyTimes()

	usecase := NewGetScheduleUseCase(mockDB)
	tenantID := "tenantID"
	ctx := context.Background()

	t.Run("should execute use case successfully", func(t *testing.T) {
		schedule := testutils.FakeSchedule()
		expectedResponse := utils.FormatScheduleResponse(schedule)

		mockScheduleDA.EXPECT().FindOneByUUID(ctx, schedule.UUID, tenantID).Return(schedule, nil)
		mockJobDA.EXPECT().FindOneByUUID(ctx, schedule.Jobs[0].UUID, tenantID).Return(schedule.Jobs[0], nil) // Necessary because of data agent not fetching in cascade

		scheduleResponse, err := usecase.Execute(ctx, schedule.UUID, tenantID)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, scheduleResponse)
	})

	t.Run("should fail with same error if FindOne fails for schedules", func(t *testing.T) {
		uuid := "uuid"
		expectedErr := errors.NotFoundError("error")

		mockScheduleDA.EXPECT().FindOneByUUID(ctx, uuid, tenantID).Return(nil, expectedErr)

		scheduleResponse, err := usecase.Execute(ctx, uuid, tenantID)

		assert.Nil(t, scheduleResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createScheduleComponent), err)
	})

	t.Run("should fail with same error if FindOne fails for jobs", func(t *testing.T) {
		schedule := testutils.FakeSchedule()
		expectedErr := errors.NotFoundError("error")

		mockScheduleDA.EXPECT().FindOneByUUID(ctx, schedule.UUID, tenantID).Return(schedule, nil)
		mockJobDA.EXPECT().FindOneByUUID(ctx, gomock.Any(), tenantID).Return(nil, expectedErr) // Necessary because of data agent not fetching in cascade

		scheduleResponse, err := usecase.Execute(ctx, schedule.UUID, tenantID)

		assert.Nil(t, scheduleResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createScheduleComponent), err)
	})
}
