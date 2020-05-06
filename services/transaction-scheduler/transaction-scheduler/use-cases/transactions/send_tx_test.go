// +build unit

package transactions

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	mocks3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
	mocks4 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators/mocks"
	"testing"
	"time"
)

func TestSendTx_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := mocks4.NewMockTransactionValidator(ctrl)
	mockTxRequestDA := mocks3.NewMockTransactionRequestAgent(ctrl)
	mockStartJob := mocks.NewMockStartJobUseCase(ctrl)

	usecase := NewSendTxUseCase(mockValidator, mockTxRequestDA, mockStartJob)
	ctx := context.Background()

	t.Run("should execute use case successfully", func(t *testing.T) {
		timeNow := time.Now()
		tenantID := "tenantID"
		txRequest := testutils.FakeTransactionRequest()
		expectedJsonParams, _ := utils.ObjectToJSON(txRequest.Params)
		expectedRequestModel := &models.TransactionRequest{
			IdempotencyKey: txRequest.IdempotencyKey,
			Schedule: &models.Schedule{
				TenantID: tenantID,
				ChainID:  txRequest.ChainID,
				Jobs: []*models.Job{{
					Type: types.JobConstantinopleTransaction,
					Transaction: &models.Transaction{
						Sender:    txRequest.Params.To,
						Recipient: txRequest.Params.From,
						Value:     txRequest.Params.Value,
						GasPrice:  txRequest.Params.GasPrice,
						GasLimit:  txRequest.Params.Gas,
						Data:      "", // TODO: Add expected txData here
					},
					Logs: []*models.Log{{
						Status:  types.LogStatusCreated,
						Message: "Job created for contract transaction request",
					}},
					Labels: txRequest.Labels,
				}},
			},
			RequestHash: "requestHash",
			Params:      expectedJsonParams,
		}

		mockValidator.EXPECT().ValidateRequestHash(ctx, txRequest.Params, txRequest.IdempotencyKey).Return("requestHash", nil)
		mockTxRequestDA.EXPECT().SelectOrInsert(ctx, expectedRequestModel).DoAndReturn(func(ctx context.Context, txRequestModel *models.TransactionRequest) error {
			txRequestModel.Schedule.UUID = "scheduleUUID"
			txRequestModel.Schedule.CreatedAt = timeNow
			txRequestModel.Schedule.Jobs[0].UUID = "jobUUID"
			txRequestModel.CreatedAt = timeNow

			return nil
		})
		fmt.Println(expectedRequestModel.Schedule.Jobs[0])
		mockStartJob.EXPECT().Execute(ctx, "jobUUID").Return(nil)

		txResponse, err := usecase.Execute(ctx, txRequest, tenantID)

		expectedResponse := &types.TransactionResponse{
			IdempotencyKey: txRequest.IdempotencyKey,
			Schedule: types.ScheduleResponse{
				UUID:      "scheduleUUID",
				ChainID:   txRequest.ChainID,
				CreatedAt: timeNow,
			},
			CreatedAt: timeNow,
		}
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse.IdempotencyKey, txResponse.IdempotencyKey)
		assert.Equal(t, expectedResponse.Schedule, txResponse.Schedule)
		assert.Equal(t, expectedResponse.CreatedAt, txResponse.CreatedAt)
	})

	t.Run("should fail with same error if validator fails", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.InvalidParameterError("error")

		mockValidator.EXPECT().ValidateRequestHash(ctx, gomock.Any(), gomock.Any()).Return("", expectedErr)

		txResponse, err := usecase.Execute(ctx, txRequest, "tenantID")

		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	t.Run("should fail with same error if SelectOrInsert fails", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.PostgresConnectionError("error")

		mockValidator.EXPECT().ValidateRequestHash(ctx, gomock.Any(), gomock.Any()).Return("requestHash", nil)
		mockTxRequestDA.EXPECT().SelectOrInsert(ctx, gomock.Any()).Return(expectedErr)

		txResponse, err := usecase.Execute(ctx, txRequest, "tenantID")

		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	t.Run("should fail with same error if start job fails", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.KafkaConnectionError("error")

		mockValidator.EXPECT().ValidateRequestHash(ctx, gomock.Any(), gomock.Any()).Return("requestHash", nil)
		mockTxRequestDA.EXPECT().SelectOrInsert(ctx, gomock.Any()).Return(nil)
		mockStartJob.EXPECT().Execute(ctx, gomock.Any()).Return(expectedErr)

		txResponse, err := usecase.Execute(ctx, txRequest, "tenantID")

		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})
}
