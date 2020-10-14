// +build unit

package ethereum

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

func TestCreateAccount_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVault := mocks.NewMockVault(ctrl)
	mockEthereumDA := mocks.NewMockEthereumAgent(ctrl)
	ctx := context.Background()

	mockVault.EXPECT().Ethereum().Return(mockEthereumDA).AnyTimes()

	usecase := NewCreateAccountUseCase(mockVault)

	t.Run("should execute use case successfully", func(t *testing.T) {
		fakeAccount := testutils.FakeETHAccount()

		mockEthereumDA.EXPECT().Insert(ctx, gomock.Any(), gomock.Any(), fakeAccount.Namespace).Return(nil)

		account, err := usecase.Execute(ctx, fakeAccount)

		assert.NoError(t, err)
		assert.Equal(t, account.Namespace, fakeAccount.Namespace)
		assert.True(t, common.IsHexAddress(account.Address))
	})

	t.Run("should fail with same error if Insert fails", func(t *testing.T) {
		fakeAccount := testutils.FakeETHAccount()
		expectedErr := errors.HashicorpVaultConnectionError("error")

		mockEthereumDA.EXPECT().Insert(ctx, gomock.Any(), gomock.Any(), fakeAccount.Namespace).Return(expectedErr)

		account, err := usecase.Execute(ctx, fakeAccount)
		assert.Nil(t, account)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})
}
