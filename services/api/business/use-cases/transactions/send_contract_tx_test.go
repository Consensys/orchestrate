package transactions

import (
	"context"
	"fmt"
	"testing"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	testutils2 "github.com/consensys/orchestrate/pkg/types/testutils"

	"github.com/consensys/orchestrate/services/api/business/use-cases/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSendContractTx_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSendTxUC := mocks.NewMockSendTxUseCase(ctrl)
	mockGetContractUC := mocks.NewMockGetContractUseCase(ctrl)

	ctx := context.Background()
	txRequest := testutils2.FakeTxRequest()
	c := testutils2.FakeContract()
	txRequestResponse := testutils2.FakeTxRequest()

	userInfo := multitenancy.NewUserInfo("tenantOne", "username")
	usecase := NewSendContractTxUseCase(mockSendTxUC, mockGetContractUC)

	t.Run("should execute use case successfully", func(t *testing.T) {
		mockGetContractUC.EXPECT().Execute(gomock.Any(), txRequest.Params.ContractName, txRequest.Params.ContractTag).Return(c, nil)
		mockSendTxUC.EXPECT().Execute(gomock.Any(), txRequest, gomock.Any(), userInfo).Return(txRequestResponse, nil)

		response, err := usecase.Execute(ctx, txRequest, userInfo)

		assert.NoError(t, err)
		assert.Equal(t, txRequestResponse, response)
	})

	t.Run("should parse arguments of a complex method successfully", func(t *testing.T) {
		newContract := testutils2.FakeContract()
		newContract.RawABI = testutils2.ContractABIStruct
		newTxRequest := testutils2.FakeTxRequest()
		newTxRequest.Params.MethodSignature = "multipleTransfer((address,(uint256,address)[]))"
		newTxRequest.Params.Args = []interface{}{
			map[string]interface{}{
				"token": "0xdbb881a51CD4023E4400CEF3ef73046743f08da3",
				"recipients": []map[string]interface{}{
					{
						"amount":    500,
						"recipient": "0xdbb881a51CD4023E4400CEF3ef73046743f08da3",
					},
				},
			},
		}
		expectedTxData := hexutil.MustDecode("0x52ca78230000000000000000000000000000000000000000000000000000000000000020000000000000000000000000dbb881a51cd4023e4400cef3ef73046743f08da30000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000001f4000000000000000000000000dbb881a51cd4023e4400cef3ef73046743f08da3")

		mockGetContractUC.EXPECT().Execute(gomock.Any(), newTxRequest.Params.ContractName, newTxRequest.Params.ContractTag).Return(newContract, nil)
		mockSendTxUC.EXPECT().Execute(gomock.Any(), newTxRequest, expectedTxData, userInfo).Return(txRequestResponse, nil)

		response, err := usecase.Execute(ctx, newTxRequest, userInfo)

		assert.NoError(t, err)
		assert.Equal(t, txRequestResponse, response)
	})

	t.Run("should parse arguments of a simple method successfully", func(t *testing.T) {
		newContract := testutils2.FakeContract()
		newContract.RawABI = testutils2.ContractABIStruct
		newTxRequest := testutils2.FakeTxRequest()
		newTxRequest.Params.MethodSignature = "singleTransfer(address,address,uint256)"
		newTxRequest.Params.Args = []interface{}{"0xdbb881a51CD4023E4400CEF3ef73046743f08da3", "0xdbb881a51CD4023E4400CEF3ef73046743f08da3", 500}
		expectedTxData := hexutil.MustDecode("0xed629438000000000000000000000000dbb881a51cd4023e4400cef3ef73046743f08da3000000000000000000000000dbb881a51cd4023e4400cef3ef73046743f08da300000000000000000000000000000000000000000000000000000000000001f4")

		mockGetContractUC.EXPECT().Execute(gomock.Any(), newTxRequest.Params.ContractName, newTxRequest.Params.ContractTag).Return(newContract, nil)
		mockSendTxUC.EXPECT().Execute(gomock.Any(), newTxRequest, expectedTxData, userInfo).Return(txRequestResponse, nil)

		response, err := usecase.Execute(ctx, newTxRequest, userInfo)

		assert.NoError(t, err)
		assert.Equal(t, txRequestResponse, response)
	})

	t.Run("should fail with same error if get contract use case fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockGetContractUC.EXPECT().Execute(gomock.Any(), txRequest.Params.ContractName, txRequest.Params.ContractTag).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, txRequest, userInfo)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).SetComponent(sendContractTxComponent), err)
	})

	t.Run("should fail with same error if send tx use case fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockGetContractUC.EXPECT().Execute(gomock.Any(), txRequest.Params.ContractName, txRequest.Params.ContractTag).Return(c, nil)
		mockSendTxUC.EXPECT().Execute(gomock.Any(), txRequest, gomock.Any(), userInfo).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, txRequest, userInfo)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).SetComponent(sendContractTxComponent), err)
	})
}
