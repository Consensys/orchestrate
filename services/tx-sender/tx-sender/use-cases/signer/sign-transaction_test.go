// +build unit

package signer

import (
	"context"
	"testing"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	qkm "github.com/consensys/orchestrate/pkg/quorum-key-manager"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	qkmmock "github.com/consensys/quorum-key-manager/pkg/client/mock"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignTransaction_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	globalStoreName := "test-store-name"
	qkm.SetGlobalStoreName(globalStoreName)
	mockKeyManagerClient := qkmmock.NewMockKeyManagerClient(ctrl)
	ctx := context.Background()

	usecase := NewSignETHTransactionUseCase(mockKeyManagerClient)

	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		signedRaw := "0xf86501822710825208944fed1fc4144c223ae3c1553be203cdfcbd38c58182c35080820713a09a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d52a065bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b"
		acc := qkm.FakeEthAccountResponse(job.Transaction.From, []string{multitenancy.DefaultTenant})
		mockKeyManagerClient.EXPECT().GetEthAccount(gomock.Any(), globalStoreName, job.Transaction.From).Return(acc, nil)
		mockKeyManagerClient.EXPECT().SignTransaction(gomock.Any(), globalStoreName, job.Transaction.From, gomock.AssignableToTypeOf(&types.SignETHTransactionRequest{})).Return(signedRaw, nil)

		raw, txHash, err := usecase.Execute(ctx, job)

		require.NoError(t, err)
		assert.Equal(t, signedRaw, raw)
		assert.Equal(t, "0xbb07e6a2f123a19d98b890eab6cb3947c0b55786b98bd09a412496d8d09cabfb", txHash)
	})

	t.Run("should execute use case successfully for deployment transactions", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.To = ""
		signedRaw := "0xf851018227108252088082c35080820713a09a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d52a065bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b"
		acc := qkm.FakeEthAccountResponse(job.Transaction.From, []string{job.TenantID})
		mockKeyManagerClient.EXPECT().GetEthAccount(gomock.Any(), globalStoreName, job.Transaction.From).Return(acc, nil)
		mockKeyManagerClient.EXPECT().SignTransaction(gomock.Any(), globalStoreName, job.Transaction.From, gomock.AssignableToTypeOf(&types.SignETHTransactionRequest{})).Return(signedRaw, nil)

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.Equal(t, signedRaw, raw)
		assert.Equal(t, "0x2c8bdef96dca7d037618fcf799a4bbfec6a6e1299b27dbc5d7cd79594e31ee54", txHash)
	})

	t.Run("should execute use case successfully for one time key transactions", func(t *testing.T) {
		job := testutils.FakeJob()
		job.InternalData.OneTimeKey = true

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.NotEmpty(t, raw)
		assert.NotEmpty(t, txHash)
	})

	t.Run("should fail with same error if ETHSignTransaction fails", func(t *testing.T) {
		expectedErr := errors.InvalidFormatError("error")
		job := testutils.FakeJob()
		acc := qkm.FakeEthAccountResponse(job.Transaction.From, []string{multitenancy.DefaultTenant})

		mockKeyManagerClient.EXPECT().GetEthAccount(gomock.Any(), globalStoreName, job.Transaction.From).Return(acc, nil)
		mockKeyManagerClient.EXPECT().SignTransaction(gomock.Any(), globalStoreName, gomock.Any(), gomock.Any()).Return("", expectedErr)

		raw, txHash, err := usecase.Execute(ctx, testutils.FakeJob())

		assert.True(t, errors.IsDependencyFailureError(err))
		assert.Empty(t, raw)
		assert.Empty(t, txHash)
	})

	t.Run("should fail with IsDependencyFailureError if fails to find account", func(t *testing.T) {
		job := testutils.FakeJob()

		mockKeyManagerClient.EXPECT().GetEthAccount(gomock.Any(), globalStoreName, job.Transaction.From).
			Return(nil, errors.NotFoundError("account no found"))

		_, _, err := usecase.Execute(ctx, job)

		assert.True(t, errors.IsDependencyFailureError(err))
	})

	t.Run("should fail with IsInvalidAuthenticationError if tenant is not allowed to access account", func(t *testing.T) {
		job := testutils.FakeJob()
		acc := qkm.FakeEthAccountResponse(job.Transaction.From, []string{})

		mockKeyManagerClient.EXPECT().GetEthAccount(gomock.Any(), globalStoreName, job.Transaction.From).Return(acc, nil)

		_, _, err := usecase.Execute(ctx, testutils.FakeJob())

		assert.True(t, errors.IsInvalidAuthenticationError(err))
	})
}
