// +build unit

package account

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager/ethereum"
	testutils3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/use-cases/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client/mock"
)

func TestCreateAccount_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	accountAgent := mocks.NewMockAccountAgent(ctrl)
	mockDB.EXPECT().Account().Return(accountAgent).AnyTimes()
	mockSearchUC := mocks2.NewMockSearchAccountsUseCase(ctrl)
	mockFundingUC := mocks2.NewMockFundingAccountUseCase(ctrl)
	mockClient := mock.NewMockKeyManagerClient(ctrl)

	usecase := NewCreateAccountUseCase(mockDB, mockSearchUC, mockFundingUC, mockClient)

	tenantID := "tenantID"
	tenants := []string{tenantID}

	t.Run("should create new account account successfully", func(t *testing.T) {
		accountEntity := testutils3.FakeAccount()
		accountEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(ctx, &entities.AccountFilters{Aliases: []string{accountEntity.Alias}}, tenants)
		mockClient.EXPECT().ETHCreateAccount(ctx, &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(&types.ETHAccountResponse{
			Address:   accountEntity.Address,
			PublicKey: accountEntity.PublicKey,
			Namespace: tenantID,
		}, nil)

		accountAgent.EXPECT().FindOneByAddress(ctx, accountEntity.Address, []string{tenantID}).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(ctx, gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, accountEntity, "", "", tenantID)

		assert.NoError(t, err)
		assert.Equal(t, resp.PublicKey, accountEntity.PublicKey)
		assert.Equal(t, resp.Address, accountEntity.Address)
	})

	t.Run("should import new account account successfully", func(t *testing.T) {
		accEntity := testutils3.FakeAccount()
		accEntity.TenantID = tenantID
		privateKey := "ETHPrivateKey"

		mockSearchUC.EXPECT().Execute(ctx, &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		mockClient.EXPECT().ETHImportAccount(ctx, &types.ImportETHAccountRequest{
			Namespace:  tenantID,
			PrivateKey: privateKey,
		}).Return(&types.ETHAccountResponse{
			Address:   accEntity.Address,
			PublicKey: accEntity.PublicKey,
			Namespace: tenantID,
		}, nil)

		accountAgent.EXPECT().FindOneByAddress(ctx, accEntity.Address, []string{tenantID}).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(ctx, gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, accEntity, privateKey, "", tenantID)

		assert.NoError(t, err)
		assert.Equal(t, resp.PublicKey, accEntity.PublicKey)
		assert.Equal(t, resp.Address, accEntity.Address)
	})

	t.Run("should create account and trigger funding successfully", func(t *testing.T) {
		accEntity := testutils3.FakeAccount()
		accEntity.TenantID = tenantID
		chainName := "besu"

		mockSearchUC.EXPECT().Execute(ctx, &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		mockFundingUC.EXPECT().Execute(ctx, gomock.Any(), chainName).Return(nil)
		mockClient.EXPECT().ETHCreateAccount(ctx, &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(&types.ETHAccountResponse{
			Address:   accEntity.Address,
			PublicKey: accEntity.PublicKey,
			Namespace: tenantID,
		}, nil)

		accountAgent.EXPECT().FindOneByAddress(ctx, accEntity.Address, []string{tenantID}).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(ctx, gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, accEntity, "", chainName, tenantID)

		assert.NoError(t, err)
		assert.Equal(t, resp.PublicKey, accEntity.PublicKey)
		assert.Equal(t, resp.Address, accEntity.Address)
	})

	t.Run("should fail with same error if search identities fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		accEntity := testutils3.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(ctx, &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants).Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, accEntity, "", "", tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})

	t.Run("should fail with same error if search identities returns values", func(t *testing.T) {
		accEntity := testutils3.FakeAccount()
		foundAccEntity := testutils3.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(ctx, &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants).
			Return([]*entities.Account{foundAccEntity}, nil)

		_, err := usecase.Execute(ctx, accEntity, "", "", tenantID)
		assert.Error(t, err)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})

	t.Run("should fail with same error create account fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		accEntity := testutils3.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(ctx, &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		mockClient.EXPECT().ETHCreateAccount(ctx, &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, accEntity, "", "", tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})
	
	t.Run("should fail with same error if cannot findOneByAddress account", func(t *testing.T) {
		expectedErr := errors.ConnectionError("error")
		accEntity := testutils3.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(ctx, &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		mockClient.EXPECT().ETHCreateAccount(ctx, &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(&types.ETHAccountResponse{
			Address:   accEntity.Address,
			PublicKey: accEntity.PublicKey,
			Namespace: tenantID,
		}, nil)

		accountAgent.EXPECT().FindOneByAddress(ctx, accEntity.Address, []string{tenantID}).
			Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, accEntity, "", "", tenantID)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})
	
	t.Run("should fail with AlreadyExistsError if account already exists", func(t *testing.T) {
		accEntity := testutils3.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(ctx, &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		mockClient.EXPECT().ETHCreateAccount(ctx, &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(&types.ETHAccountResponse{
			Address:   accEntity.Address,
			PublicKey: accEntity.PublicKey,
			Namespace: tenantID,
		}, nil)

		accountAgent.EXPECT().FindOneByAddress(ctx, accEntity.Address, []string{tenantID}).
			Return(nil, nil)

		_, err := usecase.Execute(ctx, accEntity, "", "", tenantID)

		assert.Error(t, err)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})

	t.Run("should fail with same error if cannot insert account", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		accEntity := testutils3.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(ctx, &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		mockClient.EXPECT().ETHCreateAccount(ctx, &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(&types.ETHAccountResponse{
			Address:   accEntity.Address,
			PublicKey: accEntity.PublicKey,
			Namespace: tenantID,
		}, nil)

		accountAgent.EXPECT().FindOneByAddress(ctx, accEntity.Address, []string{tenantID}).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(ctx, accEntity, "", "", tenantID)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})

	t.Run("should fail with same error if cannot trigger funding account", func(t *testing.T) {
		expectedErr := errors.ConnectionError("error")
		accEntity := testutils3.FakeAccount()
		accEntity.TenantID = tenantID
		chainName := "besu"

		mockSearchUC.EXPECT().Execute(ctx, &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		mockClient.EXPECT().ETHCreateAccount(ctx, &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(&types.ETHAccountResponse{
			Address:   accEntity.Address,
			PublicKey: accEntity.PublicKey,
		}, nil)

		accountAgent.EXPECT().FindOneByAddress(ctx, accEntity.Address, []string{tenantID}).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		mockFundingUC.EXPECT().Execute(ctx, gomock.Any(), chainName).Return(expectedErr)
		_, err := usecase.Execute(ctx, accEntity, "", chainName, tenantID)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})
}
