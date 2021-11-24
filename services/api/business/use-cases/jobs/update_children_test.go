// +build unit

package jobs

import (
	"context"
	"fmt"
	"testing"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/store/models"
	"github.com/consensys/orchestrate/services/api/store/models/testutils"

	"github.com/consensys/orchestrate/services/api/store/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateChildren_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	status := entities.StatusNeverMined

	mockDB := mocks.NewMockDB(ctrl)
	mockJobDA := mocks.NewMockJobAgent(ctrl)
	mockLogDA := mocks.NewMockLogAgent(ctrl)
	mockDB.EXPECT().Job().Return(mockJobDA).AnyTimes()
	mockDB.EXPECT().Log().Return(mockLogDA).AnyTimes()

	userInfo := multitenancy.NewUserInfo("tenantOne", "username")
	usecase := NewUpdateChildrenUseCase(mockDB)

	t.Run("should execute use case successfully", func(t *testing.T) {
		parentJobUUID := "parentJobUUID"
		jobUUID := "jobUUID"

		jobsToUpdate := []*models.Job{testutils.FakeJobModel(1), testutils.FakeJobModel(1)}
		jobsToUpdate[0].Status = entities.StatusPending
		jobsToUpdate[1].Status = entities.StatusPending

		mockJobDA.EXPECT().Search(gomock.Any(), &entities.JobFilters{ParentJobUUID: parentJobUUID, Status: entities.StatusPending},
			userInfo.AllowedTenants, userInfo.Username).Return(jobsToUpdate, nil)
		mockJobDA.EXPECT().Update(gomock.Any(), jobsToUpdate[0]).Return(nil)
		mockJobDA.EXPECT().Update(gomock.Any(), jobsToUpdate[1]).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), &models.Log{
			JobID:   &jobsToUpdate[0].ID,
			Status:  status,
			Message: fmt.Sprintf("sibling (or parent) job %s was mined instead", jobUUID),
		}).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), &models.Log{
			JobID:   &jobsToUpdate[1].ID,
			Status:  status,
			Message: fmt.Sprintf("sibling (or parent) job %s was mined instead", jobUUID),
		}).
			Return(nil)

		err := usecase.Execute(ctx, jobUUID, parentJobUUID, status, userInfo)

		assert.NoError(t, err)
	})

	t.Run("should not update status of the jobUUID job", func(t *testing.T) {
		parentJobUUID := "parentJobUUID"
		jobUUID := "jobUUID"

		jobsToUpdate := []*models.Job{testutils.FakeJobModel(1), testutils.FakeJobModel(1)}
		jobsToUpdate[0].UUID = jobUUID
		jobsToUpdate[0].Status = entities.StatusPending
		jobsToUpdate[1].Status = entities.StatusPending

		mockJobDA.EXPECT().Search(gomock.Any(), &entities.JobFilters{ParentJobUUID: parentJobUUID, Status: entities.StatusPending},
			userInfo.AllowedTenants, userInfo.Username).Return(jobsToUpdate, nil)
		mockJobDA.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), &models.Log{
			JobID:   &jobsToUpdate[1].ID,
			Status:  status,
			Message: fmt.Sprintf("sibling (or parent) job %s was mined instead", jobUUID),
		}).Return(nil)

		err := usecase.Execute(ctx, jobUUID, parentJobUUID, status, userInfo)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if Search fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockJobDA.EXPECT().Search(gomock.Any(), gomock.Any(), userInfo.AllowedTenants, userInfo.Username).Return(nil, expectedErr)

		err := usecase.Execute(ctx, "jobUUID", "parentJobUUID", status, userInfo)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateChildrenComponent), err)
	})

	t.Run("should fail with same error if Insert fails", func(t *testing.T) {
		jobsToUpdate := []*models.Job{testutils.FakeJobModel(1)}
		jobsToUpdate[0].Status = entities.StatusPending

		expectedErr := fmt.Errorf("error")

		mockJobDA.EXPECT().Search(gomock.Any(), gomock.Any(), userInfo.AllowedTenants, userInfo.Username).Return(jobsToUpdate, nil)
		mockJobDA.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		err := usecase.Execute(ctx, "jobUUID", "parentJobUUID", status, userInfo)

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateChildrenComponent), err)
	})
}
