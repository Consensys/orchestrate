package faucets

import (
	"context"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/mocks"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegisterFaucet_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	faucet := testutils.FakeFaucet()
	mockDB := mocks.NewMockDB(ctrl)
	faucetAgent := mocks.NewMockFaucetAgent(ctrl)
	mockSearchFaucetsUC := mocks2.NewMockSearchFaucetsUseCase(ctrl)

	mockDB.EXPECT().Faucet().Return(faucetAgent).AnyTimes()

	usecase := NewRegisterFaucetUseCase(mockDB, mockSearchFaucetsUC)

	t.Run("should execute use case successfully", func(t *testing.T) {
		faucetModel := parsers.NewFaucetModelFromEntity(faucet)

		mockSearchFaucetsUC.EXPECT().Execute(gomock.Any(), &entities.FaucetFilters{Names: []string{faucet.Name}}, []string{faucet.TenantID}).Return([]*entities.Faucet{}, nil)
		faucetAgent.EXPECT().Insert(gomock.Any(), faucetModel).Return(nil)

		resp, err := usecase.Execute(ctx, faucet)

		assert.NoError(t, err)
		assert.Equal(t, parsers.NewFaucetFromModel(faucetModel), resp)
	})

	t.Run("should fail with AlreadyExistsError if search faucets returns results", func(t *testing.T) {
		mockSearchFaucetsUC.EXPECT().
			Execute(gomock.Any(), &entities.FaucetFilters{Names: []string{faucet.Name}}, []string{faucet.TenantID}).
			Return([]*entities.Faucet{faucet}, nil)

		resp, err := usecase.Execute(ctx, faucet)

		assert.Nil(t, resp)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})

	t.Run("should fail with same error if search faucets fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		mockSearchFaucetsUC.EXPECT().Execute(gomock.Any(), &entities.FaucetFilters{Names: []string{faucet.Name}}, []string{faucet.TenantID}).Return(nil, expectedErr)

		resp, err := usecase.Execute(ctx, faucet)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(registerFaucetComponent), err)
	})

	t.Run("should fail with same error if insert faucet fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		mockSearchFaucetsUC.EXPECT().Execute(gomock.Any(), &entities.FaucetFilters{Names: []string{faucet.Name}}, []string{faucet.TenantID}).Return([]*entities.Faucet{}, nil)
		faucetAgent.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		resp, err := usecase.Execute(ctx, faucet)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(registerFaucetComponent), err)
	})
}
