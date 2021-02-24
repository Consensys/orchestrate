// +build unit

package accounts

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	types "github.com/ConsenSys/orchestrate/pkg/types/keymanager/ethereum"
	"github.com/ConsenSys/orchestrate/pkg/types/testutils"
	mocks2 "github.com/ConsenSys/orchestrate/services/api/business/use-cases/mocks"
	"github.com/ConsenSys/orchestrate/services/api/store/mocks"

	"github.com/ConsenSys/orchestrate/services/key-manager/client/mock"
)

func TestCreateAccount_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tenantID := "tenantID"
	tenants := []string{tenantID}
	mockDB := mocks.NewMockDB(ctrl)
	accountAgent := mocks.NewMockAccountAgent(ctrl)
	mockSearchUC := mocks2.NewMockSearchAccountsUseCase(ctrl)
	mockFundAccountUC := mocks2.NewMockFundAccountUseCase(ctrl)
	mockClient := mock.NewMockKeyManagerClient(ctrl)

	mockDB.EXPECT().Account().Return(accountAgent).AnyTimes()

	usecase := NewCreateAccountUseCase(mockDB, mockSearchUC, mockFundAccountUC, mockClient)

	t.Run("should create new account account successfully", func(t *testing.T) {
		accountEntity := testutils.FakeAccount()
		accountEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accountEntity.Alias}}, tenants)
		mockClient.EXPECT().ETHCreateAccount(gomock.Any(), &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(&types.ETHAccountResponse{
			Address:   accountEntity.Address,
			PublicKey: accountEntity.PublicKey,
			Namespace: tenantID,
		}, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accountEntity.Address, []string{}).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, accountEntity, "", "", tenantID)

		assert.NoError(t, err)
		assert.Equal(t, resp.PublicKey, accountEntity.PublicKey)
		assert.Equal(t, resp.Address, accountEntity.Address)
	})

	t.Run("should import new account account successfully", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID
		privateKey := "ETHPrivateKey"

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		mockClient.EXPECT().ETHImportAccount(gomock.Any(), &types.ImportETHAccountRequest{
			Namespace:  tenantID,
			PrivateKey: privateKey,
		}).Return(&types.ETHAccountResponse{
			Address:   accEntity.Address,
			PublicKey: accEntity.PublicKey,
			Namespace: tenantID,
		}, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address, []string{}).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, accEntity, privateKey, "", tenantID)

		assert.NoError(t, err)
		assert.Equal(t, resp.PublicKey, accEntity.PublicKey)
		assert.Equal(t, resp.Address, accEntity.Address)
	})

	t.Run("should create account and trigger funding successfully", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID
		chainName := "besu"

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		mockFundAccountUC.EXPECT().Execute(gomock.Any(), gomock.Any(), chainName, tenantID).Return(nil)
		mockClient.EXPECT().ETHCreateAccount(gomock.Any(), &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(&types.ETHAccountResponse{
			Address:   accEntity.Address,
			PublicKey: accEntity.PublicKey,
			Namespace: tenantID,
		}, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address, []string{}).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, accEntity, "", chainName, tenantID)

		assert.NoError(t, err)
		assert.Equal(t, resp.PublicKey, accEntity.PublicKey)
		assert.Equal(t, resp.Address, accEntity.Address)
	})

	t.Run("should fail with same error if search identities fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants).Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, accEntity, "", "", tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})

	t.Run("should fail with same error if search identities returns values", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		foundAccEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants).
			Return([]*entities.Account{foundAccEntity}, nil)

		_, err := usecase.Execute(ctx, accEntity, "", "", tenantID)
		assert.Error(t, err)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})

	t.Run("should fail with same error create account fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		mockClient.EXPECT().ETHCreateAccount(gomock.Any(), &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, accEntity, "", "", tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})

	t.Run("should fail with same error if cannot findOneByAddress account", func(t *testing.T) {
		expectedErr := errors.ConnectionError("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		mockClient.EXPECT().ETHCreateAccount(gomock.Any(), &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(&types.ETHAccountResponse{
			Address:   accEntity.Address,
			PublicKey: accEntity.PublicKey,
			Namespace: tenantID,
		}, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address, []string{}).
			Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, accEntity, "", "", tenantID)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})

	t.Run("should fail with AlreadyExistsError if account already exists", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		mockClient.EXPECT().ETHCreateAccount(gomock.Any(), &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(&types.ETHAccountResponse{
			Address:   accEntity.Address,
			PublicKey: accEntity.PublicKey,
			Namespace: tenantID,
		}, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address, []string{}).
			Return(nil, nil)

		_, err := usecase.Execute(ctx, accEntity, "", "", tenantID)

		assert.Error(t, err)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})

	t.Run("should fail with same error if cannot insert account", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		mockClient.EXPECT().ETHCreateAccount(gomock.Any(), &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(&types.ETHAccountResponse{
			Address:   accEntity.Address,
			PublicKey: accEntity.PublicKey,
			Namespace: tenantID,
		}, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address, []string{}).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(ctx, accEntity, "", "", tenantID)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})

	t.Run("should fail with same error if cannot trigger funding account", func(t *testing.T) {
		expectedErr := errors.ConnectionError("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID
		chainName := "besu"

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		mockClient.EXPECT().ETHCreateAccount(gomock.Any(), &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(&types.ETHAccountResponse{
			Address:   accEntity.Address,
			PublicKey: accEntity.PublicKey,
		}, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address, []string{}).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockFundAccountUC.EXPECT().Execute(gomock.Any(), gomock.Any(), chainName, tenantID).Return(expectedErr)
		_, err := usecase.Execute(ctx, accEntity, "", chainName, tenantID)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})
}
