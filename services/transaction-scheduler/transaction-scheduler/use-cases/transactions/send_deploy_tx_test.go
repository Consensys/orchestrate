// +build unit

package transactions

import (
	"context"
	"fmt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"
)

func TestSendDeployTx_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := mocks.NewMockTransactionValidator(ctrl)
	mockSendTxUC := mocks2.NewMockSendTxUseCase(ctrl)

	ctx := context.Background()
	chainUUID := "chainUUID"
	tenantID := "tenantID"
	txRequest := testutils.FakeTxRequestEntity()
	txRequest.Params.ContractTag = "contractTag"
	txRequest.Params.ContractName = "contractName"

	usecase := NewSendDeployTxUseCase(mockValidator, mockSendTxUC)

	t.Run("should execute use case successfully", func(t *testing.T) {
		txData := "txData"
		txRequestResponse := testutils.FakeTxRequestEntity()

		mockValidator.EXPECT().ValidateContract(ctx, txRequest.Params).Return(txData, nil)
		mockSendTxUC.EXPECT().Execute(ctx, txRequest, txData, chainUUID, tenantID).Return(txRequestResponse, nil)

		response, err := usecase.Execute(ctx, txRequest, chainUUID, tenantID)

		assert.Nil(t, err)
		assert.Equal(t, txRequestResponse, response)
	})

	t.Run("should fail with same error if validator fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockValidator.EXPECT().ValidateContract(ctx, txRequest.Params).Return("", expectedErr)

		response, err := usecase.Execute(ctx, txRequest, chainUUID, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendDeployTxComponent), err)
	})

	t.Run("should fail with same error if send tx use case fails", func(t *testing.T) {
		txData := "txData"
		expectedErr := fmt.Errorf("error")

		mockValidator.EXPECT().ValidateContract(ctx, txRequest.Params).Return(txData, nil)
		mockSendTxUC.EXPECT().Execute(ctx, txRequest, txData, chainUUID, tenantID).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, txRequest, chainUUID, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, expectedErr, err)
	})
}
