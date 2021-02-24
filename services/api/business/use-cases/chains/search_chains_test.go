package chains

import (
	"context"
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/multitenancy"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/services/api/business/parsers"
	"github.com/ConsenSys/orchestrate/services/api/store/mocks"
	"github.com/ConsenSys/orchestrate/services/api/store/models"
	"github.com/ConsenSys/orchestrate/services/api/store/models/testutils"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSearchChains_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	chainAgent := mocks.NewMockChainAgent(ctrl)

	mockDB.EXPECT().Chain().Return(chainAgent).AnyTimes()

	usecase := NewSearchChainsUseCase(mockDB)

	tenantID := multitenancy.DefaultTenant
	tenants := []string{tenantID}

	t.Run("should execute use case successfully", func(t *testing.T) {
		filters := &entities.ChainFilters{
			Names: []string{"name1", "name2"},
		}
		chainModel := testutils.FakeChainModel()

		chainAgent.EXPECT().Search(gomock.Any(), filters, tenants).Return([]*models.Chain{chainModel}, nil)

		resp, err := usecase.Execute(ctx, filters, tenants)

		assert.NoError(t, err)
		assert.Equal(t, []*entities.Chain{parsers.NewChainFromModel(chainModel)}, resp)
	})

	t.Run("should fail with same error if search chains fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")

		chainAgent.EXPECT().Search(gomock.Any(), nil, tenants).Return(nil, expectedErr)

		resp, err := usecase.Execute(ctx, nil, tenants)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(searchChainsComponent), err)
	})
}
