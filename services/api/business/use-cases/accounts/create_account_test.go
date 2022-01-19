// +build unit

package accounts

import (
	"context"
	"fmt"
	testutils2 "github.com/consensys/orchestrate/services/api/store/models/testutils"
	"testing"

	"github.com/consensys/orchestrate/pkg/errors"
	qkm "github.com/consensys/orchestrate/pkg/quorum-key-manager"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	mocks2 "github.com/consensys/orchestrate/services/api/business/use-cases/mocks"
	"github.com/consensys/orchestrate/services/api/store/mocks"
	qkmtypes "github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	qkmmock "github.com/consensys/quorum-key-manager/pkg/client/mock"
)

func TestCreateAccount_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
		acc := qkm.FakeEthAccountResponse(accEntity.Address, userInfo.AllowedTenants)

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo)
		mockClient.EXPECT().CreateEthAccount(gomock.Any(), accEntity.StoreID, &qkmtypes.CreateEthAccountRequest{
			KeyID: generateKeyID(userInfo.TenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		}).Return(acc, nil)
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, accEntity, nil, "", userInfo)

		assert.NoError(t, err)
		assert.Equal(t, resp.PublicKey.String(), accEntity.PublicKey.String())
		assert.Equal(t, resp.Address, accEntity.Address)
	})

	t.Run("should import new account account successfully", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = userInfo.TenantID
		accEntity.OwnerID = userInfo.Username
		privKey := hexutil.MustDecode("0xdb337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c")
		acc := qkm.FakeEthAccountResponse(accEntity.Address, userInfo.AllowedTenants)

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo)
		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), "0x83a0254be47813BBff771F4562744676C4e793F0", userInfo.AllowedTenants, userInfo.Username).Return(nil, errors.NotFoundError("not found"))
		mockClient.EXPECT().ImportEthAccount(gomock.Any(), accEntity.StoreID, &qkmtypes.ImportEthAccountRequest{
			KeyID:      generateKeyID(userInfo.TenantID, accEntity.Alias),
			PrivateKey: privKey,
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		}).Return(acc, nil)
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, accEntity, privKey, "", userInfo)

		assert.NoError(t, err)
		assert.Equal(t, resp.PublicKey.String(), accEntity.PublicKey.String())
		assert.Equal(t, resp.Address, accEntity.Address)
	})

	t.Run("should create account and trigger funding successfully", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = userInfo.TenantID
		accEntity.OwnerID = userInfo.Username
		accEntity.StoreID = "personal-qkm-store-id"
		chainName := "besu"
		acc := qkm.FakeEthAccountResponse(accEntity.Address, []string{userInfo.TenantID})

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo)
		mockFundAccountUC.EXPECT().Execute(gomock.Any(), gomock.Any(), chainName, userInfo).Return(nil)
		mockClient.EXPECT().CreateEthAccount(gomock.Any(), accEntity.StoreID, &qkmtypes.CreateEthAccountRequest{
			KeyID: generateKeyID(userInfo.TenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		}).Return(acc, nil)
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, accEntity, nil, chainName, userInfo)

		assert.NoError(t, err)
		assert.Equal(t, resp.PublicKey.String(), accEntity.PublicKey.String())
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

	t.Run("should fail with AlreadyExistsError if search identities returns values", func(t *testing.T) {
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
		mockClient.EXPECT().CreateEthAccount(gomock.Any(), accEntity.StoreID, &qkmtypes.CreateEthAccountRequest{
			KeyID: generateKeyID(userInfo.TenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		}).Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, accEntity, nil, "", userInfo)
		assert.True(t, errors.IsDependencyFailureError(err))
	})

	t.Run("should fail with InvalidParameterError if private key is invalid", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = userInfo.TenantID
		accEntity.OwnerID = userInfo.Username

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo)

		_, err := usecase.Execute(ctx, accEntity, []byte("invalidPrivKey"), "", userInfo)

		assert.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with AlreadyExistsError if account already exists", func(t *testing.T) {
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = userInfo.TenantID
		accEntity.OwnerID = userInfo.Username
		privKey := hexutil.MustDecode("0xdb337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c")

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo)
		accountAgent.EXPECT().
			FindOneByAddress(gomock.Any(), "0x83a0254be47813BBff771F4562744676C4e793F0", userInfo.AllowedTenants, userInfo.Username).
			Return(testutils2.FakeAccountModel(), nil)

		_, err := usecase.Execute(ctx, accEntity, privKey, "", userInfo)

		assert.Error(t, err)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})

	t.Run("should fail with same error if fail to get account when importing", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = userInfo.TenantID
		accEntity.OwnerID = userInfo.Username
		privKey := hexutil.MustDecode("0xdb337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c")

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo)
		accountAgent.EXPECT().FindOneByAddress(gomock.Any(), "0x83a0254be47813BBff771F4562744676C4e793F0", userInfo.AllowedTenants, userInfo.Username).Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, accEntity, privKey, "", userInfo)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})

	t.Run("should fail with same error if cannot insert account", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		accEntity := testutils.FakeAccount()
		accEntity.TenantID = userInfo.TenantID
		accEntity.OwnerID = userInfo.Username
		acc := qkm.FakeEthAccountResponse(accEntity.Address, userInfo.AllowedTenants)

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo)
		mockClient.EXPECT().CreateEthAccount(gomock.Any(), accEntity.StoreID, &qkmtypes.CreateEthAccountRequest{
			KeyID: generateKeyID(userInfo.TenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		}).Return(acc, nil)
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
		acc := qkm.FakeEthAccountResponse(accEntity.Address, userInfo.AllowedTenants)

		mockSearchUC.EXPECT().Execute(gomock.Any(), &entities.AccountFilters{Aliases: []string{accEntity.Alias}, TenantID: userInfo.TenantID}, userInfo)
		mockClient.EXPECT().CreateEthAccount(gomock.Any(), accEntity.StoreID, &qkmtypes.CreateEthAccountRequest{
			KeyID: generateKeyID(userInfo.TenantID, accEntity.Alias),
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		}).Return(acc, nil)
		accountAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockFundAccountUC.EXPECT().Execute(gomock.Any(), gomock.Any(), chainName, userInfo).Return(expectedErr)
		_, err := usecase.Execute(ctx, accEntity, nil, chainName, userInfo)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})
}
