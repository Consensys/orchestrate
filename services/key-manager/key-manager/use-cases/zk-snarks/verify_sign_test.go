package zksnarks

import (
	"context"
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestVerifySignature_Execute(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	publicKey := "0xefaecb7b08beec692cc136ad5ad7c249c337ae0890e88c6afd2f67ad51d1ad15"
	payload := "data to sign"
	signature := "0x68eb9b75aa8a0ae94a130fa2da013f281ccbecfea2572fc597451fb80b8acc92008c0d011e3f7524e6f3d98b528ffe26d984f9d154ef14d71d2899cca2101705"

	usecase := NewVerifySignatureUseCase()

	t.Run("should execute use case successfully", func(t *testing.T) {
		err := usecase.Execute(ctx, publicKey, signature, payload)

		assert.NoError(t, err)
	})

	t.Run("should fail with InvalidParameterError if fails to decode signature", func(t *testing.T) {
		err := usecase.Execute(ctx, publicKey, "invalid signature", payload)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with InvalidParameterError if fails to decode public key", func(t *testing.T) {
		err := usecase.Execute(ctx, "invalidPublicKey", signature, payload)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with InvalidParameterError if payload does not match signature", func(t *testing.T) {
		err := usecase.Execute(ctx, publicKey, signature, "invalid payload")
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}
