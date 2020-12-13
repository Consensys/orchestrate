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

func TestSignEEATransaction_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKeyManagerClient := mock.NewMockKeyManagerClient(ctrl)
	ctx := context.Background()

	usecase := NewSignEEATransactionUseCase(mockKeyManagerClient)

	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJob()
		signature := "0x9a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d5265bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b00"
		nonce, _ := strconv.ParseUint(job.Transaction.Nonce, 10, 64)
		expectedRequest0 := &ethereum.SignEEATransactionRequest{
			Namespace:      multitenancy.DefaultTenant,
			Nonce:          nonce,
			To:             job.Transaction.To,
			Data:           job.Transaction.Data,
			ChainID:        job.InternalData.ChainID,
			PrivateFrom:    job.Transaction.PrivateFrom,
			PrivateFor:     job.Transaction.PrivateFor,
			PrivacyGroupID: job.Transaction.PrivacyGroupID,
		}

		expectedRequest1 := &ethereum.SignEEATransactionRequest{
			Namespace:      job.TenantID,
			Nonce:          nonce,
			To:             job.Transaction.To,
			Data:           job.Transaction.Data,
			ChainID:        job.InternalData.ChainID,
			PrivateFrom:    job.Transaction.PrivateFrom,
			PrivateFor:     job.Transaction.PrivateFor,
			PrivacyGroupID: job.Transaction.PrivacyGroupID,
		}

		gomock.InOrder(
			mockKeyManagerClient.EXPECT().ETHSignEEATransaction(ctx, job.Transaction.From, expectedRequest0).Return("", errors.NotFoundError("not found")),
			mockKeyManagerClient.EXPECT().ETHSignEEATransaction(ctx, job.Transaction.From, expectedRequest1).Return(signature, nil),
		)

		raw, txHash, err := usecase.Execute(ctx, job)

		assert.NoError(t, err)
		assert.Equal(t, "0xf8d501822710825208944fed1fc4144c223ae3c1553be203cdfcbd38c58182c35080820713a09a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d52a065bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8ba0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486af842a0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486aa0075695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486a8a72657374726963746564", raw)
		assert.Empty(t, txHash)
	})

	t.Run("should execute use case successfully for deployment transactions", func(t *testing.T) {
		job := testutils.FakeJob()
		job.Transaction.To = ""
		signature := "0x9a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d5265bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b00"
		nonce, _ := strconv.ParseUint(job.Transaction.Nonce, 10, 64)
		expectedRequest := &ethereum.SignEEATransactionRequest{
			Namespace:      multitenancy.DefaultTenant,
			Nonce:          nonce,
			Data:           job.Transaction.Data,
			ChainID:        job.InternalData.ChainID,
			PrivateFrom:    job.Transaction.PrivateFrom,
			PrivateFor:     job.Transaction.PrivateFor,
			PrivacyGroupID: job.Transaction.PrivacyGroupID,
		}
		mockKeyManagerClient.EXPECT().ETHSignEEATransaction(ctx, job.Transaction.From, expectedRequest).Return(signature, nil)

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
		expectedErr := errors.InvalidFormatError("error")
		mockKeyManagerClient.EXPECT().ETHSignEEATransaction(ctx, gomock.Any(), gomock.Any()).Return("", expectedErr)

		raw, txHash, err := usecase.Execute(ctx, testutils.FakeJob())

		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(signEEATransactionComponent), err)
		assert.Empty(t, raw)
		assert.Empty(t, txHash)
	})

	t.Run("should fail with EncodingError if signature cannot be decoded", func(t *testing.T) {
		signature := "invalidSignature"
		mockKeyManagerClient.EXPECT().ETHSignEEATransaction(ctx, gomock.Any(), gomock.Any()).Return(signature, nil)

		raw, txHash, err := usecase.Execute(ctx, testutils.FakeJob())

		assert.True(t, errors.IsEncodingError(err))
		assert.Empty(t, raw)
		assert.Empty(t, txHash)
	})

	t.Run("should fail with InvalidParameterError if ETHSignEEATransaction fails to find tenant", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		mockKeyManagerClient.EXPECT().ETHSignEEATransaction(ctx, gomock.Any(), gomock.Any()).Return("", expectedErr).Times(2)

		raw, txHash, err := usecase.Execute(ctx, testutils.FakeJob())

		assert.True(t, errors.IsInvalidParameterError(err))
		assert.Empty(t, raw)
		assert.Empty(t, txHash)
	})
}
