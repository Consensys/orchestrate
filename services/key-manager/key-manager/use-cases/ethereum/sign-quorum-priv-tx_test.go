// +build unit

package ethereum

import (
	"context"
	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store/mocks"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

func TestSignQuorumPrivateTransaction_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVault := mocks.NewMockVault(ctrl)
	mockEthereumDA := mocks.NewMockEthereumAgent(ctrl)
	ctx := context.Background()
	address := "0xaddress"
	namespace := "namespace"
	tx := quorumtypes.NewTransaction(
		0,
		common.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"),
		big.NewInt(10000000000),
		21000,
		big.NewInt(10000000000),
		[]byte{},
	)

	mockVault.EXPECT().Ethereum().Return(mockEthereumDA).AnyTimes()

	usecase := NewSignQuorumPrivateTransactionUseCase(mockVault)

	t.Run("should execute use case successfully", func(t *testing.T) {
		privKey := "5385714a2f6d69ca034f56a5268833216ffb8fba7229c39569bc4c5f42cde97c"
		mockEthereumDA.EXPECT().FindOne(ctx, address, namespace).Return(privKey, nil)

		signature, err := usecase.Execute(ctx, address, namespace, tx)

		assert.NoError(t, err)
		assert.Equal(t, "0xefa9c4498397ee12e341f6acf81072bbf0c8fb4e4e1813ac96fd3860baa28bb931aecd59811beaffc71a4ef008882d3c13537a2f733be7643fdfea4ea77f3ded00", signature)
	})

	t.Run("should fail with CryptoOperationError if creation of ECDSA private key fails", func(t *testing.T) {
		mockEthereumDA.EXPECT().FindOne(ctx, address, namespace).Return("invalidPrivKey", nil)

		signature, err := usecase.Execute(ctx, address, namespace, tx)

		assert.Empty(t, signature)
		assert.True(t, errors.IsCryptoOperationError(err))
	})

	t.Run("should fail with same error if FindOne fails", func(t *testing.T) {
		expectedErr := errors.HashicorpVaultConnectionError("error")

		mockEthereumDA.EXPECT().FindOne(ctx, gomock.Any(), gomock.Any()).Return("", expectedErr)

		signature, err := usecase.Execute(ctx, address, namespace, tx)

		assert.Empty(t, signature)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(signQuorumPrivateTransactionComponent), err)
	})
}
