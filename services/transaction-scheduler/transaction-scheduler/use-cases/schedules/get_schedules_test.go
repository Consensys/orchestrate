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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/orm/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
)

func TestGetSchedules_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockORM := mocks2.NewMockORM(ctrl)
	mockDB := mocks.NewMockDB(ctrl)

	usecase := NewGetSchedulesUseCase(mockDB, mockORM)
	tenantID := "tenantID"
	ctx := context.Background()

	t.Run("should execute use case successfully", func(t *testing.T) {
		schedule := testutils.FakeSchedule("")
		expectedResponse := []*types.ScheduleResponse{utils.FormatScheduleResponse(schedule)}

		mockORM.EXPECT().FetchAllSchedules(ctx, mockDB, tenantID).Return([]*models.Schedule{schedule}, nil)

		scheduleResponse, err := usecase.Execute(ctx, tenantID)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, scheduleResponse)
	})

	t.Run("should fail with same error if FindAll fails for schedules", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		mockORM.EXPECT().FetchAllSchedules(ctx, mockDB, tenantID).Return([]*models.Schedule{}, expectedErr)

		scheduleResponse, err := usecase.Execute(ctx, tenantID)

		assert.Nil(t, scheduleResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createScheduleComponent), err)
	})
}
