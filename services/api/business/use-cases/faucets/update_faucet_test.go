// +build unit

package faucets

import (
	"context"
	"testing"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"

	testutils2 "github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/services/api/business/parsers"
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
	faucetAgent := mocks.NewMockFaucetAgent(ctrl)
	mockDB.EXPECT().Faucet().Return(faucetAgent).AnyTimes()
	faucet := testutils2.FakeFaucet()
	userInfo := multitenancy.NewUserInfo("tenantOne", "username")
	usecase := NewUpdateFaucetUseCase(mockDB)

	t.Run("should execute use case successfully", func(t *testing.T) {
		faucetModel := parsers.NewFaucetModelFromEntity(faucet)
		faucetModel.TenantID = userInfo.TenantID
		faucetAgent.EXPECT().Update(gomock.Any(), faucetModel, userInfo.AllowedTenants).Return(nil)
		faucetAgent.EXPECT().FindOneByUUID(gomock.Any(), faucet.UUID, userInfo.AllowedTenants).Return(faucetModel, nil)

		resp, err := usecase.Execute(ctx, faucet, userInfo)

		assert.NoError(t, err)
		assert.Equal(t, parsers.NewFaucetFromModel(faucetModel), resp)
	})

	t.Run("should fail with same error if update faucet fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		faucetAgent.EXPECT().Update(gomock.Any(), gomock.Any(), userInfo.AllowedTenants).Return(expectedErr)

		resp, err := usecase.Execute(ctx, faucet, userInfo)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateFaucetComponent), err)
	})

	t.Run("should fail with same error if findOne faucet fails", func(t *testing.T) {
		faucetModel := parsers.NewFaucetModelFromEntity(faucet)
		faucetModel.TenantID = userInfo.TenantID
		expectedErr := errors.NotFoundError("error")

		faucetAgent.EXPECT().Update(gomock.Any(), faucetModel, userInfo.AllowedTenants).Return(nil)
		faucetAgent.EXPECT().FindOneByUUID(gomock.Any(), faucet.UUID, userInfo.AllowedTenants).Return(nil, expectedErr)

		resp, err := usecase.Execute(ctx, faucet, userInfo)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateFaucetComponent), err)
	})
}
