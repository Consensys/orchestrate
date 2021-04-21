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

	publicKey := "0x5fd633ff9f8ee36f9e3a874709406103854c0f6650cb908c010ea55eabc35191"
	payload := "0xda"
	signature := "0xbdb22e268e765473720646ade5df2b35df131e97f1ea85a98cb3c2b88858f79c02d3ca05397c2a65f38063e55aad69d9479838605559f67f9fb75860df766497"

	usecase := NewVerifySignatureUseCase()

	t.Run("should execute use case successfully", func(t *testing.T) {
		err := usecase.Execute(ctx, publicKey, signature, payload)

		assert.NoError(t, err)
	})

	t.Run("should fail with InvalidParameterError if data not a hex string", func(t *testing.T) {
		err := usecase.Execute(ctx, publicKey, signature, "invalid data")
		assert.True(t, errors.IsInvalidParameterError(err))
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
