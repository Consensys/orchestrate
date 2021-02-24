// +build unit

package ethereum

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/types/testutils"
	"github.com/ConsenSys/orchestrate/services/key-manager/service/formatters"
	"github.com/ConsenSys/orchestrate/services/key-manager/store/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestSignTypedData_Execute(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVault := mocks.NewMockVault(ctrl)
	address := "0xaddress"
	namespace := "namespace"

	usecase := NewSignTypedDataUseCase(mockVault)

	t.Run("should execute use case successfully", func(t *testing.T) {
		expectedSignature := "0xsignature"
		expectedHashDomainSeparator := "0x19412cf2425f8663d8b4be7f725bde1ffad8fdff6ed2d9cf0084d785a715912a"
		expectedHashMessage := "0x93d8044d5438b567977f89a53282195effc64659ef5a29e60960e5727dd8580c"
		expectedEncodedData := fmt.Sprintf("\x19\x01%s%s", expectedHashDomainSeparator, expectedHashMessage)

		mockVault.EXPECT().ETHSign(address, namespace, expectedEncodedData).Return(expectedSignature, nil)

		typedData := formatters.FormatSignTypedDataRequest(testutils.FakeSignTypedDataRequest())
		signature, err := usecase.Execute(ctx, address, namespace, typedData)

		assert.NoError(t, err)
		assert.Equal(t, expectedSignature, signature)
	})

	t.Run("should fail with InvalidParameterError if fails to hash data", func(t *testing.T) {
		typedData := formatters.FormatSignTypedDataRequest(testutils.FakeSignTypedDataRequest())
		typedData.Message = map[string]interface{}{
			"invalid": "invalid",
		}
		signature, err := usecase.Execute(ctx, address, namespace, typedData)

		assert.Empty(t, signature)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	t.Run("should fail with InvalidParameterError if fails to hash domain separator", func(t *testing.T) {
		typedData := formatters.FormatSignTypedDataRequest(testutils.FakeSignTypedDataRequest())
		typedData.Domain.Version = ""
		signature, err := usecase.Execute(ctx, address, namespace, typedData)

		assert.Empty(t, signature)
		assert.True(t, errors.IsInvalidParameterError(err))

	})

	t.Run("should fail with same error if ETHSign fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")
		mockVault.EXPECT().ETHSign(address, namespace, gomock.Any()).Return("", expectedErr)

		typedData := formatters.FormatSignTypedDataRequest(testutils.FakeSignTypedDataRequest())
		signature, err := usecase.Execute(ctx, address, namespace, typedData)

		assert.Empty(t, signature)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(signTypedDataComponent), err)
	})
}
