// +build unit

package ethereum

import (
	"context"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

func TestSignPayload_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVault := mocks.NewMockVault(ctrl)
	mockEthereumDA := mocks.NewMockEthereumAgent(ctrl)
	ctx := context.Background()
	address := "0xaddress"
	namespace := "namespace"

	mockVault.EXPECT().Ethereum().Return(mockEthereumDA).AnyTimes()

	usecase := NewSignUseCase(mockVault)

	t.Run("should execute use case successfully", func(t *testing.T) {
		privKey := "5385714a2f6d69ca034f56a5268833216ffb8fba7229c39569bc4c5f42cde97c"
		mockEthereumDA.EXPECT().FindOne(ctx, address, namespace).Return(privKey, nil)

		signature, err := usecase.Execute(ctx, address, namespace, "my data to sign")

		assert.NoError(t, err)
		assert.Equal(t, signature, "0x7107193a8683e258ada2dfa76b5e6fc145ebd98f0e6eee77cb91381201fe7bca5445beccebe164e23abe0639f089e17b24ce867be9fece8b4872cfe13d91464601")
	})

	t.Run("should fail with CryptoOperationError if creation of ECDSA private key fails", func(t *testing.T) {
		mockEthereumDA.EXPECT().FindOne(ctx, address, namespace).Return("invalidPrivKey", nil)

		signature, err := usecase.Execute(ctx, address, namespace, "my data to sign")

		assert.Empty(t, signature)
		assert.True(t, errors.IsCryptoOperationError(err))
	})

	t.Run("should fail with same error if FindOne fails", func(t *testing.T) {
		expectedErr := errors.HashicorpVaultConnectionError("error")

		mockEthereumDA.EXPECT().FindOne(ctx, gomock.Any(), gomock.Any()).Return("", expectedErr)

		signature, err := usecase.Execute(ctx, address, namespace, "my data to sign")

		assert.Empty(t, signature)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})
}
