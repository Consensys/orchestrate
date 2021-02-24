package chains

import (
	"context"
	"testing"

	mocks2 "github.com/ConsenSys/orchestrate/services/api/business/use-cases/mocks"

	"github.com/ConsenSys/orchestrate/pkg/multitenancy"

	"github.com/ConsenSys/orchestrate/pkg/types/testutils"
	"github.com/ConsenSys/orchestrate/services/api/business/parsers"
	"github.com/ConsenSys/orchestrate/services/api/store/mocks"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateChain_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockDBTX := mocks.NewMockTx(ctrl)
	chainAgent := mocks.NewMockChainAgent(ctrl)
	mockGetChainUC := mocks2.NewMockGetChainUseCase(ctrl)
	privateTxManagerAgent := mocks.NewMockPrivateTxManagerAgent(ctrl)

	mockDB.EXPECT().Begin().Return(mockDBTX, nil).AnyTimes()
	mockDB.EXPECT().Chain().Return(chainAgent).AnyTimes()
	mockDBTX.EXPECT().Chain().Return(chainAgent).AnyTimes()
	mockDBTX.EXPECT().PrivateTxManager().Return(privateTxManagerAgent).AnyTimes()
	mockDBTX.EXPECT().Commit().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Rollback().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Close().Return(nil).AnyTimes()

	tenantID := multitenancy.DefaultTenant
	tenants := []string{tenantID}

	usecase := NewUpdateChainUseCase(mockDB, mockGetChainUC)

	t.Run("should execute use case successfully", func(t *testing.T) {
		chain := testutils.FakeChain()
		chain.PrivateTxManager = nil
		chainModel := parsers.NewChainModelFromEntity(chain)

		mockGetChainUC.EXPECT().Execute(gomock.Any(), chain.UUID, tenants).Return(chain, nil)
		chainAgent.EXPECT().Update(gomock.Any(), chainModel, tenants).Return(nil)
		mockGetChainUC.EXPECT().Execute(gomock.Any(), chain.UUID, tenants).Return(chain, nil)

		resp, err := usecase.Execute(ctx, chain, tenants)

		assert.NoError(t, err)
		assert.Equal(t, parsers.NewChainFromModel(chainModel), resp)
	})

	t.Run("should execute use case successfully with private tx manager", func(t *testing.T) {
		chain := testutils.FakeChain()
		chainModel := parsers.NewChainModelFromEntity(chain)

		mockGetChainUC.EXPECT().Execute(gomock.Any(), chain.UUID, tenants).Return(chain, nil)
		chainAgent.EXPECT().Update(gomock.Any(), gomock.Any(), tenants).Return(nil)
		privateTxManagerAgent.EXPECT().Update(gomock.Any(), chainModel.PrivateTxManagers[0]).Return(nil)
		mockGetChainUC.EXPECT().Execute(gomock.Any(), chain.UUID, tenants).Return(chain, nil)

		resp, err := usecase.Execute(ctx, chain, tenants)

		assert.NoError(t, err)
		assert.Equal(t, parsers.NewChainFromModel(chainModel), resp)
	})

	t.Run("should execute use case successfully with private tx manager to insert", func(t *testing.T) {
		chainUpdate := testutils.FakeChain()
		chainRetrieved := testutils.FakeChain()
		chainRetrieved.PrivateTxManager = nil
		chainRetrieved.UUID = chainUpdate.UUID
		chainModel := parsers.NewChainModelFromEntity(chainUpdate)
		chainModel.PrivateTxManagers[0].ChainUUID = chainRetrieved.UUID

		mockGetChainUC.EXPECT().Execute(gomock.Any(), chainUpdate.UUID, tenants).Return(chainRetrieved, nil)
		privateTxManagerAgent.EXPECT().Insert(gomock.Any(), chainModel.PrivateTxManagers[0]).Return(nil)
		chainAgent.EXPECT().Update(gomock.Any(), gomock.Any(), tenants).Return(nil)
		mockGetChainUC.EXPECT().Execute(gomock.Any(), chainUpdate.UUID, tenants).Return(chainUpdate, nil)

		_, err := usecase.Execute(ctx, chainUpdate, tenants)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if get chain fails", func(t *testing.T) {
		chain := testutils.FakeChain()
		expectedErr := errors.NotFoundError("error")

		mockGetChainUC.EXPECT().Execute(gomock.Any(), chain.UUID, tenants).Return(nil, expectedErr)

		resp, err := usecase.Execute(ctx, chain, tenants)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateChainComponent), err)
	})

	t.Run("should fail with same error if update private tx manager fails", func(t *testing.T) {
		chain := testutils.FakeChain()
		expectedErr := errors.NotFoundError("error")

		mockGetChainUC.EXPECT().Execute(gomock.Any(), chain.UUID, tenants).Return(chain, nil)
		privateTxManagerAgent.EXPECT().Update(gomock.Any(), gomock.Any()).Return(expectedErr)

		resp, err := usecase.Execute(ctx, chain, tenants)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateChainComponent), err)
	})

	t.Run("should fail with same error if update chain fails", func(t *testing.T) {
		chain := testutils.FakeChain()
		expectedErr := errors.NotFoundError("error")

		mockGetChainUC.EXPECT().Execute(gomock.Any(), chain.UUID, tenants).Return(chain, nil)
		privateTxManagerAgent.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		chainAgent.EXPECT().Update(gomock.Any(), gomock.Any(), tenants).Return(expectedErr)

		resp, err := usecase.Execute(ctx, chain, tenants)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateChainComponent), err)
	})

	t.Run("should fail with same error if get chain fails", func(t *testing.T) {
		chain := testutils.FakeChain()
		expectedErr := errors.NotFoundError("error")

		mockGetChainUC.EXPECT().Execute(gomock.Any(), chain.UUID, tenants).Return(chain, nil)
		chainAgent.EXPECT().Update(gomock.Any(), gomock.Any(), tenants).Return(nil)
		privateTxManagerAgent.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockGetChainUC.EXPECT().Execute(gomock.Any(), chain.UUID, tenants).Return(nil, expectedErr)

		resp, err := usecase.Execute(ctx, chain, tenants)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateChainComponent), err)
	})
}
