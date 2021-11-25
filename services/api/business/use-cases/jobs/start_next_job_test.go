// +build unit

package jobs

import (
	"context"
	"testing"
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/consensys/orchestrate/pkg/types/entities"
	mocks2 "github.com/consensys/orchestrate/services/api/business/use-cases/mocks"
	"github.com/consensys/orchestrate/services/api/store/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/consensys/orchestrate/services/api/store/mocks"
	testutils2 "github.com/consensys/orchestrate/services/api/store/models/testutils"
)

func TestStartNextJob_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockJobDA := mocks.NewMockJobAgent(ctrl)
	mockTxDA := mocks.NewMockTransactionAgent(ctrl)
	mockStartJobUC := mocks2.NewMockStartJobUseCase(ctrl)

	mockDB.EXPECT().Job().Return(mockJobDA).AnyTimes()
	mockDB.EXPECT().Transaction().Return(mockTxDA).AnyTimes()

	userInfo := multitenancy.NewUserInfo("tenantOne", "username")
	usecase := NewStartNextJobUseCase(mockDB, mockStartJobUC)

	ctx := context.Background()

	t.Run("should execute use case for EEA marking transaction successfully", func(t *testing.T) {
		jobModel := testutils2.FakeJobModel(0)
		nextJobModel := testutils2.FakeJobModel(0)
		txHash := ethcommon.HexToHash("0x123")

		jobModel.NextJobUUID = nextJobModel.UUID
		jobModel.Transaction.Hash = txHash.String()
		jobModel.Status = entities.StatusStored
		jobModel.Logs = append(jobModel.Logs, &models.Log{
			ID:        1,
			Status:    entities.StatusStored,
			CreatedAt: time.Now(),
		})
		jobModel.Type = entities.EEAPrivateTransaction
		nextJobModel.Type = entities.EEAMarkingTransaction

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobModel.UUID, userInfo.AllowedTenants, userInfo.Username, false).
			Return(jobModel, nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), nextJobModel.UUID, userInfo.AllowedTenants, userInfo.Username, false).
			Return(nextJobModel, nil)
		nextJobModel.Transaction.Data = txHash.String()
		mockTxDA.EXPECT().Update(gomock.Any(), nextJobModel.Transaction).Return(nil)

		mockStartJobUC.EXPECT().Execute(gomock.Any(), nextJobModel.UUID, userInfo)
		err := usecase.Execute(ctx, jobModel.UUID, userInfo)

		assert.NoError(t, err)
	})

	t.Run("should execute use case for tessera marking transaction successfully", func(t *testing.T) {
		jobModel := testutils2.FakeJobModel(0)
		nextJobModel := testutils2.FakeJobModel(0)
		enclaveKey := ethcommon.HexToHash("0x123").String()

		jobModel.NextJobUUID = nextJobModel.UUID
		jobModel.Transaction.EnclaveKey = enclaveKey
		jobModel.Transaction.Gas = "0x1"
		jobModel.Status = entities.StatusStored
		jobModel.Logs = append(jobModel.Logs, &models.Log{
			ID:        1,
			Status:    entities.StatusStored,
			CreatedAt: time.Now(),
		})
		jobModel.Type = entities.TesseraPrivateTransaction
		nextJobModel.Type = entities.TesseraMarkingTransaction

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), jobModel.UUID, userInfo.AllowedTenants, userInfo.Username, false).
			Return(jobModel, nil)
		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), nextJobModel.UUID, userInfo.AllowedTenants, userInfo.Username, false).
			Return(nextJobModel, nil)
		nextJobModel.Transaction.Data = enclaveKey
		nextJobModel.Transaction.Gas = "0x1"
		mockTxDA.EXPECT().Update(gomock.Any(), nextJobModel.Transaction).Return(nil)

		mockStartJobUC.EXPECT().Execute(gomock.Any(), nextJobModel.UUID, userInfo)
		err := usecase.Execute(ctx, jobModel.UUID, userInfo)

		assert.NoError(t, err)
	})
}
