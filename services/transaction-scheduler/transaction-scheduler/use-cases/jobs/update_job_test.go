// +build unit

package jobs

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"
)

func TestUpdateJob_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockDBTX := mocks.NewMockTx(ctrl)
	mockTransactionDA := mocks.NewMockTransactionAgent(ctrl)
	mockJobDA := mocks.NewMockJobAgent(ctrl)
	mockLogDA := mocks.NewMockLogAgent(ctrl)

	usecase := NewUpdateJobUseCase(mockDB)

	tenantID := "tenantID"

	t.Run("should execute use case successfully", func(t *testing.T) {
		jobEntity := testutils.FakeJobEntity()
		jobModel := testutils2.FakeJob(0)
		parsers.UpdateJobModelFromEntities(jobModel, jobEntity)

		mockDB.EXPECT().Job().Return(mockJobDA).Times(1)
		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, tenantID).Return(jobModel, nil)

		mockDB.EXPECT().Begin().Return(mockDBTX, nil).Times(1)
		mockDBTX.EXPECT().Transaction().Return(mockTransactionDA).Times(1)
		mockDBTX.EXPECT().Job().Return(mockJobDA).Times(1)
		mockDBTX.EXPECT().Log().Return(mockLogDA).Times(1)
		mockDBTX.EXPECT().Commit().Return(nil).Times(1)
		mockDBTX.EXPECT().Close().Return(nil).Times(1)

		mockTransactionDA.EXPECT().Update(gomock.Any(), jobModel.Transaction).Return(nil).Times(1)
		mockJobDA.EXPECT().Update(gomock.Any(), jobModel).Return(nil).Times(1)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)

		assert.Nil(t, err)
	})

	t.Run("should fail with the same error if find one fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils.FakeJobEntity()
		jobModel := testutils2.FakeJob(0)
		parsers.UpdateJobModelFromEntities(jobModel, jobEntity)

		mockDB.EXPECT().Job().Return(mockJobDA).Times(1)
		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, tenantID).Return(jobModel, expectedErr)


		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)
			assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with the same error if update transaction fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils.FakeJobEntity()
		jobModel := testutils2.FakeJob(0)
		parsers.UpdateJobModelFromEntities(jobModel, jobEntity)

		mockDB.EXPECT().Job().Return(mockJobDA).Times(1)
		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, tenantID).Return(jobModel, nil)

		mockDB.EXPECT().Begin().Return(mockDBTX, nil).Times(1)
		mockDBTX.EXPECT().Transaction().Return(mockTransactionDA).Times(1)
		mockDBTX.EXPECT().Rollback().Return(nil).Times(1)
		mockDBTX.EXPECT().Close().Return(nil).Times(1)

		mockTransactionDA.EXPECT().Update(gomock.Any(), jobModel.Transaction).Return(expectedErr).Times(1)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)
			assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with the same error if update job fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils.FakeJobEntity()
		jobModel := testutils2.FakeJob(0)
		parsers.UpdateJobModelFromEntities(jobModel, jobEntity)

		mockDB.EXPECT().Job().Return(mockJobDA).Times(1)
		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, tenantID).Return(jobModel, nil)

		mockDB.EXPECT().Begin().Return(mockDBTX, nil).Times(1)
		mockDBTX.EXPECT().Transaction().Return(mockTransactionDA).Times(1)
		mockDBTX.EXPECT().Job().Return(mockJobDA).Times(1)
		mockDBTX.EXPECT().Rollback().Return(nil).Times(1)
		mockDBTX.EXPECT().Close().Return(nil).Times(1)

		mockTransactionDA.EXPECT().Update(gomock.Any(), jobModel.Transaction).Return(nil).Times(1)
		mockJobDA.EXPECT().Update(gomock.Any(), jobModel).Return(expectedErr).Times(1)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)
			assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})

	t.Run("should fail with the same error if update log fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		jobEntity := testutils.FakeJobEntity()
		jobModel := testutils2.FakeJob(0)
		parsers.UpdateJobModelFromEntities(jobModel, jobEntity)

		mockDB.EXPECT().Job().Return(mockJobDA).Times(1)
		mockJobDA.EXPECT().FindOneByUUID(ctx, jobEntity.UUID, tenantID).Return(jobModel, nil)

		mockDB.EXPECT().Begin().Return(mockDBTX, nil).Times(1)
		mockDBTX.EXPECT().Transaction().Return(mockTransactionDA).Times(1)
		mockDBTX.EXPECT().Job().Return(mockJobDA).Times(1)
		mockDBTX.EXPECT().Log().Return(mockLogDA).Times(1)
		mockDBTX.EXPECT().Rollback().Return(nil).Times(1)
		mockDBTX.EXPECT().Close().Return(nil).Times(1)

		mockTransactionDA.EXPECT().Update(gomock.Any(), jobModel.Transaction).Return(nil).Times(1)
		mockJobDA.EXPECT().Update(gomock.Any(), jobModel).Return(nil).Times(1)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr).Times(1)

		_, err := usecase.Execute(context.Background(), jobEntity, tenantID)
			assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateJobComponent), err)
	})
}
