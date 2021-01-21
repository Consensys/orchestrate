// +build unit

package accounts

import (
	"context"
	parsers2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	models2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models/testutils"
)

func TestSearchAccounts_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	accountAgent := mocks.NewMockAccountAgent(ctrl)
	mockDB.EXPECT().Account().Return(accountAgent).AnyTimes()

	usecase := NewSearchAccountsUseCase(mockDB)

	tenantID := "tenantID"
	tenants := []string{tenantID}

	t.Run("should execute use case successfully", func(t *testing.T) {
		acc := testutils.FakeAccountModel()

		filter := &entities.AccountFilters{
			Aliases: []string{"alias1"},
		}

		accountAgent.EXPECT().Search(gomock.Any(), filter, tenants).Return([]*models2.Account{acc}, nil)

		resp, err := usecase.Execute(ctx, filter, tenants)

		assert.NoError(t, err)
		assert.Equal(t, parsers2.NewAccountEntityFromModels(acc), resp[0])
	})

	t.Run("should fail with same error if search identities fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		filter := &entities.AccountFilters{
			Aliases: []string{"alias1"},
		}

		accountAgent.EXPECT().Search(gomock.Any(), filter, tenants).Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, filter, tenants)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})
}
