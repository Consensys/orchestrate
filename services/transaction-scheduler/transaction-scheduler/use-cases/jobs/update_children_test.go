// +build unit

package jobs

import (
	"context"
	"fmt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
)

func TestUpdateChildren_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tenants := []string{"tenantID"}
	ctx := context.Background()
	status := utils.StatusNeverMined

	mockDB := mocks.NewMockDB(ctrl)
	mockJobDA := mocks.NewMockJobAgent(ctrl)
	mockLogDA := mocks.NewMockLogAgent(ctrl)
	mockDB.EXPECT().Job().Return(mockJobDA).AnyTimes()
	mockDB.EXPECT().Log().Return(mockLogDA).AnyTimes()

	usecase := NewUpdateChildrenUseCase(mockDB)

	t.Run("should execute use case successfully if parentJobUUID is set", func(t *testing.T) {
		parentJobUUID := "parentJobUUID"
		jobUUID := "jobUUID"

		jobsToUpdate := []*models.Job{testutils.FakeJobModel(1), testutils.FakeJobModel(1)}
		jobsToUpdate[0].Logs[0].Status = utils.StatusPending
		jobsToUpdate[1].Logs[0].Status = utils.StatusPending

		mockJobDA.EXPECT().LockOneByUUID(ctx, parentJobUUID).Return(nil)
		mockJobDA.EXPECT().Search(ctx, &entities.JobFilters{ParentJobUUID: parentJobUUID}, tenants).Return(jobsToUpdate, nil)
		mockLogDA.EXPECT().Insert(ctx, &models.Log{
			JobID:   &jobsToUpdate[0].ID,
			Status:  status,
			Message: fmt.Sprintf("sibling (or parent) job %s was mined instead", jobUUID),
		}).
			Return(nil)
		mockLogDA.EXPECT().Insert(ctx, &models.Log{
			JobID:   &jobsToUpdate[1].ID,
			Status:  status,
			Message: fmt.Sprintf("sibling (or parent) job %s was mined instead", jobUUID),
		}).
			Return(nil)

		err := usecase.Execute(ctx, jobUUID, parentJobUUID, status, tenants)

		assert.NoError(t, err)
	})

	t.Run("should execute use case successfully if parentJobUUID is not set", func(t *testing.T) {
		parentJobUUID := ""
		jobUUID := "jobUUID"

		jobsToUpdate := []*models.Job{testutils.FakeJobModel(1), testutils.FakeJobModel(1)}
		jobsToUpdate[0].Logs[0].Status = utils.StatusPending
		jobsToUpdate[1].Logs[0].Status = utils.StatusPending

		mockJobDA.EXPECT().LockOneByUUID(ctx, jobUUID).Return(nil)
		mockJobDA.EXPECT().Search(ctx, &entities.JobFilters{ParentJobUUID: jobUUID}, tenants).Return(jobsToUpdate, nil)
		mockLogDA.EXPECT().Insert(ctx, &models.Log{
			JobID:   &jobsToUpdate[0].ID,
			Status:  status,
			Message: fmt.Sprintf("sibling (or parent) job %s was mined instead", jobUUID),
		}).
			Return(nil)
		mockLogDA.EXPECT().Insert(ctx, &models.Log{
			JobID:   &jobsToUpdate[1].ID,
			Status:  status,
			Message: fmt.Sprintf("sibling (or parent) job %s was mined instead", jobUUID),
		}).
			Return(nil)

		err := usecase.Execute(ctx, jobUUID, parentJobUUID, status, tenants)

		assert.NoError(t, err)
	})

	t.Run("should not update status of the jobUUID job", func(t *testing.T) {
		parentJobUUID := "parentJobUUID"
		jobUUID := "jobUUID"

		jobsToUpdate := []*models.Job{testutils.FakeJobModel(1), testutils.FakeJobModel(1)}
		jobsToUpdate[0].UUID = jobUUID
		jobsToUpdate[0].Logs[0].Status = utils.StatusPending
		jobsToUpdate[1].Logs[0].Status = utils.StatusPending

		mockJobDA.EXPECT().LockOneByUUID(ctx, parentJobUUID).Return(nil)
		mockJobDA.EXPECT().Search(ctx, &entities.JobFilters{ParentJobUUID: parentJobUUID}, tenants).Return(jobsToUpdate, nil)
		mockLogDA.EXPECT().Insert(ctx, &models.Log{
			JobID:   &jobsToUpdate[1].ID,
			Status:  status,
			Message: fmt.Sprintf("sibling (or parent) job %s was mined instead", jobUUID),
		}).
			Return(nil)

		err := usecase.Execute(ctx, jobUUID, parentJobUUID, status, tenants)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if LockOneByUUID fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockJobDA.EXPECT().LockOneByUUID(ctx, gomock.Any()).Return(expectedErr)

		err := usecase.Execute(ctx, "jobUUID", "parentJobUUID", status, tenants)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateChildrenComponent), err)
	})

	t.Run("should fail with same error if Search fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockJobDA.EXPECT().LockOneByUUID(ctx, gomock.Any()).Return(nil)
		mockJobDA.EXPECT().Search(ctx, gomock.Any(), tenants).Return(nil, expectedErr)

		err := usecase.Execute(ctx, "jobUUID", "parentJobUUID", status, tenants)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateChildrenComponent), err)
	})

	t.Run("should fail with same error if Insert fails", func(t *testing.T) {
		jobsToUpdate := []*models.Job{testutils.FakeJobModel(1)}
		jobsToUpdate[0].Logs[0].Status = utils.StatusPending

		expectedErr := fmt.Errorf("error")

		mockJobDA.EXPECT().LockOneByUUID(ctx, gomock.Any()).Return(nil)
		mockJobDA.EXPECT().Search(ctx, gomock.Any(), tenants).Return(jobsToUpdate, nil)
		mockLogDA.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)

		err := usecase.Execute(ctx, "jobUUID", "parentJobUUID", status, tenants)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateChildrenComponent), err)
	})
}
