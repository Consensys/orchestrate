package faucets

import (
	"context"
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/multitenancy"

	testutils2 "github.com/ConsenSys/orchestrate/pkg/types/testutils"
	"github.com/ConsenSys/orchestrate/services/api/business/parsers"
	"github.com/ConsenSys/orchestrate/services/api/store/mocks"

	"github.com/ConsenSys/orchestrate/pkg/errors"
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
	tenantID := multitenancy.DefaultTenant
	tenants := []string{tenantID}

	usecase := NewUpdateFaucetUseCase(mockDB)

	t.Run("should execute use case successfully", func(t *testing.T) {
		faucetModel := parsers.NewFaucetModelFromEntity(faucet)

		faucetAgent.EXPECT().Update(gomock.Any(), faucetModel, tenants).Return(nil)
		faucetAgent.EXPECT().FindOneByUUID(gomock.Any(), faucet.UUID, tenants).Return(faucetModel, nil)

		resp, err := usecase.Execute(ctx, faucet, tenants)

		assert.NoError(t, err)
		assert.Equal(t, parsers.NewFaucetFromModel(faucetModel), resp)
	})

	t.Run("should fail with same error if update faucet fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		faucetAgent.EXPECT().Update(gomock.Any(), gomock.Any(), tenants).Return(expectedErr)

		resp, err := usecase.Execute(ctx, faucet, tenants)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateFaucetComponent), err)
	})

	t.Run("should fail with same error if findOne faucet fails", func(t *testing.T) {
		faucetModel := parsers.NewFaucetModelFromEntity(faucet)
		expectedErr := errors.NotFoundError("error")

		faucetAgent.EXPECT().Update(gomock.Any(), faucetModel, tenants).Return(nil)
		faucetAgent.EXPECT().FindOneByUUID(gomock.Any(), faucet.UUID, tenants).Return(nil, expectedErr)

		resp, err := usecase.Execute(ctx, faucet, tenants)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(updateFaucetComponent), err)
	})
}
