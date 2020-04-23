// +build unit

package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
	"net/http"

	"net/http/httptest"

	"testing"
)

const tenantId = "tenantId"

type transactionsControllerTestSuite struct {
	suite.Suite
	controller             *TransactionsController
	sendTransactionUseCase *mocks.MockSendTxUseCase
	ctx                    context.Context
}

func TestTransactionsController(t *testing.T) {
	s := new(transactionsControllerTestSuite)
	suite.Run(t, s)
}

func (s *transactionsControllerTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.sendTransactionUseCase = mocks.NewMockSendTxUseCase(ctrl)
	s.ctx = context.WithValue(context.Background(), multitenancy.TenantIDKey, tenantId)

	s.controller = NewTransactionsController(s.sendTransactionUseCase)
}

func (s *transactionsControllerTestSuite) TestTransactionsController_Send() {
	txRequest := testutils.FakeTransactionRequest()
	requestBytes, _ := json.Marshal(txRequest)
	txResponse := &types.TransactionResponse{
		IdempotencyKey: txRequest.IdempotencyKey,
		ChainID:        txRequest.ChainID,
		Method:         types.MethodSendRawTransaction,
		Schedule:       types.ScheduleResponse{},
	}

	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/transactions/send", bytes.NewReader(requestBytes)).WithContext(s.ctx)
		s.sendTransactionUseCase.EXPECT().Execute(s.ctx, txRequest, tenantId).Return(txResponse, nil).Times(1)

		s.controller.Send(rw, httpRequest)

		expectedBody, _ := utils.ObjectToJSON(txResponse)
		assert.Equal(t, expectedBody+"\n", rw.Body.String())
		assert.Equal(t, http.StatusAccepted, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with not found if use case fails with NotFoundError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/transactions/send", bytes.NewReader(requestBytes)).WithContext(s.ctx)
		s.sendTransactionUseCase.EXPECT().Execute(s.ctx, txRequest, tenantId).Return(nil, errors.NotFoundError("error")).Times(1)

		s.controller.Send(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}
