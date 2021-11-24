// +build unit

package transactions

import (
	"context"
	"fmt"
	"testing"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	testutils2 "github.com/consensys/orchestrate/pkg/types/testutils"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks2 "github.com/consensys/orchestrate/services/api/business/use-cases/mocks"
)

func TestSendContractTx_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSendTxUC := mocks2.NewMockSendTxUseCase(ctrl)

	ctx := context.Background()
	txRequest := testutils2.FakeTxRequest()

	userInfo := multitenancy.NewUserInfo("tenantOne", "username")
	usecase := NewSendContractTxUseCase(mockSendTxUC)

	t.Run("should execute use case successfully", func(t *testing.T) {
		txRequestResponse := testutils2.FakeTxRequest()

		mockSendTxUC.EXPECT().Execute(gomock.Any(), txRequest, gomock.Any(), userInfo).Return(txRequestResponse, nil)

		response, err := usecase.Execute(ctx, txRequest, userInfo)

		assert.NoError(t, err)
		assert.Equal(t, txRequestResponse, response)
	})

	t.Run("should fail with same error if send tx use case fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockSendTxUC.EXPECT().Execute(gomock.Any(), txRequest, gomock.Any(), userInfo).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, txRequest, userInfo)

		assert.Nil(t, response)
		assert.Equal(t, expectedErr, err)
	})
}
