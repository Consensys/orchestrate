// +build unit

package signer

import (
	"context"
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	qkm "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager"
	qkmmock "github.com/consensys/quorum-key-manager/pkg/client/mock"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/ConsenSys/orchestrate/pkg/types/testutils"
	"github.com/consensys/quorum/common/hexutil"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
		signature := "0x9a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d5265bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b00"
		acc := qkm.FakeEth1AccountResponse(job.Transaction.From, []string{job.TenantID})
		txData, _ := hexutil.Decode("0xde9325591179dad0b11f84231b5b5dcfce57f4b16247ab76010cef71a4e6aa0d")
		expectedRequest := &types.SignHexPayloadRequest{
			Data: txData,
		}
		mockKeyManagerClient.EXPECT().GetEth1Account(gomock.Any(), globalStoreName, job.Transaction.From).Return(acc, nil)
		mockKeyManagerClient.EXPECT().SignEth1Data(gomock.Any(), globalStoreName, job.Transaction.From, expectedRequest).Return(signature, nil)

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.Equal(t, "0xf86301822710825208944fed1fc4144c223ae3c1553be203cdfcbd38c58182c3508025a09a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d52a065bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b", raw)
		assert.Equal(t, "0x087251e41ae0f2d339053abbd8bef70f2b6d0c18d092b0df11f990ba95526f43", txHash)
	})

	t.Run("should execute use case successfully for deployment transactions", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.To = ""
		signature := "0x9a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d5265bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b00"
		acc := qkm.FakeEth1AccountResponse(job.Transaction.From, []string{job.TenantID})
		txData, _ := hexutil.Decode("0xa646977b9596105d021f2df1a6f4cdd83cb32c795e91c1cfc46197075b0a3009")
		expectedRequest := &types.SignHexPayloadRequest{
			Data: txData,
		}
		mockKeyManagerClient.EXPECT().GetEth1Account(gomock.Any(), globalStoreName, job.Transaction.From).Return(acc, nil)
		mockKeyManagerClient.EXPECT().SignEth1Data(gomock.Any(), globalStoreName, job.Transaction.From, expectedRequest).Return(signature, nil)

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.Equal(t, "0xf84f018227108252088082c3508025a09a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d52a065bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b", raw)
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
		acc := qkm.FakeEth1AccountResponse(job.Transaction.From, []string{job.TenantID})
		mockKeyManagerClient.EXPECT().GetEth1Account(gomock.Any(), globalStoreName, job.Transaction.From).Return(acc, nil)
		mockKeyManagerClient.EXPECT().SignEth1Data(gomock.Any(), globalStoreName, gomock.Any(), gomock.Any()).Return("", expectedErr)

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.True(t, errors.IsDependencyFailureError(err))
		assert.Empty(t, raw)
		assert.Empty(t, txHash)
	})

	t.Run("should fail with EncodingError if signature cannot be decoded", func(t *testing.T) {
		signature := "invalidSignature"
		job := testutils.FakeJob()
		acc := qkm.FakeEth1AccountResponse(job.Transaction.From, []string{job.TenantID})
		mockKeyManagerClient.EXPECT().GetEth1Account(gomock.Any(), globalStoreName, job.Transaction.From).Return(acc, nil)
		mockKeyManagerClient.EXPECT().SignEth1Data(gomock.Any(), globalStoreName, gomock.Any(), gomock.Any()).Return(signature, nil)

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.True(t, errors.IsEncodingError(err))
		assert.Empty(t, raw)
		assert.Empty(t, txHash)
	})

	t.Run("should fail with IsDependencyFailureError if fails to find account", func(t *testing.T) {
		job := testutils.FakeJob()

		mockKeyManagerClient.EXPECT().GetEth1Account(gomock.Any(), globalStoreName, job.Transaction.From).
			Return(nil,  errors.NotFoundError("account no found"))

		_, _, err := usecase.Execute(ctx, job)

		assert.True(t, errors.IsDependencyFailureError(err))
	})

	t.Run("should fail with IsInvalidAuthenticationError if tenant is not allowed to access account", func(t *testing.T) {
		job := testutils.FakeJob()
		acc := qkm.FakeEth1AccountResponse(job.Transaction.From, []string{})

		mockKeyManagerClient.EXPECT().GetEth1Account(gomock.Any(), globalStoreName, job.Transaction.From).Return(acc, nil)

		_, _, err := usecase.Execute(ctx, testutils.FakeJob())

		assert.True(t, errors.IsInvalidAuthenticationError(err))
	})
}
