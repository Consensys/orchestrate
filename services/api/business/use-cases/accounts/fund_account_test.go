// +build unit

package accounts

import (
	"context"
	"fmt"
	"github.com/consensys/orchestrate/pkg/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/business/use-cases/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/types/testutils"
)

var (
	faucetNotFoundErr = errors.NotFoundError("not found faucet candidate")
)

func TestFundingAccount_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tenantID := "tenantID"
	allowedTenants := []string{tenantID, multitenancy.DefaultTenant}

	mockSearchChainsUC := mocks.NewMockSearchChainsUseCase(ctrl)
	mockGetFaucetCandidate := mocks.NewMockGetFaucetCandidateUseCase(ctrl)
	mockSendTxUC := mocks.NewMockSendTxUseCase(ctrl)

	usecase := NewFundAccountUseCase(mockSearchChainsUC, mockSendTxUC, mockGetFaucetCandidate)

	t.Run("should trigger funding identity successfully", func(t *testing.T) {
		account := testutils.FakeAccount()
		chains := []*entities.Chain{testutils.FakeChain()}
		faucet := testutils.FakeFaucet()
		chainName := "besu"

		mockSearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{chainName}}, allowedTenants).Return(chains, nil)
		mockGetFaucetCandidate.EXPECT().Execute(gomock.Any(), account.Address, chains[0], allowedTenants).Return(faucet, nil)
		mockSendTxUC.EXPECT().Execute(gomock.Any(), gomock.Any(), "", tenantID).Return(nil, nil)

		err := usecase.Execute(ctx, account, chainName, tenantID)

		assert.NoError(t, err)
	})

	t.Run("should do nothing if there is not faucet candidates", func(t *testing.T) {
		account := testutils.FakeAccount()
		chains := []*entities.Chain{testutils.FakeChain()}
		chainName := "besu"

		mockSearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{chainName}}, allowedTenants).Return(chains, nil)
		mockGetFaucetCandidate.EXPECT().Execute(gomock.Any(), account.Address, chains[0], allowedTenants).Return(nil, faucetNotFoundErr)

		err := usecase.Execute(ctx, account, chainName, tenantID)

		assert.NoError(t, err)
	})

	t.Run("should fail with InvalidParameter if no chains are found", func(t *testing.T) {
		account := testutils.FakeAccount()
		chainName := "besu"

		mockSearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{chainName}}, allowedTenants).Return([]*entities.Chain{}, nil)

		err := usecase.Execute(ctx, account, chainName, tenantID)

		assert.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with same error if search chain fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")
		account := testutils.FakeAccount()
		chainName := "besu"

		mockSearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{chainName}}, allowedTenants).Return(nil, expectedErr)

		err := usecase.Execute(ctx, account, chainName, tenantID)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(fundAccountComponent), err)
	})

	t.Run("should fail with same error if get faucet candidate fails", func(t *testing.T) {
		expectedErr := errors.ConnectionError("error")
		account := testutils.FakeAccount()
		chains := []*entities.Chain{testutils.FakeChain()}
		chainName := "besu"

		mockSearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{chainName}}, allowedTenants).Return(chains, nil)
		mockGetFaucetCandidate.EXPECT().
			Execute(gomock.Any(), account.Address, gomock.Any(), allowedTenants).
			Return(nil, expectedErr)

		err := usecase.Execute(ctx, account, chainName, tenantID)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(fundAccountComponent), err)
	})

	t.Run("should fail with same error if send funding transaction fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		account := testutils.FakeAccount()
		chains := []*entities.Chain{testutils.FakeChain()}
		faucet := testutils.FakeFaucet()
		chainName := "besu"

		mockSearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{chainName}}, allowedTenants).Return(chains, nil)
		mockGetFaucetCandidate.EXPECT().Execute(gomock.Any(), account.Address, gomock.Any(), allowedTenants).Return(faucet, nil)
		mockSendTxUC.EXPECT().Execute(gomock.Any(), gomock.Any(), "", tenantID).Return(nil, expectedErr)

		err := usecase.Execute(ctx, account, chainName, tenantID)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(fundAccountComponent), err)
	})
}
