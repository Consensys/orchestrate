// +build unit

package transactions

import (
	"context"
	"fmt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
)

func TestSearchTxs_Execute(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockTransactionRequestDA := mocks.NewMockTransactionRequestAgent(ctrl)
	mockGetTxUC := mocks2.NewMockGetTxUseCase(ctrl)
	tenantID := "tenantID"
	filter := &entities.TransactionFilters{}

	mockDB.EXPECT().TransactionRequest().Return(mockTransactionRequestDA).AnyTimes()

	usecase := NewSearchTransactionsUseCase(mockDB, mockGetTxUC)

	t.Run("should execute use case successfully", func(t *testing.T) {
		txRequestModels := []*models.TransactionRequest{testutils.FakeTxRequest(0), testutils.FakeTxRequest(1)}
		txRequest0 := testutils2.FakeTxRequestEntity()
		txRequest1 := testutils2.FakeTxRequestEntity()

		mockTransactionRequestDA.EXPECT().Search(ctx, tenantID, filter).Return(txRequestModels, nil)
		mockGetTxUC.EXPECT().Execute(ctx, txRequestModels[0].UUID, tenantID).Return(txRequest0, nil)
		mockGetTxUC.EXPECT().Execute(ctx, txRequestModels[1].UUID, tenantID).Return(txRequest1, nil)

		result, err := usecase.Execute(ctx, filter, tenantID)

		assert.Nil(t, err)

		assert.Equal(t, txRequest0.UUID, result[0].UUID)
		assert.Equal(t, txRequest0.IdempotencyKey, result[0].IdempotencyKey)
		assert.Equal(t, txRequest0.CreatedAt, result[0].CreatedAt)
		assert.Equal(t, txRequest0.Params, result[0].Params)

		assert.Equal(t, txRequest1.UUID, result[1].UUID)
		assert.Equal(t, txRequest1.IdempotencyKey, result[1].IdempotencyKey)
		assert.Equal(t, txRequest1.CreatedAt, result[1].CreatedAt)
		assert.Equal(t, txRequest1.Params, result[1].Params)
	})

	t.Run("should fail with same error if Search fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		mockTransactionRequestDA.EXPECT().Search(ctx, tenantID, filter).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, filter, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(searchTxsComponent), err)
	})

	t.Run("should fail with same error if GetTxUseCase fails", func(t *testing.T) {
		txRequestModels := []*models.TransactionRequest{testutils.FakeTxRequest(0)}
		expectedErr := fmt.Errorf("error")

		mockTransactionRequestDA.EXPECT().Search(ctx, tenantID, filter).Return(txRequestModels, nil)
		mockGetTxUC.EXPECT().Execute(ctx, txRequestModels[0].UUID, tenantID).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, filter, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(searchTxsComponent), err)
	})
}
