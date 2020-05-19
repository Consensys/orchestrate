// +build unit

package schedules

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/orm/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
	"testing"
)

func TestGetSchedule_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockORM := mocks2.NewMockORM(ctrl)
	mockDB := mocks.NewMockDB(ctrl)

	usecase := NewGetScheduleUseCase(mockDB, mockORM)
	tenantID := "tenantID"
	ctx := context.Background()

	t.Run("should execute use case successfully", func(t *testing.T) {
		schedule := testutils.FakeSchedule("")
		expectedResponse := utils.FormatScheduleResponse(schedule)

		mockORM.EXPECT().FetchScheduleByUUID(ctx, mockDB, schedule.UUID, tenantID).Return(schedule, nil)

		scheduleResponse, err := usecase.Execute(ctx, schedule.UUID, tenantID)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, scheduleResponse)
	})

	t.Run("should fail with same error if FindOne fails for schedules", func(t *testing.T) {
		uuid := "uuid"
		expectedErr := errors.NotFoundError("error")

		mockORM.EXPECT().FetchScheduleByUUID(ctx, mockDB, uuid, tenantID).Return(nil, expectedErr)

		scheduleResponse, err := usecase.Execute(ctx, uuid, tenantID)

		assert.Nil(t, scheduleResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createScheduleComponent), err)
	})
}
