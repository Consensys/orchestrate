// +build unit

package ethereum

import (
	"context"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client/mock"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignTransaction_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKeyManagerClient := mock.NewMockKeyManagerClient(ctrl)
	ctx := context.Background()

	usecase := NewSignTransactionUseCase(mockKeyManagerClient)

	t.Run("should execute use case successfully for deployment transactions", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.To = ""
		signature := "0x9a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d5265bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b00"
		nonce, _ := strconv.ParseUint(job.Transaction.Nonce, 10, 64)
		gasLimit, _ := strconv.ParseUint(job.Transaction.Gas, 10, 64)
		expectedRequest := &ethereum.SignETHTransactionRequest{
			Namespace: job.TenantID,
			Nonce:     nonce,
			Amount:    job.Transaction.Value,
			GasPrice:  job.Transaction.GasPrice,
			GasLimit:  gasLimit,
			Data:      job.Transaction.Data,
			To:        job.Transaction.To,
			ChainID:   job.InternalData.ChainID,
		}
		mockKeyManagerClient.EXPECT().ETHSignTransaction(ctx, job.Transaction.From, expectedRequest).Return(signature, nil)

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.Equal(t, "0xf851018227108252088082c35080820713a09a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d52a065bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b", raw)
		assert.Equal(t, "0xbd3684c53d31580ab9ed1a1f133611effbe23f4544ff184886de4749fbc66644", txHash)
	})

	t.Run("should execute use case successfully for contract transactions", func(t *testing.T) {
		job := testutils.FakeJob()
		signature := "0x9a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d5265bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b00"
		nonce, _ := strconv.ParseUint(job.Transaction.Nonce, 10, 64)
		gasLimit, _ := strconv.ParseUint(job.Transaction.Gas, 10, 64)
		expectedRequest := &ethereum.SignETHTransactionRequest{
			Namespace: job.TenantID,
			Nonce:     nonce,
			Amount:    job.Transaction.Value,
			GasPrice:  job.Transaction.GasPrice,
			GasLimit:  gasLimit,
			Data:      job.Transaction.Data,
			To:        job.Transaction.To,
			ChainID:   job.InternalData.ChainID,
		}
		mockKeyManagerClient.EXPECT().ETHSignTransaction(ctx, job.Transaction.From, expectedRequest).Return(signature, nil)

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.Equal(t, "0xf86501822710825208944fed1fc4144c223ae3c1553be203cdfcbd38c58182c35080820713a09a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d52a065bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b", raw)
		assert.Equal(t, "0x5b8dc6b3ce158ca2d4bafdf171f182f8fe451baac81c3fc8703dbf1d237ea2d4", txHash)
	})

	t.Run("should fail with same error if ETHSignTransaction fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		mockKeyManagerClient.EXPECT().ETHSignTransaction(ctx, gomock.Any(), gomock.Any()).Return("", expectedErr)

		raw, txHash, err := usecase.Execute(ctx, testutils.FakeJob())

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(signTransactionComponent), err)
		assert.Empty(t, raw)
		assert.Empty(t, txHash)
	})

	t.Run("should fail with EncodingError if signature cannot be decoded", func(t *testing.T) {
		signature := "invalidSignature"
		mockKeyManagerClient.EXPECT().ETHSignTransaction(ctx, gomock.Any(), gomock.Any()).Return(signature, nil)

		raw, txHash, err := usecase.Execute(ctx, testutils.FakeJob())

		assert.True(t, errors.IsEncodingError(err))
		assert.Empty(t, raw)
		assert.Empty(t, txHash)
	})
}
