// +build unit

package ethereum

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
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
	namespace := "namespace"

	mockVault.EXPECT().Ethereum().Return(mockEthereumDA).AnyTimes()

	usecase := NewCreateAccountUseCase(mockVault)

	t.Run("should execute use case successfully by generating a private key", func(t *testing.T) {
		mockEthereumDA.EXPECT().Insert(ctx, gomock.Any(), gomock.Any(), namespace).Return(nil)

		account, err := usecase.Execute(ctx, namespace, "")

		assert.NoError(t, err)
		assert.Equal(t, account.Namespace, namespace)
		assert.True(t, common.IsHexAddress(account.Address))
	})

	t.Run("should execute use case successfully by importing a private key", func(t *testing.T) {
		privKey := "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"

		mockEthereumDA.EXPECT().Insert(ctx, gomock.Any(), gomock.Any(), namespace).Return(nil)

		account, err := usecase.Execute(ctx, namespace, privKey)

		assert.NoError(t, err)
		assert.Equal(t, account.Namespace, namespace)
		assert.Equal(t, "0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5", account.Address)
	})

	t.Run("should fail with same error if Insert fails", func(t *testing.T) {
		expectedErr := errors.HashicorpVaultConnectionError("error")

		mockEthereumDA.EXPECT().Insert(ctx, gomock.Any(), gomock.Any(), namespace).Return(expectedErr)

		account, err := usecase.Execute(ctx, namespace, "")
		assert.Nil(t, account)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})
}
