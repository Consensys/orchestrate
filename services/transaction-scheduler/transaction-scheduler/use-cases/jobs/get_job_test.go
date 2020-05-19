// +build unit

package jobs

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
)

func TestGetJob_Execute(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockJobDA := mocks.NewMockJobAgent(ctrl)

	mockDB.EXPECT().Job().Return(mockJobDA).AnyTimes()
	usecase := NewGetJobUseCase(mockDB)
	
	tenantID := "tenantID"

	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJob(0)
		expectedResponse := utils.FormatJobResponse(job)
		
		mockJobDA.EXPECT().FindOneByUUID(ctx, job.UUID, tenantID).Return(job, nil)
		jobResponse, err := usecase.Execute(ctx, job.UUID, tenantID)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, jobResponse)
	})

	t.Run("should fail with same error if FindOneByUUID fails for job", func(t *testing.T) {
		uuid := "uuid"
		expectedErr := errors.NotFoundError("error")

		mockJobDA.EXPECT().FindOneByUUID(ctx, uuid, tenantID).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, uuid, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})
}
