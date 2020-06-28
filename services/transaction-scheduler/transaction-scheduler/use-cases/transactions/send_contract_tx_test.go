// +build unit

package transactions

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators/mocks"
)

func TestSendContractTx_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := mocks.NewMockTransactionValidator(ctrl)
	mockSendTxUC := mocks2.NewMockSendTxUseCase(ctrl)

	ctx := context.Background()
	chainUUID := "chainUUID"
	tenantID := "tenantID"
	txRequest := testutils.FakeTxRequestEntity()

	usecase := NewSendContractTxUseCase(mockValidator, mockSendTxUC)

	t.Run("should execute use case successfully", func(t *testing.T) {
		txData := "txData"
		txRequestResponse := testutils.FakeTxRequestEntity()

		mockValidator.EXPECT().ValidateMethodSignature(txRequest.Params.MethodSignature, txRequest.Params.Args).Return(txData, nil)
		mockSendTxUC.EXPECT().Execute(ctx, txRequest, txData, chainUUID, tenantID).Return(txRequestResponse, nil)

		response, err := usecase.Execute(ctx, txRequest, chainUUID, tenantID)

		assert.NoError(t, err)
		assert.Equal(t, txRequestResponse, response)
	})

	t.Run("should fail with same error if validator fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockValidator.EXPECT().ValidateMethodSignature(txRequest.Params.MethodSignature, txRequest.Params.Args).Return("", expectedErr)

		response, err := usecase.Execute(ctx, txRequest, chainUUID, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendContractTxComponent), err)
	})

	t.Run("should fail with same error if send tx use case fails", func(t *testing.T) {
		txData := "txData"
		expectedErr := fmt.Errorf("error")

		mockValidator.EXPECT().ValidateMethodSignature(txRequest.Params.MethodSignature, txRequest.Params.Args).Return(txData, nil)
		mockSendTxUC.EXPECT().Execute(ctx, txRequest, txData, chainUUID, tenantID).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, txRequest, chainUUID, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, expectedErr, err)
	})
}
