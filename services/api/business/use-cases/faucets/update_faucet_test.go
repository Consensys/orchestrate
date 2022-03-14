// +build unit

package faucets

import (
	"context"
	"testing"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	"github.com/consensys/orchestrate/services/api/store/models/testutils"
	"github.com/stretchr/testify/require"

	testutils2 "github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/services/api/store/mocks"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateFaucet_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	accountAgent := mocks.NewMockAccountAgent(ctrl)
	chainAgent := mocks.NewMockChainAgent(ctrl)
	faucetAgent := mocks.NewMockFaucetAgent(ctrl)
	mockDB.EXPECT().Chain().AnyTimes().Return(chainAgent)
	mockDB.EXPECT().Account().AnyTimes().Return(accountAgent)
	mockDB.EXPECT().Faucet().AnyTimes().Return(faucetAgent)

	chain := testutils.FakeChainModel()
	account := testutils.FakeAccountModel()
	userInfo := multitenancy.NewUserInfo("tenantOne", "username")
	usecase := NewUpdateFaucetUseCase(mockDB)
	faucet := testutils2.FakeFaucet()
	faucet.TenantID = userInfo.TenantID
	faucetModel := parsers.NewFaucetModelFromEntity(faucet)

	expectedErr := errors.NotFoundError("error")
	t.Run("should execute use case successfully", func(t *testing.T) {

		chainAgent.EXPECT().FindOneByUUID(gomock.Any(), faucet.ChainRule, userInfo.AllowedTenants, userInfo.Username).
			Return(chain, nil)
		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), faucetModel.CreditorAccount, userInfo.AllowedTenants,
			userInfo.Username).Return(account, nil)
		faucetAgent.EXPECT().Update(gomock.Any(), faucetModel, userInfo.AllowedTenants).Return(nil)
		faucetAgent.EXPECT().FindOneByUUID(gomock.Any(), faucet.UUID, userInfo.AllowedTenants).Return(faucetModel, nil)

		resp, err := usecase.Execute(ctx, faucet, userInfo)

		assert.NoError(t, err)
		assert.Equal(t, faucet, resp)
	})

	t.Run("should fail with the same error if update faucet fails", func(t *testing.T) {
		chainAgent.EXPECT().FindOneByUUID(gomock.Any(), faucetModel.ChainRule, userInfo.AllowedTenants, userInfo.Username).
			Return(chain, nil)
		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), faucetModel.CreditorAccount, userInfo.AllowedTenants,
			userInfo.Username).Return(account, nil)
		faucetAgent.EXPECT().Update(gomock.Any(), faucetModel, userInfo.AllowedTenants).Return(expectedErr)

		resp, err := usecase.Execute(ctx, faucet, userInfo)

		assert.Nil(t, resp)
		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should fail with invalid params error if cannot find account", func(t *testing.T) {
		chainAgent.EXPECT().FindOneByUUID(gomock.Any(), faucet.ChainRule, userInfo.AllowedTenants, userInfo.Username).
			Return(chain, nil)
		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), faucet.CreditorAccount.String(), userInfo.AllowedTenants,
			userInfo.Username).Return(nil, expectedErr)

		resp, err := usecase.Execute(ctx, faucet, userInfo)

		assert.Nil(t, resp)
		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with invalid params error if cannot find chain", func(t *testing.T) {
		chainAgent.EXPECT().FindOneByUUID(gomock.Any(), faucet.ChainRule, userInfo.AllowedTenants, userInfo.Username).
			Return(nil, expectedErr)

		resp, err := usecase.Execute(ctx, faucet, userInfo)

		assert.Nil(t, resp)
		require.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}
