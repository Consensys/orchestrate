// +build unit

package chains

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

func TestGetJob_Execute(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chainRegistryClient := mock.NewMockChainRegistryClient(ctrl)
	usecase := NewGetChainByNameUseCase(chainRegistryClient)

	chainModel := &models.Chain{
		Name:     "ChainName",
		UUID:     "ChainUUID",
		TenantID: "tenantID",
	}

	t.Run("should execute use case successfully", func(t *testing.T) {
		chainRegistryClient.EXPECT().
			GetChainByName(gomock.Any(), chainModel.Name).
			Return(chainModel, nil)

		response, err := usecase.Execute(ctx, chainModel.Name, chainModel.TenantID)

		expectedRes := parsers.NewChainFromModels(chainModel)
		assert.NoError(t, err)
		assert.Equal(t, expectedRes, response)
	})

	t.Run("should fail with same error if GetChainByName fails", func(t *testing.T) {
		expectedErr := errors.InvalidArgError("error")

		chainRegistryClient.EXPECT().
			GetChainByName(gomock.Any(), chainModel.Name).
			Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, chainModel.Name, chainModel.TenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(getChainByNameComponent), err)
	})

	t.Run("should fail with same error if GetChainByName by not found", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		chainRegistryClient.EXPECT().
			GetChainByName(gomock.Any(), chainModel.Name).
			Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, chainModel.Name, chainModel.TenantID)

		assert.Nil(t, response)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}
