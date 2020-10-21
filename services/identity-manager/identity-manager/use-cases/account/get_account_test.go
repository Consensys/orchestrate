// +build unit

package account

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/parsers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store/models/testutils"
)

func TestGetAccount_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	accountAgent := mocks.NewMockAccountAgent(ctrl)
	mockDB.EXPECT().Account().Return(accountAgent).AnyTimes()

	usecase := NewGetAccountUseCase(mockDB)

	tenantID := "tenantID"
	tenants := []string{tenantID}

	t.Run("should execute use case successfully", func(t *testing.T) {
		iden := testutils.FakeAccountModel()
		
		accountAgent.EXPECT().FindOneByAddress(ctx, iden.Address, tenants).Return(iden, nil)

		resp, err := usecase.Execute(ctx, iden.Address, tenants)

		assert.NoError(t, err)
		assert.Equal(t, parsers.NewAccountEntityFromModels(iden), resp)
	})
	
	t.Run("should fail with same error if get account fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		acc := testutils.FakeAccountModel()

		accountAgent.EXPECT().FindOneByAddress(ctx, acc.Address, tenants).Return(nil, expectedErr)

		_, err := usecase.Execute(ctx,  acc.Address, tenants)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(getAccountComponent), err)
	})
}
