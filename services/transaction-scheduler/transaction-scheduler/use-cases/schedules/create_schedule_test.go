// +build unit

package schedules

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types/testutils"
	"testing"
	"time"
)

func TestCreateSchedule_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScheduleDA := mocks.NewMockScheduleAgent(ctrl)
	mockChainRegistryClient := mock.NewMockChainRegistryClient(ctrl)
	usecase := NewCreateScheduleUseCase(mockChainRegistryClient, mockScheduleDA)
	ctx := context.Background()

	t.Run("should execute use case successfully", func(t *testing.T) {
		timeNow := time.Now()
		scheduleRequest := testutils.FakeScheduleRequest()
		expectedResponse := &types.ScheduleResponse{
			UUID:      "testUUID",
			ChainID:   scheduleRequest.ChainID,
			CreatedAt: timeNow,
		}

		mockChainRegistryClient.EXPECT().GetChainByUUID(ctx, scheduleRequest.ChainID).Return(nil, nil)
		mockScheduleDA.EXPECT().Insert(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, schedule *models.Schedule) error {
			schedule.UUID = "testUUID"
			schedule.ID = 1
			schedule.CreatedAt = timeNow
			return nil
		})
		scheduleResponse, err := usecase.Execute(ctx, scheduleRequest, "tenantID")

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse.UUID, scheduleResponse.UUID)
		assert.Equal(t, expectedResponse.ChainID, scheduleResponse.ChainID)
		assert.Equal(t, expectedResponse.CreatedAt, scheduleResponse.CreatedAt)
	})

	t.Run("should fail with InvalidParameterError error if it fails to validate request", func(t *testing.T) {
		scheduleRequest := testutils.FakeScheduleRequest()
		scheduleRequest.ChainID = ""

		scheduleResponse, err := usecase.Execute(ctx, scheduleRequest, "tenantID")
		assert.True(t, errors.IsInvalidParameterError(err))
		assert.Nil(t, scheduleResponse)
	})

	t.Run("should fail with same error if chain registry fails", func(t *testing.T) {
		scheduleRequest := testutils.FakeScheduleRequest()
		expectedErr := errors.DataError("error")

		mockChainRegistryClient.EXPECT().GetChainByUUID(ctx, scheduleRequest.ChainID).Return(nil, expectedErr)

		scheduleResponse, err := usecase.Execute(ctx, scheduleRequest, "tenantID")
		assert.Nil(t, scheduleResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createScheduleComponent), err)
	})

	t.Run("should fail with same error if Insert fails", func(t *testing.T) {
		scheduleRequest := testutils.FakeScheduleRequest()
		expectedErr := errors.PostgresConnectionError("error")

		mockChainRegistryClient.EXPECT().GetChainByUUID(ctx, scheduleRequest.ChainID).Return(nil, nil)
		mockScheduleDA.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)

		scheduleResponse, err := usecase.Execute(ctx, scheduleRequest, "tenantID")
		assert.Nil(t, scheduleResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createScheduleComponent), err)
	})
}
