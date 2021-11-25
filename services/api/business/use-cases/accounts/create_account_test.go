// +build unit

package accounts

import (
	"context"
	"testing"

	"github.com/consensys/orchestrate/pkg/errors"
	qkm "github.com/consensys/orchestrate/pkg/quorum-key-manager"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	mocks2 "github.com/consensys/orchestrate/services/api/business/use-cases/mocks"
	"github.com/consensys/orchestrate/services/api/store/mocks"
	qkmtypes "github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum/common/hexutil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	qkmmock "github.com/consensys/quorum-key-manager/pkg/client/mock"
)

func TestCreateAccount_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	globalStoreName := "test-store-name"
	qkm.SetGlobalStoreName(globalStoreName)

	mockDB := mocks.NewMockDB(ctrl)
	accountAgent := mocks.NewMockAccountAgent(ctrl)
	mockSearchUC := mocks2.NewMockSearchAccountsUseCase(ctrl)
	mockFundAccountUC := mocks2.NewMockFundAccountUseCase(ctrl)
	mockClient := qkmmock.NewMockKeyManagerClient(ctrl)

	mockDB.EXPECT().Account().Return(accountAgent).AnyTimes()

	userInfo := multitenancy.NewUserInfo("tenantOne", "username")
	usecase := NewCreateAccountUseCase(mockDB, mockSearchUC, mockFundAccountUC, mockClient)

	t.Run("should create new account account successfully", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = userInfo.TenantID
		accEntity.OwnerID = userInfo.Username
		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo)
		acc := qkm.FakeEthAccountResponse(accEntity.Address, userInfo.AllowedTenants)
		mockClient.EXPECT().CreateEthAccount(gomock.Any(), globalStoreName, &qkmtypes.CreateEthAccountRequest{
			KeyID: generateKeyID(userInfo.TenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		}).Return(acc, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address.Hex(), userInfo.AllowedTenants, userInfo.Username).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, accEntity, nil, "", userInfo)

		assert.NoError(t, err)
		assert.Equal(t, resp.PublicKey, accEntity.PublicKey)
		assert.Equal(t, resp.Address, accEntity.Address)
	})

	t.Run("should import new account account successfully", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = userInfo.TenantID
		accEntity.OwnerID = userInfo.Username
		privateKey := "1234"
		bPrivKey, _ := hexutil.Decode("0x" + privateKey)

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo)
		acc := qkm.FakeEthAccountResponse(accEntity.Address, userInfo.AllowedTenants)
		mockClient.EXPECT().ImportEthAccount(gomock.Any(), globalStoreName, &qkmtypes.ImportEthAccountRequest{
			KeyID:      generateKeyID(userInfo.TenantID, accEntity.Alias),
			PrivateKey: bPrivKey,
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		}).Return(acc, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address.Hex(), userInfo.AllowedTenants, userInfo.Username).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, accEntity, bPrivKey, "", userInfo)

		assert.NoError(t, err)
		assert.Equal(t, resp.PublicKey, accEntity.PublicKey)
		assert.Equal(t, resp.Address, accEntity.Address)
	})

	t.Run("should create account and trigger funding successfully", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = userInfo.TenantID
		accEntity.OwnerID = userInfo.Username
		chainName := "besu"

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo)
		mockFundAccountUC.EXPECT().Execute(gomock.Any(), gomock.Any(), chainName, userInfo).Return(nil)
		acc := qkm.FakeEthAccountResponse(accEntity.Address, []string{userInfo.TenantID})
		mockClient.EXPECT().CreateEthAccount(gomock.Any(), globalStoreName, &qkmtypes.CreateEthAccountRequest{
			KeyID: generateKeyID(userInfo.TenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		}).Return(acc, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address.Hex(), userInfo.AllowedTenants, userInfo.Username).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, accEntity, nil, chainName, userInfo)

		assert.NoError(t, err)
		assert.Equal(t, resp.PublicKey, accEntity.PublicKey)
		assert.Equal(t, resp.Address, accEntity.Address)
	})

	t.Run("should fail with same error if search identities fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = userInfo.TenantID
		accEntity.OwnerID = userInfo.Username

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo).
			Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, accEntity, nil, "", userInfo)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})

	t.Run("should fail with same error if search identities returns values", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		foundAccEntity := testutils.FakeAccount()
		accEntity.TenantID = userInfo.TenantID
		accEntity.OwnerID = userInfo.Username

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo).
			Return([]*entities.Account{foundAccEntity}, nil)

		_, err := usecase.Execute(ctx, accEntity, nil, "", userInfo)
		assert.Error(t, err)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})

	t.Run("should fail with same error create account fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = userInfo.TenantID
		accEntity.OwnerID = userInfo.Username

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo)
		mockClient.EXPECT().CreateEthAccount(gomock.Any(), globalStoreName, &qkmtypes.CreateEthAccountRequest{
			KeyID: generateKeyID(userInfo.TenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		}).Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, accEntity, nil, "", userInfo)
		assert.True(t, errors.IsDependencyFailureError(err))
	})

	t.Run("should fail with same error if cannot findOneByAddress account", func(t *testing.T) {
		expectedErr := errors.ConnectionError("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = userInfo.TenantID
		accEntity.OwnerID = userInfo.Username

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo)
		acc := qkm.FakeEthAccountResponse(accEntity.Address, userInfo.AllowedTenants)
		mockClient.EXPECT().CreateEthAccount(gomock.Any(), globalStoreName, &qkmtypes.CreateEthAccountRequest{
			KeyID: generateKeyID(userInfo.TenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		}).Return(acc, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address.Hex(), userInfo.AllowedTenants, userInfo.Username).
			Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, accEntity, nil, "", userInfo)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})

	t.Run("should fail with AlreadyExistsError if account already exists", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = userInfo.TenantID
		accEntity.OwnerID = userInfo.Username

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo)
		acc := qkm.FakeEthAccountResponse(accEntity.Address, userInfo.AllowedTenants)
		mockClient.EXPECT().CreateEthAccount(gomock.Any(), globalStoreName, &qkmtypes.CreateEthAccountRequest{
			KeyID: generateKeyID(userInfo.TenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		}).Return(acc, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address.Hex(), userInfo.AllowedTenants, userInfo.Username).
			Return(nil, nil)

		_, err := usecase.Execute(ctx, accEntity, nil, "", userInfo)

		assert.Error(t, err)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})

	t.Run("should fail with same error if cannot insert account", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = userInfo.TenantID
		accEntity.OwnerID = userInfo.Username

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo)
		acc := qkm.FakeEthAccountResponse(accEntity.Address, userInfo.AllowedTenants)
		mockClient.EXPECT().CreateEthAccount(gomock.Any(), globalStoreName, &qkmtypes.CreateEthAccountRequest{
			KeyID: generateKeyID(userInfo.TenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		}).Return(acc, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address.Hex(), userInfo.AllowedTenants, userInfo.Username).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(ctx, accEntity, nil, "", userInfo)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})

	t.Run("should fail with same error if cannot trigger funding account", func(t *testing.T) {
		expectedErr := errors.ConnectionError("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = userInfo.TenantID
		accEntity.OwnerID = userInfo.Username
		chainName := "besu"

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo)
		acc := qkm.FakeEthAccountResponse(accEntity.Address, userInfo.AllowedTenants)
		mockClient.EXPECT().CreateEthAccount(gomock.Any(), globalStoreName, &qkmtypes.CreateEthAccountRequest{
			KeyID: generateKeyID(userInfo.TenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		}).Return(acc, nil)

		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), accEntity.Address.Hex(), userInfo.AllowedTenants, userInfo.Username).
			Return(nil, errors.NotFoundError("not found"))
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockFundAccountUC.EXPECT().Execute(gomock.Any(), gomock.Any(), chainName, userInfo).Return(expectedErr)
		_, err := usecase.Execute(ctx, accEntity, nil, chainName, userInfo)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})
}
