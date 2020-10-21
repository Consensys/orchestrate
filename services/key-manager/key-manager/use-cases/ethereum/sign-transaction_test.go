// +build unit

package ethereum

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store/mocks"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

func TestSignTransaction_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVault := mocks.NewMockVault(ctrl)
	mockEthereumDA := mocks.NewMockEthereumAgent(ctrl)
	ctx := context.Background()
	address := "0xaddress"
	namespace := "namespace"
	chainID := big.NewInt(1)
	tx := types.NewTransaction(
		0,
		common.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"),
		big.NewInt(10000000000),
		21000,
		big.NewInt(10000000000),
		[]byte{},
	)

	mockVault.EXPECT().Ethereum().Return(mockEthereumDA).AnyTimes()

	usecase := NewSignTransactionUseCase(mockVault)

	t.Run("should execute use case successfully", func(t *testing.T) {
		privKey := "5385714a2f6d69ca034f56a5268833216ffb8fba7229c39569bc4c5f42cde97c"
		mockEthereumDA.EXPECT().FindOne(ctx, address, namespace).Return(privKey, nil)

		signature, err := usecase.Execute(ctx, address, namespace, chainID, tx)

		assert.NoError(t, err)
		assert.Equal(t, "0xd35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e01", signature)
	})

	t.Run("should fail with CryptoOperationError if creation of ECDSA private key fails", func(t *testing.T) {
		mockEthereumDA.EXPECT().FindOne(ctx, address, namespace).Return("invalidPrivKey", nil)

		signature, err := usecase.Execute(ctx, address, namespace, chainID, tx)

		assert.Empty(t, signature)
		assert.True(t, errors.IsCryptoOperationError(err))
	})

	t.Run("should fail with same error if FindOne fails", func(t *testing.T) {
		expectedErr := errors.HashicorpVaultConnectionError("error")

		mockEthereumDA.EXPECT().FindOne(ctx, gomock.Any(), gomock.Any()).Return("", expectedErr)

		signature, err := usecase.Execute(ctx, address, namespace, chainID, tx)

		assert.Empty(t, signature)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(createAccountComponent), err)
	})
}
