// +build unit

package signer

import (
	"context"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client/mock"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignQuorumPrivateTransaction_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKeyManagerClient := mock.NewMockKeyManagerClient(ctrl)
	ctx := context.Background()

	usecase := NewSignQuorumPrivateTransactionUseCase(mockKeyManagerClient)

	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		signature := "0x9a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d5265bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b00"
		nonce, _ := strconv.ParseUint(job.Transaction.Nonce, 10, 64)
		gasLimit, _ := strconv.ParseUint(job.Transaction.Gas, 10, 64)
		expectedRequest := &ethereum.SignQuorumPrivateTransactionRequest{
			Namespace: multitenancy.DefaultTenant,
			Nonce:     nonce,
			To:        job.Transaction.To,
			Amount:    job.Transaction.Value,
			GasPrice:  job.Transaction.GasPrice,
			GasLimit:  gasLimit,
			Data:      job.Transaction.Data,
		}
		mockKeyManagerClient.EXPECT().ETHSignQuorumPrivateTransaction(ctx, job.Transaction.From, expectedRequest).Return(signature, nil)

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.Equal(t, "0xf86301822710825208944fed1fc4144c223ae3c1553be203cdfcbd38c58182c3508025a09a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d52a065bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b", raw)
		assert.Equal(t, "0x087251e41ae0f2d339053abbd8bef70f2b6d0c18d092b0df11f990ba95526f43", txHash)
	})

	t.Run("should execute use case successfully for deployment transactions", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.To = ""
		signature := "0x9a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d5265bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b00"
		nonce, _ := strconv.ParseUint(job.Transaction.Nonce, 10, 64)
		gasLimit, _ := strconv.ParseUint(job.Transaction.Gas, 10, 64)
		expectedRequest := &ethereum.SignQuorumPrivateTransactionRequest{
			Namespace: multitenancy.DefaultTenant,
			Nonce:     nonce,
			Amount:    job.Transaction.Value,
			GasPrice:  job.Transaction.GasPrice,
			GasLimit:  gasLimit,
			Data:      job.Transaction.Data,
		}
		mockKeyManagerClient.EXPECT().ETHSignQuorumPrivateTransaction(ctx, job.Transaction.From, expectedRequest).Return(signature, nil)

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
		mockKeyManagerClient.EXPECT().ETHSignQuorumPrivateTransaction(ctx, gomock.Any(), gomock.Any()).Return("", expectedErr)

		raw, txHash, err := usecase.Execute(ctx, testutils.FakeJob())

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(signQuorumPrivateTransactionComponent), err)
		assert.Empty(t, raw)
		assert.Empty(t, txHash)
	})

	t.Run("should fail with EncodingError if signature cannot be decoded", func(t *testing.T) {
		signature := "invalidSignature"
		mockKeyManagerClient.EXPECT().ETHSignQuorumPrivateTransaction(ctx, gomock.Any(), gomock.Any()).Return(signature, nil)

		raw, txHash, err := usecase.Execute(ctx, testutils.FakeJob())

		assert.True(t, errors.IsEncodingError(err))
		assert.Empty(t, raw)
		assert.Empty(t, txHash)
	})

	t.Run("should fail with InvalidParameterError if ETHSignQuorumPrivateTransaction fails to find tenant", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		mockKeyManagerClient.EXPECT().ETHSignQuorumPrivateTransaction(ctx, gomock.Any(), gomock.Any()).Return("", expectedErr).Times(2)

		raw, txHash, err := usecase.Execute(ctx, testutils.FakeJob())

		assert.True(t, errors.IsInvalidParameterError(err))
		assert.Empty(t, raw)
		assert.Empty(t, txHash)
	})
}
