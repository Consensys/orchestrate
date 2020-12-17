// +build unit

package transactions

import (
	"context"
	"fmt"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/validators/mocks"
)

func TestSendDeployTx_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := mocks.NewMockTransactionValidator(ctrl)
	mockSendTxUC := mocks2.NewMockSendTxUseCase(ctrl)

	ctx := context.Background()
	tenantID := "tenantID"
	txRequest := testutils2.FakeTxRequest()
	txRequest.Params.ContractTag = "contractTag"
	txRequest.Params.ContractName = "contractName"

	usecase := NewSendDeployTxUseCase(mockValidator, mockSendTxUC)

	t.Run("should execute use case successfully", func(t *testing.T) {
		txData := "txData"
		txRequestResponse := testutils2.FakeTxRequest()

		mockValidator.EXPECT().ValidateContract(ctx, txRequest.Params).Return(txData, nil)
		mockSendTxUC.EXPECT().Execute(ctx, txRequest, txData, tenantID).Return(txRequestResponse, nil)

		response, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.NoError(t, err)
		assert.Equal(t, txRequestResponse, response)
	})

	t.Run("should fail with same error if validator fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockValidator.EXPECT().ValidateContract(ctx, txRequest.Params).Return("", expectedErr)

		response, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendDeployTxComponent), err)
	})

	t.Run("should fail with same error if send tx use case fails", func(t *testing.T) {
		txData := "txData"
		expectedErr := fmt.Errorf("error")

		mockValidator.EXPECT().ValidateContract(ctx, txRequest.Params).Return(txData, nil)
		mockSendTxUC.EXPECT().Execute(ctx, txRequest, txData, tenantID).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, expectedErr, err)
	})
}
