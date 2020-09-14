// +build unit

package transactions

import (
	"context"
	"fmt"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
)

func TestGetTx_Execute(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockTransactionRequestDA := mocks.NewMockTransactionRequestAgent(ctrl)
	mockGetScheduleUC := mocks2.NewMockGetScheduleUseCase(ctrl)
	tenants := []string{"tenantID"}

	mockDB.EXPECT().TransactionRequest().Return(mockTransactionRequestDA).AnyTimes()

	usecase := NewGetTxUseCase(mockDB, mockGetScheduleUC)

	t.Run("should execute use case successfully", func(t *testing.T) {
		txRequest := testutils.FakeTxRequest(0)
		schedule := testutils2.FakeSchedule()

		mockTransactionRequestDA.EXPECT().FindOneByUUID(ctx, txRequest.UUID, tenants).Return(txRequest, nil)
		mockGetScheduleUC.EXPECT().Execute(ctx, txRequest.Schedule.UUID, tenants).Return(schedule, nil)

		result, err := usecase.Execute(ctx, txRequest.UUID, tenants)

		assert.NoError(t, err)
		assert.Equal(t, txRequest.UUID, result.UUID)
		assert.Equal(t, txRequest.IdempotencyKey, result.IdempotencyKey)
		assert.Equal(t, txRequest.ChainName, result.ChainName)
		assert.Equal(t, txRequest.CreatedAt, result.CreatedAt)
		assert.Equal(t, txRequest.Params, result.Params)
		assert.Equal(t, schedule, result.Schedule)
	})

	t.Run("should fail with same error if FindOneByUUID fails", func(t *testing.T) {
		uuid := "uuid"
		expectedErr := errors.NotFoundError("error")

		mockTransactionRequestDA.EXPECT().FindOneByUUID(ctx, uuid, tenants).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, uuid, tenants)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(getTxComponent), err)
	})

	t.Run("should fail with same error if GetScheduleUseCase fails", func(t *testing.T) {
		txRequest := testutils.FakeTxRequest(0)
		expectedErr := fmt.Errorf("error")

		mockTransactionRequestDA.EXPECT().FindOneByUUID(ctx, txRequest.UUID, tenants).Return(txRequest, nil)
		mockGetScheduleUC.EXPECT().Execute(ctx, txRequest.Schedule.UUID, tenants).Return(nil, expectedErr)

		response, err := usecase.Execute(ctx, txRequest.UUID, tenants)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(getTxComponent), err)
	})
}
