package faucets

import (
	"context"
	"testing"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	"github.com/consensys/orchestrate/services/api/store/mocks"
	"github.com/consensys/orchestrate/services/api/store/models"
	"github.com/consensys/orchestrate/services/api/store/models/testutils"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSearchFaucets_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	faucetAgent := mocks.NewMockFaucetAgent(ctrl)
	mockDB.EXPECT().Faucet().Return(faucetAgent).AnyTimes()

	usecase := NewSearchFaucets(mockDB)

	tenantID := multitenancy.DefaultTenant
	tenants := []string{tenantID}

	t.Run("should execute use case successfully", func(t *testing.T) {
		filters := &entities.FaucetFilters{
			Names:     []string{"name1", "name2"},
			ChainRule: "chainRule",
		}
		faucet := testutils.FakeFaucetModel()
		faucetAgent.EXPECT().Search(gomock.Any(), filters, tenants).Return([]*models.Faucet{faucet}, nil)

		resp, err := usecase.Execute(ctx, filters, tenants)

		assert.NoError(t, err)
		assert.Equal(t, []*entities.Faucet{parsers.NewFaucetFromModel(faucet)}, resp)
	})

	t.Run("should fail with same error if search faucets fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")

		faucetAgent.EXPECT().Search(gomock.Any(), nil, tenants).Return(nil, expectedErr)

		resp, err := usecase.Execute(ctx, nil, tenants)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(searchFaucetsComponent), err)
	})
}
