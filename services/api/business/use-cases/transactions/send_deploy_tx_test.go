// +build unit

package transactions

import (
	"context"
	"fmt"
	"testing"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	testutils2 "github.com/consensys/orchestrate/pkg/types/testutils"

	mocks2 "github.com/consensys/orchestrate/services/api/business/use-cases/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSendDeployTx_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSendTxUC := mocks2.NewMockSendTxUseCase(ctrl)
	mockGetContractUC := mocks2.NewMockGetContractUseCase(ctrl)

	ctx := context.Background()
	txRequest := testutils2.FakeTxRequest()
	txRequest.Params.Args = nil

	userInfo := multitenancy.NewUserInfo("tenantOne", "username")
	usecase := NewSendDeployTxUseCase(mockSendTxUC, mockGetContractUC)

	t.Run("should execute use case successfully", func(t *testing.T) {
		txRequestResponse := testutils2.FakeTxRequest()
		fakeContract := testutils2.FakeContract()

		mockGetContractUC.EXPECT().Execute(gomock.Any(), txRequest.Params.ContractName, txRequest.Params.ContractTag).Return(fakeContract, nil)
		mockSendTxUC.EXPECT().Execute(gomock.Any(), txRequest, gomock.Any(), userInfo).Return(txRequestResponse, nil)

		response, err := usecase.Execute(ctx, txRequest, userInfo)

		assert.NoError(t, err)
		assert.Equal(t, txRequestResponse, response)
	})

	t.Run("should fail with same error if validator fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockGetContractUC.EXPECT().Execute(gomock.Any(), txRequest.Params.ContractName, txRequest.Params.ContractTag).Return(nil, expectedErr)
		response, err := usecase.Execute(ctx, txRequest, userInfo)

		assert.Nil(t, response)
		assert.Error(t, err)
	})

	t.Run("should fail with same error if send tx use case fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")
		fakeContract := testutils2.FakeContract()

		mockGetContractUC.EXPECT().Execute(gomock.Any(), txRequest.Params.ContractName, txRequest.Params.ContractTag).Return(fakeContract, nil)
		mockSendTxUC.EXPECT().Execute(gomock.Any(), txRequest, gomock.Any(), userInfo).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, txRequest, userInfo)

		assert.Nil(t, response)
		assert.Equal(t, expectedErr, err)
	})
}
