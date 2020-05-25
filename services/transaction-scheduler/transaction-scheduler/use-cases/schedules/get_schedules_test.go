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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"
)

func TestGetSchedules_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockScheduleDA := mocks.NewMockScheduleAgent(ctrl)
	mockJobDA := mocks.NewMockJobAgent(ctrl)

	usecase := NewGetSchedulesUseCase(mockDB)
	tenantID := "tenantID"
	chainUUID := "ChainUUID"
	ctx := context.Background()

	t.Run("should execute use case successfully", func(t *testing.T) {
		scheduleEntity := testutils.FakeScheduleEntity(chainUUID)
		scheduleModel := parsers.NewScheduleModelFromEntities(scheduleEntity, tenantID)
		expectedResponse := []*entities.Schedule{parsers.NewScheduleEntityFromModels(scheduleModel)}
		
		mockDB.EXPECT().Schedule().Return(mockScheduleDA).Times(1)
		mockDB.EXPECT().Job().Return(mockJobDA).Times(1)
		
		mockScheduleDA.EXPECT().
			FindAll(gomock.Any(), tenantID).
			Return([]*models.Schedule{scheduleModel}, nil)

		mockJobDA.EXPECT().
			FindOneByUUID(gomock.Any(), scheduleModel.Jobs[0].UUID, tenantID).
			Return(scheduleModel.Jobs[0], nil)

		schedulesResponse, err := usecase.Execute(ctx, tenantID)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, schedulesResponse)
	})

	t.Run("should fail with same error if FindAll fails for schedules", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		mockDB.EXPECT().Schedule().Return(mockScheduleDA).Times(1)
		
		mockScheduleDA.EXPECT().
			FindAll(gomock.Any(), tenantID).
			Return(nil, expectedErr)

		scheduleResponse, err := usecase.Execute(ctx, tenantID)

		assert.Nil(t, scheduleResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createScheduleComponent), err)
	})

	t.Run("should fail with same error if FindOne fails for jobs", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		scheduleEntity := testutils.FakeScheduleEntity(chainUUID)
		scheduleModel := parsers.NewScheduleModelFromEntities(scheduleEntity, tenantID)

		mockDB.EXPECT().Schedule().Return(mockScheduleDA).Times(1)
		mockDB.EXPECT().Job().Return(mockJobDA).Times(1)

		mockScheduleDA.EXPECT().
			FindAll(gomock.Any(), tenantID).
			Return([]*models.Schedule{scheduleModel}, nil)
		
		mockJobDA.EXPECT().
			FindOneByUUID(gomock.Any(), scheduleModel.Jobs[0].UUID, tenantID).
			Return(scheduleModel.Jobs[0], expectedErr)

		scheduleResponse, err := usecase.Execute(ctx, tenantID)

		assert.Nil(t, scheduleResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createScheduleComponent), err)
	})
}
