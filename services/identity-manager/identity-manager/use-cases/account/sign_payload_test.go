// +build unit

package account

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client/mock"
)

func TestSignPayload_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockKeyManagerClient(ctrl)

	usecase := NewSignPayloadUseCase(mockClient)

	tenantID := "tenantID"

	t.Run("should execute use case successfully", func(t *testing.T) {
		acc := testutils.FakeAccountModel()
		payload := "messageToSign"
		signedPayload := "0xACDEF01234567890"

		mockClient.EXPECT().ETHSign(ctx, acc.Address, &keymanager.PayloadRequest{
			Data:      payload,
			Namespace: tenantID,
		}).Return(signedPayload, nil)

		resp, err := usecase.Execute(ctx, acc.Address, payload, tenantID)

		assert.NoError(t, err)
		assert.Equal(t, signedPayload, resp)
	})

	t.Run("should fail with same error if ETH sign payload fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		acc := testutils.FakeAccountModel()

		mockClient.EXPECT().ETHSign(ctx, acc.Address, gomock.Any()).Return("", expectedErr)

		_, err := usecase.Execute(ctx, acc.Address, "payload", tenantID)

		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(signPayloadComponent), err)
	})
}
