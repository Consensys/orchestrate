// +build unit

package signer

import (
	"context"
	"testing"

	"github.com/consensys/orchestrate/pkg/errors"
	qkm "github.com/consensys/orchestrate/pkg/quorum-key-manager"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	qkmmock "github.com/consensys/quorum-key-manager/pkg/client/mock"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignQuorumPrivateTransaction_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	globalStoreName := "test-store-name"
	qkm.SetGlobalStoreName(globalStoreName)
	mockKeyManagerClient := qkmmock.NewMockKeyManagerClient(ctrl)
	ctx := context.Background()

	usecase := NewSignQuorumPrivateTransactionUseCase(mockKeyManagerClient)

	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		signedRaw := "0xf851018227108252088082c35080820713a09a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d52a065bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b"
		mockKeyManagerClient.EXPECT().SignQuorumPrivateTransaction(gomock.Any(), globalStoreName, job.Transaction.From, gomock.AssignableToTypeOf(&types.SignQuorumPrivateTransactionRequest{})).Return(signedRaw, nil)

		raw, txHash, err := usecase.Execute(ctx, job)

		require.NoError(t, err)
		assert.Equal(t, signedRaw, raw)
		assert.Equal(t, "0x2c8bdef96dca7d037618fcf799a4bbfec6a6e1299b27dbc5d7cd79594e31ee54", txHash)
	})

	t.Run("should execute use case successfully for deployment transactions", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.To = ""
		signedRaw := "0xf84f018227108252088082c3508025a09a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d52a065bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b"
		mockKeyManagerClient.EXPECT().SignQuorumPrivateTransaction(gomock.Any(), globalStoreName, job.Transaction.From, gomock.AssignableToTypeOf(&types.SignQuorumPrivateTransactionRequest{})).Return(signedRaw, nil)

		raw, txHash, err := usecase.Execute(ctx, job)

		require.NoError(t, err)
		assert.Equal(t, signedRaw, raw)
		assert.Equal(t, "0xdcd5ec31fe3201903d2c580d4b2ff3d6c48e2980f77f7410d9868e788c19dacb", txHash)
	})

	t.Run("should execute use case successfully for one time key transactions", func(t *testing.T) {
		job := testutils.FakeJob()
		job.InternalData.OneTimeKey = true

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.NotEmpty(t, raw)
		assert.NotEmpty(t, txHash)
	})

	t.Run("should fail with same error if ETHSignQuorumPrivateTransaction fails", func(t *testing.T) {
		expectedErr := errors.InvalidFormatError("error")
		job := testutils.FakeJob()
		mockKeyManagerClient.EXPECT().SignQuorumPrivateTransaction(gomock.Any(), globalStoreName, gomock.Any(), gomock.Any()).Return("", expectedErr)

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.True(t, errors.IsDependencyFailureError(err))
		assert.Empty(t, raw)
		assert.Empty(t, txHash)
	})
}
