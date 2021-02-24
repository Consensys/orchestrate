// +build unit

package ethereum

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestVerifySignature_Execute(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	address := "0x5Cc634233E4a454d47aACd9fC68801482Fb02610"
	payload := "my data to sign"

	usecase := NewVerifySignatureUseCase()

	t.Run("should execute use case successfully", func(t *testing.T) {
		signature := "0x34334af7bacf5d82bb892c838beda65331232c29e122b3485f31e14eda731dbb0ebae9d1eed72c099ff4c3b462aebf449068f717f3638a6facd0b3dddf2529a500"
		err := usecase.Execute(ctx, address, signature, payload)
		assert.NoError(t, err)
	})

	t.Run("should fail with InvalidParameterError if fails to decode signature", func(t *testing.T) {
		invalidSignature := "invalid"
		err := usecase.Execute(ctx, address, invalidSignature, payload)

		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with InvalidParameterError if fails to recover public key (invalid signature length)", func(t *testing.T) {
		invalidSignature := "0x34334af7bacf5d82bb"

		err := usecase.Execute(ctx, address, invalidSignature, payload)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with InvalidParameterError if addresses do not match", func(t *testing.T) {
		signature := "0x3399aeb23d6564b3a0b220447e9f1bb2057ffb82cfb766147620aa6bc84938e26941e7583d6460fea405d99da897e88cab07a7fd0991c6c2163645c45d25e4b201"
		err := usecase.Execute(ctx, address, signature, payload)

		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with InvalidParameterError if payload does not match", func(t *testing.T) {
		signature := "0x34334af7bacf5d82bb892c838beda65331232c29e122b3485f31e14eda731dbb0ebae9d1eed72c099ff4c3b462aebf449068f717f3638a6facd0b3dddf2529a500"
		err := usecase.Execute(ctx, address, signature, "my data that was changed")

		assert.True(t, errors.IsInvalidParameterError(err))
	})
}
