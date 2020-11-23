// +build unit

package ethereum

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store/mocks"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

func TestSignEEATransaction_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVault := mocks.NewMockVault(ctrl)
	mockEthereumDA := mocks.NewMockEthereumAgent(ctrl)
	ctx := context.Background()
	address := "0xaddress"
	namespace := "namespace"
	chainID := "1"
	tx := ethtypes.NewTransaction(
		0,
		common.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"),
		big.NewInt(10000000000),
		21000,
		big.NewInt(10000000000),
		[]byte{},
	)

	mockVault.EXPECT().Ethereum().Return(mockEthereumDA).AnyTimes()

	usecase := NewSignEEATransactionUseCase(mockVault)

	t.Run("should execute use case successfully", func(t *testing.T) {
		privKey := "5385714a2f6d69ca034f56a5268833216ffb8fba7229c39569bc4c5f42cde97c"
		mockEthereumDA.EXPECT().FindOne(ctx, address, namespace).Return(privKey, nil)
		privateArgs := testutils.FakePrivateETHTransactionParams()

		signature, err := usecase.Execute(ctx, address, namespace, chainID, tx, privateArgs)

		assert.NoError(t, err)
		assert.Equal(t, "0x2424ed4546e2039c9f132222eb361286a485a5e9eade6fc5ee1c9548d5391e146d7e794a3a5aa7b553d3905f65824633aab893985112f1c48e6f51e0a8ceb02001", signature)
	})

	t.Run("should fail with same error if FindOne fails", func(t *testing.T) {
		privateArgs := testutils.FakePrivateETHTransactionParams()
		expectedErr := errors.HashicorpVaultConnectionError("error")

		mockEthereumDA.EXPECT().FindOne(ctx, gomock.Any(), gomock.Any()).Return("", expectedErr)

		signature, err := usecase.Execute(ctx, address, namespace, chainID, tx, privateArgs)

		assert.Empty(t, signature)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(signQuorumPrivateTransactionComponent), err)
	})

	t.Run("should fail with CryptoOperationError if creation of ECDSA private key fails", func(t *testing.T) {
		privateArgs := testutils.FakePrivateETHTransactionParams()
		mockEthereumDA.EXPECT().FindOne(ctx, address, namespace).Return("invalidPrivKey", nil)

		signature, err := usecase.Execute(ctx, address, namespace, chainID, tx, privateArgs)

		assert.Empty(t, signature)
		assert.True(t, errors.IsCryptoOperationError(err))
	})

	t.Run("should fail with InvalidParameterError if privateFrom is invalid base64", func(t *testing.T) {
		privateArgs := testutils.FakePrivateETHTransactionParams()
		privKey := "5385714a2f6d69ca034f56a5268833216ffb8fba7229c39569bc4c5f42cde97c"
		mockEthereumDA.EXPECT().FindOne(ctx, address, namespace).Return(privKey, nil)

		privateArgs.PrivateFrom = "invalid privateFrom"
		signature, err := usecase.Execute(ctx, address, namespace, chainID, tx, privateArgs)

		assert.Empty(t, signature)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with InvalidParameterError if privateFor contains invalid base64", func(t *testing.T) {
		privateArgs := testutils.FakePrivateETHTransactionParams()
		privKey := "5385714a2f6d69ca034f56a5268833216ffb8fba7229c39569bc4c5f42cde97c"
		mockEthereumDA.EXPECT().FindOne(ctx, address, namespace).Return(privKey, nil)

		privateArgs.PrivateFor = []string{"invalid privateFor"}
		signature, err := usecase.Execute(ctx, address, namespace, chainID, tx, privateArgs)

		assert.Empty(t, signature)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with InvalidParameterError if privacyGroupID is invalid base64", func(t *testing.T) {
		privateArgs := testutils.FakePrivateETHTransactionParams()
		privKey := "5385714a2f6d69ca034f56a5268833216ffb8fba7229c39569bc4c5f42cde97c"
		mockEthereumDA.EXPECT().FindOne(ctx, address, namespace).Return(privKey, nil)

		privateArgs.PrivacyGroupID = "invalid privacyGroupID"
		signature, err := usecase.Execute(ctx, address, namespace, chainID, tx, privateArgs)

		assert.Empty(t, signature)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}
