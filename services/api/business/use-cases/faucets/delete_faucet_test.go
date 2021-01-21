package faucets

import (
	"context"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models/testutils"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteFaucet_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	faucetAgent := mocks.NewMockFaucetAgent(ctrl)
	mockDB.EXPECT().Faucet().Return(faucetAgent).AnyTimes()

	usecase := NewDeleteFaucetUseCase(mockDB)

	tenantID := multitenancy.DefaultTenant
	tenants := []string{tenantID}

	t.Run("should execute use case successfully", func(t *testing.T) {
		faucetModel := testutils.FakeFaucetModel()

		faucetAgent.EXPECT().FindOneByUUID(gomock.Any(), "uuid", tenants).Return(faucetModel, nil)
		faucetAgent.EXPECT().Delete(gomock.Any(), faucetModel, tenants).Return(nil)

		err := usecase.Execute(ctx, "uuid", tenants)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if findOne faucet fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		faucetAgent.EXPECT().FindOneByUUID(gomock.Any(), "uuid", tenants).Return(nil, expectedErr)

		err := usecase.Execute(ctx, "uuid", tenants)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(deleteFaucetComponent), err)
	})

	t.Run("should fail with same error if delete faucet fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		faucetAgent.EXPECT().FindOneByUUID(gomock.Any(), "uuid", tenants).Return(testutils.FakeFaucetModel(), nil)
		faucetAgent.EXPECT().Delete(gomock.Any(), gomock.Any(), tenants).Return(expectedErr)

		err := usecase.Execute(ctx, "uuid", tenants)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(deleteFaucetComponent), err)
	})
}
