// +build unit

package signer

import (
	"context"
	"testing"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/pkg/utils"
	qkmmock "github.com/consensys/quorum-key-manager/pkg/client/mock"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignTransaction_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKeyManagerClient := qkmmock.NewMockKeyManagerClient(ctrl)
	ctx := context.Background()

	usecase := NewSignETHTransactionUseCase(mockKeyManagerClient)

	signedRaw := utils.StringToHexBytes("0xf86501822710825208944fed1fc4144c223ae3c1553be203cdfcbd38c58182c35080820713a09a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d52a065bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b")
	
	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		mockKeyManagerClient.EXPECT().SignTransaction(gomock.Any(), job.InternalData.StoreID, job.Transaction.From.String(), 
			gomock.AssignableToTypeOf(&types.SignETHTransactionRequest{})).Return(signedRaw.String(), nil)

		raw, txHash, err := usecase.Execute(ctx, job)

		require.NoError(t, err)
		assert.Equal(t, signedRaw.String(), raw.String())
		assert.Equal(t, "0xbb07e6a2f123a19d98b890eab6cb3947c0b55786b98bd09a412496d8d09cabfb", txHash.String())
	})

	t.Run("should execute use case successfully for deployment transactions", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.To = nil
		mockKeyManagerClient.EXPECT().SignTransaction(gomock.Any(), job.InternalData.StoreID, job.Transaction.From.String(), 
			gomock.AssignableToTypeOf(&types.SignETHTransactionRequest{})).Return(signedRaw.String(), nil)

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.Equal(t, signedRaw.String(), raw.String())
		assert.Equal(t, "0xbb07e6a2f123a19d98b890eab6cb3947c0b55786b98bd09a412496d8d09cabfb", txHash.String())
	})

	t.Run("should execute use case successfully for one time key transactions", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.TransactionType = entities.LegacyTxType
		job.InternalData.OneTimeKey = true

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.NotEmpty(t, raw)
		assert.NotEmpty(t, txHash)
	})

	t.Run("should fail with same error if ETHSignTransaction fails", func(t *testing.T) {
		expectedErr := errors.InvalidFormatError("error")
		mockKeyManagerClient.EXPECT().SignTransaction(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return("", expectedErr)

		raw, txHash, err := usecase.Execute(ctx, testutils.FakeJob())

		assert.True(t, errors.IsDependencyFailureError(err))
		assert.Empty(t, raw)
		assert.Empty(t, txHash)
	})
}
