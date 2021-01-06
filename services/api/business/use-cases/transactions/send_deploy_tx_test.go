// +build unit

package transactions

import (
	"context"
	"fmt"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/mocks"
)

func TestSendDeployTx_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSendTxUC := mocks2.NewMockSendTxUseCase(ctrl)
	mockGetContractUC := mocks2.NewMockGetContractUseCase(ctrl)

	ctx := context.Background()
	tenantID := "tenantID"
	txRequest := testutils2.FakeTxRequest()
	txRequest.Params.ContractTag = "contractTag"
	txRequest.Params.ContractName = "contractName"
	txRequest.Params.Args = nil

	usecase := NewSendDeployTxUseCase(mockSendTxUC, mockGetContractUC)

	t.Run("should execute use case successfully", func(t *testing.T) {
		txRequestResponse := testutils2.FakeTxRequest()
		fakeContract := testutils2.FakeContract()

		mockGetContractUC.EXPECT().Execute(ctx, &entities.ContractID{
			Name: txRequest.Params.ContractName,
			Tag:  txRequest.Params.ContractTag,
		}).Return(fakeContract, nil)

		mockSendTxUC.EXPECT().Execute(ctx, txRequest, gomock.Any(), tenantID).Return(txRequestResponse, nil)

		response, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.NoError(t, err)
		assert.Equal(t, txRequestResponse, response)
	})

	t.Run("should fail with same error if validator fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockGetContractUC.EXPECT().Execute(ctx, &entities.ContractID{
			Name: txRequest.Params.ContractName,
			Tag:  txRequest.Params.ContractTag,
		}).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.Nil(t, response)
		assert.Error(t, err)
	})

	t.Run("should fail with same error if send tx use case fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")
		fakeContract := testutils2.FakeContract()

		mockGetContractUC.EXPECT().Execute(ctx, &entities.ContractID{
			Name: txRequest.Params.ContractName,
			Tag:  txRequest.Params.ContractTag,
		}).Return(fakeContract, nil)

		mockSendTxUC.EXPECT().Execute(ctx, txRequest, gomock.Any(), tenantID).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, txRequest, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, expectedErr, err)
	})
}
