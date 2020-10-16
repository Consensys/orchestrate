// +build unit

package identity

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

func TestCreateIdentity_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	identityAgent := mocks.NewMockIdentityAgent(ctrl)
	mockDB.EXPECT().Identity().Return(identityAgent).AnyTimes()
	mockSearchUC := mocks2.NewMockSearchIdentitiesUseCase(ctrl)
	mockClient := mock.NewMockKeyManagerClient(ctrl)

	usecase := NewCreateIdentityUseCase(mockDB, mockSearchUC, mockClient)

	tenantID := "tenantID"
	tenants := []string{tenantID}

	t.Run("should create account successfully", func(t *testing.T) {
		idenEntity := testutils3.FakeIdentity()
		idenEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(ctx, &entities.IdentityFilters{Aliases: []string{idenEntity.Alias}}, tenants)
		mockClient.EXPECT().CreateETHAccount(ctx, &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(&types.ETHAccountResponse{
			Address:   idenEntity.Address,
			PublicKey: idenEntity.PublicKey,
			Namespace: tenantID,
		}, nil)

		identityAgent.EXPECT().Insert(ctx, gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, idenEntity, "", tenantID)

		assert.NoError(t, err)
		assert.Equal(t, resp.PublicKey, idenEntity.PublicKey)
		assert.Equal(t, resp.Address, idenEntity.Address)
		assert.Equal(t, resp.Active, true)
	})
	
	t.Run("should import account successfully", func(t *testing.T) {
		idenEntity := testutils3.FakeIdentity()
		idenEntity.TenantID = tenantID
		privateKey := "ETHPrivateKey"

		mockSearchUC.EXPECT().Execute(ctx, &entities.IdentityFilters{Aliases: []string{idenEntity.Alias}}, tenants)
		mockClient.EXPECT().ImportETHAccount(ctx, &types.ImportETHAccountRequest{
			Namespace: tenantID,
			PrivateKey: privateKey,
		}).Return(&types.ETHAccountResponse{
			Address:   idenEntity.Address,
			PublicKey: idenEntity.PublicKey,
			Namespace: tenantID,
		}, nil)

		identityAgent.EXPECT().Insert(ctx, gomock.Any()).Return(nil)

		resp, err := usecase.Execute(ctx, idenEntity, privateKey, tenantID)

		assert.NoError(t, err)
		assert.Equal(t, resp.PublicKey, idenEntity.PublicKey)
		assert.Equal(t, resp.Address, idenEntity.Address)
		assert.Equal(t, resp.Active, true)
	})

	t.Run("should fail with same error if search identities fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		idenEntity := testutils3.FakeIdentity()
		idenEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(ctx, &entities.IdentityFilters{Aliases: []string{idenEntity.Alias}}, tenants).Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, idenEntity, "", tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createIdentityComponent), err)
	})

	t.Run("should fail with same error if search identities returns values", func(t *testing.T) {
		idenEntity := testutils3.FakeIdentity()
		foundIdenEntity := testutils3.FakeIdentity()
		idenEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(ctx, &entities.IdentityFilters{Aliases: []string{idenEntity.Alias}}, tenants).
			Return([]*entities.Identity{foundIdenEntity}, nil)

		_, err := usecase.Execute(ctx, idenEntity, "", tenantID)
		assert.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with same error create account fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		idenEntity := testutils3.FakeIdentity()
		idenEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(ctx, &entities.IdentityFilters{Aliases: []string{idenEntity.Alias}}, tenants)
		mockClient.EXPECT().CreateETHAccount(ctx, &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(nil, expectedErr)

		_, err := usecase.Execute(ctx, idenEntity, "", tenantID)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createIdentityComponent), err)
	})

	t.Run("should fail with same error if cannot insert identity", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		idenEntity := testutils3.FakeIdentity()
		idenEntity.TenantID = tenantID

		mockSearchUC.EXPECT().Execute(ctx, &entities.IdentityFilters{Aliases: []string{idenEntity.Alias}}, tenants)
		mockClient.EXPECT().CreateETHAccount(ctx, &types.CreateETHAccountRequest{
			Namespace: tenantID,
		}).Return(&types.ETHAccountResponse{
			Address:   idenEntity.Address,
			PublicKey: idenEntity.PublicKey,
			Namespace: tenantID,
		}, nil)

		identityAgent.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)

		_, err := usecase.Execute(ctx, idenEntity, "", tenantID)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createIdentityComponent), err)
	})
}
