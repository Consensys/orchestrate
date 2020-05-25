// +build unit

package jobs

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

func TestSearchJobs_Execute(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockJobDA := mocks.NewMockJobAgent(ctrl)

	mockDB.EXPECT().Job().Return(mockJobDA).AnyTimes()
	usecase := NewSearchJobsUseCase(mockDB)
	
	tenantID := "tenantID"

	t.Run("should execute use case successfully", func(t *testing.T) {
		jobs := []*models.Job{testutils.FakeJob(0)}
		filters := make(map[string]string)
		expectedResponse := []*entities.Job{parsers.NewJobEntityFromModels(jobs[0])}
		
		mockJobDA.EXPECT().Search(ctx, filters, tenantID).Return(jobs, nil)
		jobResponse, err := usecase.Execute(ctx, filters, tenantID)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, jobResponse)
	})

	t.Run("should fail with same error if search fails for jobs", func(t *testing.T) {
		filters := make(map[string]string, 0)
		expectedErr := errors.NotFoundError("error")

		mockJobDA.EXPECT().Search(ctx, filters, tenantID).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, filters, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(searchJobsComponent), err)
	})
}
