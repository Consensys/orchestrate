// +build unit

package transactions

import (
	"context"
	"fmt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	testutils3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models/testutils"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/mocks"
)

func TestSearchTxs_Execute(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockTransactionRequestDA := mocks.NewMockTransactionRequestAgent(ctrl)
	mockGetTxUC := mocks2.NewMockGetTxUseCase(ctrl)
	tenantID := "tenantID"
	filter := &entities.TransactionRequestFilters{}

	mockDB.EXPECT().TransactionRequest().Return(mockTransactionRequestDA).AnyTimes()

	usecase := NewSearchTransactionsUseCase(mockDB, mockGetTxUC)

	t.Run("should execute use case successfully", func(t *testing.T) {
		txRequestModels := []*models.TransactionRequest{testutils.FakeTxRequest(0), testutils.FakeTxRequest(1)}
		txRequest0 := testutils3.FakeTxRequest()
		txRequest1 := testutils3.FakeTxRequest()

		mockTransactionRequestDA.EXPECT().Search(gomock.Any(), filter, []string{tenantID}).Return(txRequestModels, nil)
		mockGetTxUC.EXPECT().Execute(gomock.Any(), txRequestModels[0].Schedule.UUID, []string{tenantID}).Return(txRequest0, nil)
		mockGetTxUC.EXPECT().Execute(gomock.Any(), txRequestModels[1].Schedule.UUID, []string{tenantID}).Return(txRequest1, nil)

		result, err := usecase.Execute(ctx, filter, []string{tenantID})

		assert.Nil(t, err)

		assert.Equal(t, txRequest0.Schedule.UUID, result[0].Schedule.UUID)
		assert.Equal(t, txRequest0.IdempotencyKey, result[0].IdempotencyKey)
		assert.Equal(t, txRequest0.CreatedAt, result[0].CreatedAt)
		assert.Equal(t, txRequest0.Params, result[0].Params)

		assert.Equal(t, txRequest1.Schedule.UUID, result[1].Schedule.UUID)
		assert.Equal(t, txRequest1.IdempotencyKey, result[1].IdempotencyKey)
		assert.Equal(t, txRequest1.CreatedAt, result[1].CreatedAt)
		assert.Equal(t, txRequest1.Params, result[1].Params)
	})

	t.Run("should fail with same error if Search fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		mockTransactionRequestDA.EXPECT().Search(gomock.Any(), filter, []string{tenantID}).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, filter, []string{tenantID})

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(searchTxsComponent), err)
	})

	t.Run("should fail with same error if GetTxUseCase fails", func(t *testing.T) {
		txRequestModels := []*models.TransactionRequest{testutils.FakeTxRequest(0)}
		expectedErr := fmt.Errorf("error")

		mockTransactionRequestDA.EXPECT().Search(gomock.Any(), filter, []string{tenantID}).Return(txRequestModels, nil)
		mockGetTxUC.EXPECT().Execute(gomock.Any(), txRequestModels[0].Schedule.UUID, []string{tenantID}).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, filter, []string{tenantID})

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(searchTxsComponent), err)
	})
}
