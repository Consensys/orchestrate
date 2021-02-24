// +build unit

package ethereum

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/ConsenSys/orchestrate/pkg/types/testutils"
	"github.com/ConsenSys/orchestrate/services/key-manager/key-manager/use-cases/mocks"
	"github.com/ConsenSys/orchestrate/services/key-manager/service/formatters"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestVerifyTypedDataSignature_Execute(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVerifySignatureUC := mocks.NewMockVerifyETHSignatureUseCase(ctrl)
	address := "0x5Cc634233E4a454d47aACd9fC68801482Fb02610"
	payload := formatters.FormatSignTypedDataRequest(testutils.FakeSignTypedDataRequest())
	signature := "0xsignature"

	usecase := NewVerifyTypedDataSignatureUseCase(mockVerifySignatureUC)

	t.Run("should execute use case successfully", func(t *testing.T) {
		expectedHashDomainSeparator := "0x19412cf2425f8663d8b4be7f725bde1ffad8fdff6ed2d9cf0084d785a715912a"
		expectedHashMessage := "0x93d8044d5438b567977f89a53282195effc64659ef5a29e60960e5727dd8580c"
		expectedEncodedData := fmt.Sprintf("\x19\x01%s%s", expectedHashDomainSeparator, expectedHashMessage)

		mockVerifySignatureUC.EXPECT().Execute(ctx, address, signature, expectedEncodedData).Return(nil)

		err := usecase.Execute(ctx, address, signature, payload)
		assert.NoError(t, err)
	})

	t.Run("should fail with same error if verify signature fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockVerifySignatureUC.EXPECT().Execute(ctx, address, signature, gomock.Any()).Return(expectedErr)

		err := usecase.Execute(ctx, address, signature, payload)

		assert.Equal(t, expectedErr, err)
	})
}
