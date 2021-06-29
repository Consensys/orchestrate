// +build unit

package signer

import (
	"context"
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/multitenancy"
	qkm "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager"
	qkmmock "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager/client/mocks"
	"github.com/ConsenSys/orchestrate/pkg/quorum-key-manager/types"
	"github.com/ConsenSys/orchestrate/pkg/types/testutils"
	"github.com/consensys/quorum/common/hexutil"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignEEATransaction_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	globalStoreName := "test-store-name"
	qkm.SetGlobalStoreName(globalStoreName)
	mockKeyManagerClient := qkmmock.NewMockKeyManagerClient(ctrl)
	ctx := context.Background()

	usecase := NewSignEEATransactionUseCase(mockKeyManagerClient)

	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		signature := "0x9a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d5265bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b00"
		acc := qkm.FakeEth1AccountResponse(job.Transaction.From, []string{job.TenantID})
		txData, _ := hexutil.Decode("0xbc00fb0501bfc05a2928a14848864983d89945d1088e67bb264ec308c58c409f")
		expectedRequest := &types.SignHexPayloadRequest{
			Data: txData,
		}
		mockKeyManagerClient.EXPECT().GetEth1Account(gomock.Any(), globalStoreName, job.Transaction.From).Return(acc, nil)
		mockKeyManagerClient.EXPECT().SignEth1Data(gomock.Any(), globalStoreName, job.Transaction.From, expectedRequest).Return(signature, nil)

		raw, txHash, err := usecase.Execute(ctx, job)

		require.NoError(t, err)
		assert.Equal(t, "0xf8d501822710825208944fed1fc4144c223ae3c1553be203cdfcbd38c58182c35080820713a09a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d52a065bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8ba0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486af842a0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486aa0075695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486a8a72657374726963746564", raw)
		assert.Empty(t, txHash)
	})

	t.Run("should execute use case successfully for deployment transactions", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.To = ""
		signature := "0x9a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d5265bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b00"
		acc := qkm.FakeEth1AccountResponse(job.Transaction.From, []string{job.TenantID})
		txData, _ := hexutil.Decode("0xd3b00b8f9fa470236ab7cfc080b4b790f884705740be763c3532134253b4a8ef")
		expectedRequest := &types.SignHexPayloadRequest{
			Data: txData,
		}
		mockKeyManagerClient.EXPECT().GetEth1Account(gomock.Any(), globalStoreName, job.Transaction.From).Return(acc, nil)
		mockKeyManagerClient.EXPECT().SignEth1Data(gomock.Any(), globalStoreName, job.Transaction.From, expectedRequest).Return(signature, nil)

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.Equal(t, "0xf8c1018227108252088082c35080820713a09a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d52a065bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8ba0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486af842a0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486aa0075695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486a8a72657374726963746564", raw)
		assert.Empty(t, txHash)
	})

	t.Run("should execute use case successfully for one time key transactions", func(t *testing.T) {
		job := testutils.FakeJob()
		job.InternalData.OneTimeKey = true
		job.Transaction.PrivateFrom = "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="
		job.Transaction.PrivateFor = []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="}

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.NotEmpty(t, raw)
		assert.Empty(t, txHash)
	})

	t.Run("should fail with same error if ETHSignEEATransaction fails", func(t *testing.T) {
		job := testutils.FakeJob()
		acc := qkm.FakeEth1AccountResponse(job.Transaction.From, []string{multitenancy.DefaultTenant})
		
		mockKeyManagerClient.EXPECT().GetEth1Account(gomock.Any(), globalStoreName, job.Transaction.From).Return(acc, nil)
		mockKeyManagerClient.EXPECT().SignEth1Data(gomock.Any(), globalStoreName, job.Transaction.From, gomock.Any()).Return("", errors.InvalidFormatError("error"))

		raw, txHash, err := usecase.Execute(ctx, testutils.FakeJob())

		assert.True(t, errors.IsDependencyFailureError(err))
		assert.Empty(t, raw)
		assert.Empty(t, txHash)
	})

	t.Run("should fail with EncodingError if signature cannot be decoded", func(t *testing.T) {
		job := testutils.FakeJob()
		signature := "invalidSignature"
		acc := qkm.FakeEth1AccountResponse(job.Transaction.From, []string{multitenancy.DefaultTenant})
		mockKeyManagerClient.EXPECT().GetEth1Account(gomock.Any(), globalStoreName, job.Transaction.From).Return(acc, nil)
		mockKeyManagerClient.EXPECT().SignEth1Data(gomock.Any(), globalStoreName, gomock.Any(), gomock.Any()).Return(signature, nil)

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.True(t, errors.IsEncodingError(err))
		assert.Empty(t, raw)
		assert.Empty(t, txHash)
	})

	t.Run("should fail with InvalidAuthenticationError if ETHSignEEATransaction fails to find tenant", func(t *testing.T) {
		job := testutils.FakeJob()
		acc := qkm.FakeEth1AccountResponse(job.Transaction.From, []string{})

		mockKeyManagerClient.EXPECT().GetEth1Account(gomock.Any(), globalStoreName, job.Transaction.From).Return(acc, nil)

		_, _, err := usecase.Execute(ctx, job)

		assert.True(t, errors.IsInvalidAuthenticationError(err))
	})
	
	t.Run("should fail with IsDependencyFailureError if fails to find account", func(t *testing.T) {
		job := testutils.FakeJob()

		mockKeyManagerClient.EXPECT().GetEth1Account(gomock.Any(), globalStoreName, job.Transaction.From).
			Return(nil,  errors.NotFoundError("account no found"))

		_, _, err := usecase.Execute(ctx, job)

		assert.True(t, errors.IsDependencyFailureError(err))
	})
}
