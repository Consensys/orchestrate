// +build unit

package jobs

import (
	"context"
	"testing"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	"github.com/consensys/orchestrate/services/api/store/mocks"
	"github.com/consensys/orchestrate/services/api/store/models/testutils"
)

func TestGetJob_Execute(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockJobDA := mocks.NewMockJobAgent(ctrl)

	mockDB.EXPECT().Job().Return(mockJobDA).AnyTimes()
	userInfo := multitenancy.NewUserInfo("tenantOne", "username")
	usecase := NewGetJobUseCase(mockDB)

	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJobModel(0)
		expectedResponse := parsers.NewJobEntityFromModels(job)

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), job.UUID, userInfo.AllowedTenants, userInfo.Username, true).Return(job, nil)
		jobResponse, err := usecase.Execute(ctx, job.UUID, userInfo)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, jobResponse)
	})

	t.Run("should fail with same error if FindOneByUUID fails for job", func(t *testing.T) {
		uuid := "uuid"
		expectedErr := errors.NotFoundError("error")

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), uuid, userInfo.AllowedTenants, userInfo.Username, true).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, uuid, userInfo)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createJobComponent), err)
	})
}
