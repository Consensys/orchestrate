package chains

import (
	"context"
	"testing"

	testutils2 "github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	mocks2 "github.com/consensys/orchestrate/services/api/business/use-cases/mocks"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/services/api/store/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteChain_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockDBTX := mocks.NewMockTx(ctrl)
	chainAgent := mocks.NewMockChainAgent(ctrl)
	privateTxManagerAgent := mocks.NewMockPrivateTxManagerAgent(ctrl)
	getChainUC := mocks2.NewMockGetChainUseCase(ctrl)

	mockDB.EXPECT().Begin().Return(mockDBTX, nil).AnyTimes()
	mockDB.EXPECT().Chain().Return(chainAgent).AnyTimes()
	mockDBTX.EXPECT().Chain().Return(chainAgent).AnyTimes()
	mockDBTX.EXPECT().PrivateTxManager().Return(privateTxManagerAgent).AnyTimes()
	mockDBTX.EXPECT().Commit().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Rollback().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Close().Return(nil).AnyTimes()

	usecase := NewDeleteChainUseCase(mockDB, getChainUC)
	userInfo := multitenancy.NewUserInfo("tenantOne", "username")

	t.Run("should execute use case successfully", func(t *testing.T) {
		chain := testutils2.FakeChain()
		chainModel := parsers.NewChainModelFromEntity(chain)
		chainModel.TenantID = userInfo.TenantID
		chainModel.OwnerID = userInfo.Username

		getChainUC.EXPECT().Execute(gomock.Any(), "uuid", userInfo).Return(chain, nil)
		privateTxManagerAgent.EXPECT().Delete(gomock.Any(), chainModel.PrivateTxManagers[0]).Return(nil)
		chainAgent.EXPECT().Delete(gomock.Any(), chainModel, userInfo.AllowedTenants).Return(nil)

		err := usecase.Execute(ctx, "uuid", userInfo)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if get chain fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		getChainUC.EXPECT().Execute(gomock.Any(), "uuid", userInfo).Return(nil, expectedErr)

		err := usecase.Execute(ctx, "uuid", userInfo)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(deleteChainComponent), err)
	})

	t.Run("should fail with same error if delete private tx manager fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		chain := testutils2.FakeChain()
		chainModel := parsers.NewChainModelFromEntity(chain)

		getChainUC.EXPECT().Execute(gomock.Any(), "uuid", userInfo).Return(chain, nil)
		privateTxManagerAgent.EXPECT().Delete(gomock.Any(), chainModel.PrivateTxManagers[0]).Return(expectedErr)

		err := usecase.Execute(ctx, "uuid", userInfo)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(deleteChainComponent), err)
	})

	t.Run("should fail with same error if delete chain fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		chain := testutils2.FakeChain()
		chain.TenantID = userInfo.TenantID
		chain.OwnerID = userInfo.Username

		getChainUC.EXPECT().Execute(gomock.Any(), "uuid", userInfo).Return(chain, nil)
		privateTxManagerAgent.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
		chainAgent.EXPECT().Delete(gomock.Any(), parsers.NewChainModelFromEntity(chain), userInfo.AllowedTenants).Return(expectedErr)

		err := usecase.Execute(ctx, "uuid", userInfo)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(deleteChainComponent), err)
	})
}
