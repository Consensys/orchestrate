// +build unit

package accounts

import (
	"context"
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	qkm "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager"
	qkmtypes "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager/types"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/pkg/types/testutils"
	mocks2 "github.com/ConsenSys/orchestrate/services/api/business/use-cases/mocks"
	"github.com/ConsenSys/orchestrate/services/api/store/mocks"
	"github.com/consensys/quorum/common/hexutil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	qkmmock "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager/client/mocks"
)

func TestCreateAccount_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	globalStoreName := "test-store-name"
	qkm.SetGlobalStoreName(globalStoreName)

	tenantID := "tenantID"
	tenants := []string{tenantID}
	mockDB := mocks.NewMockDB(ctrl)
	accountAgent := mocks.NewMockAccountAgent(ctrl)
	mockSearchUC := mocks2.NewMockSearchAccountsUseCase(ctrl)
	mockFundAccountUC := mocks2.NewMockFundAccountUseCase(ctrl)
	mockClient := qkmmock.NewMockKeyManagerClient(ctrl)

	mockDB.EXPECT().Account().Return(accountAgent).AnyTimes()

	usecase := NewCreateAccountUseCase(mockDB, mockSearchUC, mockFundAccountUC, mockClient)

	t.Run("should create new account account successfully", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID
		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		acc := qkm.FakeEth1AccountResponse(accEntity.Address, []string{tenantID})
		mockClient.EXPECT().CreateEth1Account(gomock.Any(), globalStoreName, &qkmtypes.CreateEth1AccountRequest{
			ID: generateAccountID(tenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants: tenantID,
			},
		}).Return(acc, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address, []string{tenantID}).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, accEntity, nil, "", tenantID)

		assert.NoError(t, err)
		assert.Equal(t, resp.PublicKey, accEntity.PublicKey)
		assert.Equal(t, resp.Address, accEntity.Address)
	})

	t.Run("should import new account account successfully", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID
		privateKey := "1234"
		bPrivKey, _ := hexutil.Decode("0x" + privateKey)

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		acc := qkm.FakeEth1AccountResponse(accEntity.Address, []string{tenantID})
		mockClient.EXPECT().ImportEth1Account(gomock.Any(), globalStoreName, &qkmtypes.ImportEth1AccountRequest{
			ID:         generateAccountID(tenantID, accEntity.Alias),
			PrivateKey: bPrivKey,
			Tags: map[string]string{
				qkm.TagIDAllowedTenants: tenantID,
			},
		}).Return(acc, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address, []string{tenantID}).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, accEntity, bPrivKey, "", tenantID)

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
		acc := qkm.FakeEth1AccountResponse(accEntity.Address, []string{tenantID})
		mockClient.EXPECT().CreateEth1Account(gomock.Any(), globalStoreName, &qkmtypes.CreateEth1AccountRequest{
			ID: generateAccountID(tenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants: tenantID,
			},
		}).Return(acc, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address, []string{tenantID}).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, accEntity, nil, chainName, tenantID)

		assert.NoError(t, err)
		assert.Equal(t, resp.PublicKey, accEntity.PublicKey)
		assert.Equal(t, resp.Address, accEntity.Address)
	})

	t.Run("should fail with same error if search identities fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants).Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, accEntity, nil, "", tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})

	t.Run("should fail with same error if search identities returns values", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		foundAccEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants).
			Return([]*entities.Account{foundAccEntity}, nil)

		_, err := usecase.Execute(ctx, accEntity, nil, "", tenantID)
		assert.Error(t, err)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})

	t.Run("should fail with same error create account fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		mockClient.EXPECT().CreateEth1Account(gomock.Any(), globalStoreName, &qkmtypes.CreateEth1AccountRequest{
			ID: generateAccountID(tenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants: tenantID,
			},
		}).Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, accEntity, nil, "", tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})

	t.Run("should fail with same error if cannot findOneByAddress account", func(t *testing.T) {
		expectedErr := errors.ConnectionError("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		acc := qkm.FakeEth1AccountResponse(accEntity.Address, []string{tenantID})
		mockClient.EXPECT().CreateEth1Account(gomock.Any(), globalStoreName, &qkmtypes.CreateEth1AccountRequest{
			ID: generateAccountID(tenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants: tenantID,
			},
		}).Return(acc, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address, []string{tenantID}).
			Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, accEntity, nil, "", tenantID)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})

	t.Run("should fail with AlreadyExistsError if account already exists", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		acc := qkm.FakeEth1AccountResponse(accEntity.Address, []string{tenantID})
		mockClient.EXPECT().CreateEth1Account(gomock.Any(), globalStoreName, &qkmtypes.CreateEth1AccountRequest{
			ID: generateAccountID(tenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants: tenantID,
			},
		}).Return(acc, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address, []string{tenantID}).
			Return(nil, nil)

		_, err := usecase.Execute(ctx, accEntity, nil, "", tenantID)

		assert.Error(t, err)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})

	t.Run("should fail with same error if cannot insert account", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		acc := qkm.FakeEth1AccountResponse(accEntity.Address, []string{tenantID})
		mockClient.EXPECT().CreateEth1Account(gomock.Any(), globalStoreName, &qkmtypes.CreateEth1AccountRequest{
			ID: generateAccountID(tenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants: tenantID,
			},
		}).Return(acc, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address, []string{tenantID}).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(ctx, accEntity, nil, "", tenantID)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})

	t.Run("should fail with same error if cannot trigger funding account", func(t *testing.T) {
		expectedErr := errors.ConnectionError("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = tenantID
		chainName := "besu"
	
		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}}, tenants)
		acc := qkm.FakeEth1AccountResponse(accEntity.Address, []string{tenantID})
		mockClient.EXPECT().CreateEth1Account(gomock.Any(), globalStoreName, &qkmtypes.CreateEth1AccountRequest{
			ID: generateAccountID(tenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants: tenantID,
			},
		}).Return(acc, nil)
	
		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address, []string{tenantID}).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockFundAccountUC.EXPECT().Execute(gomock.Any(), gomock.Any(), chainName, tenantID).Return(expectedErr)
		_, err := usecase.Execute(ctx, accEntity, nil, chainName, tenantID)
	
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})
}
