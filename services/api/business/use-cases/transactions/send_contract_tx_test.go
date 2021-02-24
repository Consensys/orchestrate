// +build unit

package transactions

import (
	"context"
	"fmt"
	"testing"

	testutils2 "github.com/ConsenSys/orchestrate/pkg/types/testutils"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks2 "github.com/ConsenSys/orchestrate/services/api/business/use-cases/mocks"
)

func TestSendContractTx_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSendTxUC := mocks2.NewMockSendTxUseCase(ctrl)

	ctx := context.Background()
	tenantID := "tenantID"
	txRequest := testutils2.FakeTxRequest()

	usecase := NewSendContractTxUseCase(mockSendTxUC)

	t.Run("should execute use case successfully", func(t *testing.T) {
		txRequestResponse := testutils2.FakeTxRequest()

		mockSendTxUC.EXPECT().Execute(gomock.Any(), txRequest, gomock.Any(), tenantID).Return(txRequestResponse, nil)

		response, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.NoError(t, err)
		assert.Equal(t, txRequestResponse, response)
	})

	t.Run("should fail with same error if send tx use case fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockSendTxUC.EXPECT().Execute(gomock.Any(), txRequest, gomock.Any(), tenantID).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, expectedErr, err)
	})
}
