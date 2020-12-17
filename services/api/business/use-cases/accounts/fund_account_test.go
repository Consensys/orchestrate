// +build unit

package accounts

import (
	"context"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/mocks"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	testutils3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client/mock"
)

var (
	faucetNotFoundErr = errors.NotFoundError("not found faucet candidate")
)

func TestFundingIdentity_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tenantID := "tenantID"
	mockRegisterClient := mock2.NewMockChainRegistryClient(ctrl)
	mockSendTxUC := mocks.NewMockSendTxUseCase(ctrl)

	usecase := NewFundAccountUseCase(mockRegisterClient, mockSendTxUC)

	t.Run("should trigger funding identity successfully", func(t *testing.T) {
		account := testutils3.FakeAccount()
		chain := testutils3.FakeChain()
		faucet := testutils3.FakeFaucet()
		chainName := "besu"

		mockRegisterClient.EXPECT().GetChainByName(ctx, chainName).Return(chain, nil)
		mockRegisterClient.EXPECT().GetFaucetCandidate(ctx, ethcommon.HexToAddress(account.Address), chain.UUID).Return(faucet, nil)
		mockSendTxUC.EXPECT().Execute(ctx, gomock.Any(), "", tenantID).Return(nil, nil)

		err := usecase.Execute(ctx, account, chainName, tenantID)

		assert.NoError(t, err)
	})

	t.Run("should do nothing if there is not faucet candidates", func(t *testing.T) {
		account := testutils3.FakeAccount()
		chain := testutils3.FakeChain()
		chainName := "besu"

		mockRegisterClient.EXPECT().GetChainByName(ctx, chainName).Return(chain, nil)
		mockRegisterClient.EXPECT().
			GetFaucetCandidate(ctx, ethcommon.HexToAddress(account.Address), chain.UUID).
			Return(nil, faucetNotFoundErr)

		err := usecase.Execute(ctx, account, chainName, tenantID)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if search chain fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		account := testutils3.FakeAccount()
		chainName := "besu"

		mockRegisterClient.EXPECT().GetChainByName(ctx, chainName).Return(nil, expectedErr)

		err := usecase.Execute(ctx, account, chainName, tenantID)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(fundAccountComponent), err)
	})

	t.Run("should fail with same error if get faucet candidate fails", func(t *testing.T) {
		expectedErr := errors.ConnectionError("error")
		account := testutils3.FakeAccount()
		chain := testutils3.FakeChain()
		chainName := "besu"

		mockRegisterClient.EXPECT().GetChainByName(ctx, chainName).Return(chain, nil)
		mockRegisterClient.EXPECT().GetFaucetCandidate(ctx, ethcommon.HexToAddress(account.Address), chain.UUID).Return(nil, expectedErr)

		err := usecase.Execute(ctx, account, chainName, tenantID)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(fundAccountComponent), err)
	})

	t.Run("should fail with same error if send funding transaction fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		account := testutils3.FakeAccount()
		chain := testutils3.FakeChain()
		faucet := testutils3.FakeFaucet()
		chainName := "besu"

		mockRegisterClient.EXPECT().GetChainByName(ctx, chainName).Return(chain, nil)
		mockRegisterClient.EXPECT().GetFaucetCandidate(ctx, ethcommon.HexToAddress(account.Address), chain.UUID).Return(faucet, nil)
		mockSendTxUC.EXPECT().Execute(ctx, gomock.Any(), "", tenantID).Return(nil, expectedErr)

		err := usecase.Execute(ctx, account, chainName, tenantID)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(fundAccountComponent), err)
	})
}
