// +build unit

package chains

import (
	"context"
	"math/big"
	"testing"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient/mock"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/consensys/orchestrate/pkg/types/entities"
	mocks2 "github.com/consensys/orchestrate/services/api/business/use-cases/mocks"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	"github.com/consensys/orchestrate/services/api/store/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegisterChain_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockDBTX := mocks.NewMockTx(ctrl)
	chainAgent := mocks.NewMockChainAgent(ctrl)
	privateTxManagerAgent := mocks.NewMockPrivateTxManagerAgent(ctrl)
	mockSearchChainsUC := mocks2.NewMockSearchChainsUseCase(ctrl)
	mockEthClient := mock.NewMockClient(ctrl)

	mockDB.EXPECT().Begin().Return(mockDBTX, nil).AnyTimes()
	mockDB.EXPECT().Chain().Return(chainAgent).AnyTimes()
	mockDBTX.EXPECT().Chain().Return(chainAgent).AnyTimes()
	mockDBTX.EXPECT().PrivateTxManager().Return(privateTxManagerAgent).AnyTimes()
	mockDBTX.EXPECT().Commit().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Rollback().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Close().Return(nil).AnyTimes()

	userInfo := multitenancy.NewUserInfo("tenantOne", "username")

	usecase := NewRegisterChainUseCase(mockDB, mockSearchChainsUC, mockEthClient)

	t.Run("should execute use case successfully", func(t *testing.T) {
		chain := testutils.FakeChain()
		chain.PrivateTxManager = nil
		chainModel := parsers.NewChainModelFromEntity(chain)
		chainModel.TenantID = userInfo.TenantID
		chainModel.OwnerID = userInfo.Username

		mockSearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{chain.Name}, TenantID: userInfo.TenantID},
			userInfo).Return([]*entities.Chain{}, nil)
		mockEthClient.EXPECT().Network(gomock.Any(), chain.URLs[0]).Return(big.NewInt(888), nil)
		chainAgent.EXPECT().Insert(gomock.Any(), chainModel).Return(nil)

		resp, err := usecase.Execute(ctx, chain, false, userInfo)

		assert.NoError(t, err)
		assert.Equal(t, parsers.NewChainFromModel(chainModel), resp)
	})

	t.Run("should execute use case successfully from latest block", func(t *testing.T) {
		chain := testutils.FakeChain()
		chain.PrivateTxManager = nil
		chain.TenantID = userInfo.TenantID
		chain.OwnerID = userInfo.Username
		chainTip := big.NewInt(1)

		mockSearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{chain.Name}, TenantID: userInfo.TenantID}, userInfo).
			Return([]*entities.Chain{}, nil)
		mockEthClient.EXPECT().Network(gomock.Any(), chain.URLs[0]).
			Return(big.NewInt(888), nil)
		mockEthClient.EXPECT().HeaderByNumber(gomock.Any(), chain.URLs[0], nil).
			Return(&types.Header{
				Number: chainTip,
			}, nil)
		chainAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, chain, true, userInfo)

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), resp.ListenerStartingBlock)
		assert.Equal(t, uint64(1), resp.ListenerCurrentBlock)
	})

	t.Run("should execute use case successfully with private tx manager", func(t *testing.T) {
		chain := testutils.FakeChain()
		chainModel := parsers.NewChainModelFromEntity(chain)
		chainModel.PrivateTxManagers[0].ChainUUID = chainModel.UUID
		chainModel.TenantID = userInfo.TenantID
		chainModel.OwnerID = userInfo.Username

		mockSearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{chain.Name}, TenantID: userInfo.TenantID},
			userInfo).Return([]*entities.Chain{}, nil)
		mockEthClient.EXPECT().Network(gomock.Any(), chain.URLs[0]).Return(big.NewInt(888), nil)
		chainAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		privateTxManagerAgent.EXPECT().Insert(gomock.Any(), chainModel.PrivateTxManagers[0]).Return(nil)

		resp, err := usecase.Execute(ctx, chain, false, userInfo)

		assert.NoError(t, err)
		assert.Equal(t, parsers.NewChainFromModel(chainModel), resp)
	})

	t.Run("should fail with AlreadyExistsError if search chains returns results", func(t *testing.T) {
		chain := testutils.FakeChain()

		mockSearchChainsUC.EXPECT().
			Execute(gomock.Any(), &entities.ChainFilters{Names: []string{chain.Name}, TenantID: userInfo.TenantID}, userInfo).
			Return([]*entities.Chain{chain}, nil)

		resp, err := usecase.Execute(ctx, chain, false, userInfo)

		assert.Nil(t, resp)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})

	t.Run("should fail with same error if search chains fails", func(t *testing.T) {
		chain := testutils.FakeChain()
		expectedErr := errors.NotFoundError("error")
		chain.TenantID = userInfo.TenantID
		chain.OwnerID = userInfo.Username

		mockSearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{chain.Name}, TenantID: userInfo.TenantID},
			userInfo).Return(nil, expectedErr)

		resp, err := usecase.Execute(ctx, chain, false, userInfo)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(registerChainComponent), err)
	})

	t.Run("should fail with same error if insert chain fails", func(t *testing.T) {
		chain := testutils.FakeChain()
		expectedErr := errors.NotFoundError("error")
		chain.TenantID = userInfo.TenantID
		chain.OwnerID = userInfo.Username

		mockSearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{chain.Name}, TenantID: userInfo.TenantID},
			userInfo).Return([]*entities.Chain{}, nil)
		mockEthClient.EXPECT().Network(gomock.Any(), chain.URLs[0]).Return(big.NewInt(888), nil)
		chainAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		resp, err := usecase.Execute(ctx, chain, false, userInfo)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(registerChainComponent), err)
	})

	t.Run("should fail with same error if insert private tx manager fails", func(t *testing.T) {
		chain := testutils.FakeChain()
		chain.TenantID = userInfo.TenantID
		chain.OwnerID = userInfo.Username
		expectedErr := errors.NotFoundError("error")

		mockSearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{chain.Name}, TenantID: userInfo.TenantID},
			userInfo).Return([]*entities.Chain{}, nil)
		mockEthClient.EXPECT().Network(gomock.Any(), chain.URLs[0]).Return(big.NewInt(888), nil)
		chainAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		privateTxManagerAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		resp, err := usecase.Execute(ctx, chain, false, userInfo)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(registerChainComponent), err)
	})
}
