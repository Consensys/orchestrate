// +build unit

package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
)

type transactionsControllerTestSuite struct {
	suite.Suite
	controller             *TransactionsController
	sendTransactionUseCase *mocks.MockSendTxUseCase
	ctx                    context.Context
	tenantID               string
	chainUUID              string
}

func (s *transactionsControllerTestSuite) SendTransaction() transactions.SendTxUseCase {
	return s.sendTransactionUseCase
}

var _ transactions.UseCases = &transactionsControllerTestSuite{}


func TestTransactionsController(t *testing.T) {
	s := new(transactionsControllerTestSuite)
	suite.Run(t, s)
}

func (s *transactionsControllerTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.sendTransactionUseCase = mocks.NewMockSendTxUseCase(ctrl)
	s.tenantID = "tenantId"
	s.chainUUID = uuid.NewV4().String()
	s.ctx = context.WithValue(context.Background(), multitenancy.TenantIDKey, s.tenantID)

	s.controller = NewTransactionsController(s)
}

func (s *transactionsControllerTestSuite) TestTransactionsController_Send() {
	txRequest := testutils.FakeTransactionRequest()
	requestBytes, _ := json.Marshal(txRequest)

	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/transactions/%s/send", s.chainUUID), 
			bytes.NewReader(requestBytes)).WithContext(s.ctx)
		txResponse := testutils.FakeTransactionResponse()

		s.sendTransactionUseCase.EXPECT().Execute(s.ctx, txRequest, gomock.Any(), s.tenantID).Return(txResponse, nil).Times(1)

		s.controller.Send(rw, httpRequest)

		expectedBody, _ := utils.ObjectToJSON(txResponse)
		assert.Equal(t, expectedBody+"\n", rw.Body.String())
		assert.Equal(t, http.StatusAccepted, rw.Code)
	})

	s.T().Run("should fail with Bad request if invalid format", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		txRequest.IdempotencyKey = ""
		requestBytes, _ := json.Marshal(txRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/transactions/%s/send", s.chainUUID), 
			bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.controller.Send(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 422 if use case fails with InvalidParameterError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/transactions/%s/send", s.chainUUID), 
			bytes.NewReader(requestBytes)).WithContext(s.ctx)
		s.sendTransactionUseCase.EXPECT().Execute(s.ctx, txRequest, gomock.Any(), s.tenantID).Return(nil, errors.InvalidParameterError("error")).Times(1)

		s.controller.Send(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})
}
