// +build unit

package controllers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/formatters"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/testutils"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions/mocks"
)

type transactionsControllerTestSuite struct {
	suite.Suite
	controller            *TransactionsController
	router                *mux.Router
	sendContractTxUseCase *mocks.MockSendContractTxUseCase
	sendDeployTxUseCase   *mocks.MockSendDeployTxUseCase
	sendTxUseCase         *mocks.MockSendTxUseCase
	ctx                   context.Context
	tenantID              string
	chainUUID             string
}

func (s *transactionsControllerTestSuite) SendContractTransaction() transactions.SendContractTxUseCase {
	return s.sendContractTxUseCase
}

func (s *transactionsControllerTestSuite) SendDeployTransaction() transactions.SendDeployTxUseCase {
	return s.sendDeployTxUseCase
}

func (s *transactionsControllerTestSuite) SendTransaction() transactions.SendTxUseCase {
	return s.sendTxUseCase
}

var _ transactions.UseCases = &transactionsControllerTestSuite{}

func TestTransactionsController(t *testing.T) {
	s := new(transactionsControllerTestSuite)
	suite.Run(t, s)
}

func (s *transactionsControllerTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.sendContractTxUseCase = mocks.NewMockSendContractTxUseCase(ctrl)
	s.sendDeployTxUseCase = mocks.NewMockSendDeployTxUseCase(ctrl)
	s.sendTxUseCase = mocks.NewMockSendTxUseCase(ctrl)
	s.tenantID = "tenantId"
	s.chainUUID = uuid.NewV4().String()
	s.ctx = context.WithValue(context.Background(), multitenancy.TenantIDKey, s.tenantID)

	s.router = mux.NewRouter()
	s.controller = NewTransactionsController(s)
	s.controller.Append(s.router)
}

func (s *transactionsControllerTestSuite) TestTransactionsController_Send() {

	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		urlPath := fmt.Sprintf("/transactions/%v/send", s.chainUUID)

		txRequest := testutils.FakeSendTransactionRequest()
		requestBytes, err := json.Marshal(txRequest)
		if err != nil {
			return
		}
		txRequestEntity := formatters.FormatSendTxRequest(txRequest)

		httpRequest := httptest.
			NewRequest(http.MethodPost, urlPath, bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		testutils2.FakeTxRequestEntity()
		txRequestEntityResp := testutils2.FakeTxRequestEntity()

		s.sendContractTxUseCase.EXPECT().
			Execute(gomock.Any(), txRequestEntity, s.chainUUID, s.tenantID).
			Return(txRequestEntityResp, nil).
			Times(1)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatTxResponse(txRequestEntityResp)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusAccepted, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 422 if use case fails with InvalidParameterError", func(t *testing.T) {
		txRequest := testutils.FakeSendTransactionRequest()
		requestBytes, _ := json.Marshal(txRequest)
		txRequestEntity := formatters.FormatSendTxRequest(txRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/transactions/%s/send", s.chainUUID),
			bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.sendContractTxUseCase.EXPECT().
			Execute(gomock.Any(), txRequestEntity, s.chainUUID, s.tenantID).
			Return(nil, errors.InvalidParameterError("error")).
			Times(1)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})

	s.T().Run("should fail with Bad request if invalid format", func(t *testing.T) {
		txRequest := testutils.FakeSendTransactionRequest()
		txRequest.IdempotencyKey = ""
		requestBytes, _ := json.Marshal(txRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/transactions/%s/send", s.chainUUID),
			bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	s.T().Run("should fail with 400 if request fails with InvalidParameterError for private txs", func(t *testing.T) {
		rw := httptest.NewRecorder()
		txRequest := testutils.FakeSendTesseraRequest()
		txRequest.Params.PrivateFrom = ""
		requestBytes, _ := json.Marshal(txRequest)

		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/transactions/%s/send", s.chainUUID),
			bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})
}

func (s *transactionsControllerTestSuite) TestTransactionsController_Deploy() {
	urlPath := fmt.Sprintf("/transactions/%v/deploy-contract", s.chainUUID)

	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()

		txRequest := testutils.FakeDeployContractRequest()
		requestBytes, _ := json.Marshal(txRequest)
		txRequestEntity := formatters.FormatDeployContractRequest(txRequest)

		httpRequest := httptest.NewRequest(http.MethodPost, urlPath, bytes.NewReader(requestBytes)).WithContext(s.ctx)

		txRequestEntityResp := testutils2.FakeTxRequestEntity()

		s.sendDeployTxUseCase.EXPECT().
			Execute(gomock.Any(), txRequestEntity, s.chainUUID, s.tenantID).
			Return(txRequestEntityResp, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatTxResponse(txRequestEntityResp)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusAccepted, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 422 if use case fails with InvalidParameterError", func(t *testing.T) {
		txRequest := testutils.FakeDeployContractRequest()
		requestBytes, _ := json.Marshal(txRequest)
		txRequestEntity := formatters.FormatDeployContractRequest(txRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, urlPath, bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.sendDeployTxUseCase.EXPECT().
			Execute(gomock.Any(), txRequestEntity, s.chainUUID, s.tenantID).
			Return(nil, errors.InvalidParameterError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})

	s.T().Run("should fail with Bad request if invalid format", func(t *testing.T) {
		txRequest := testutils.FakeDeployContractRequest()
		txRequest.IdempotencyKey = ""
		requestBytes, _ := json.Marshal(txRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, urlPath, bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	s.T().Run("should fail with 400 if request fails with InvalidParameterError for private txs", func(t *testing.T) {
		rw := httptest.NewRecorder()
		txRequest := testutils.FakeDeployContractRequest()
		txRequest.Params.PrivateFrom = "PrivateFrom"
		requestBytes, _ := json.Marshal(txRequest)

		httpRequest := httptest.NewRequest(http.MethodPost, urlPath, bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})
}

func (s *transactionsControllerTestSuite) TestTransactionsController_SendRaw() {

	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		urlPath := fmt.Sprintf("/transactions/%v/send-raw", s.chainUUID)

		txRequest := testutils.FakeSendRawTransactionRequest()
		requestBytes, err := json.Marshal(txRequest)
		if err != nil {
			return
		}
		// txRequestEntity := formatters.FormatSendRawRequest(txRequest)

		httpRequest := httptest.
			NewRequest(http.MethodPost, urlPath, bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		testutils2.FakeTxRequestEntity()
		txRequestEntityResp := testutils2.FakeTxRequestEntity()

		s.sendTxUseCase.EXPECT().
			Execute(gomock.Any(), gomock.Any(), "", s.chainUUID, s.tenantID).
			Return(txRequestEntityResp, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatTxResponse(txRequestEntityResp)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusAccepted, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 422 if use case fails with InvalidParameterError", func(t *testing.T) {
		txRequest := testutils.FakeSendRawTransactionRequest()
		requestBytes, _ := json.Marshal(txRequest)
		txRequestEntity := formatters.FormatSendRawRequest(txRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/transactions/%s/send-raw", s.chainUUID),
			bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.sendTxUseCase.EXPECT().
			Execute(gomock.Any(), txRequestEntity, "", s.chainUUID, s.tenantID).
			Return(nil, errors.InvalidParameterError("error")).
			Times(1)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})

	s.T().Run("should fail with Bad request if invalid format", func(t *testing.T) {
		txRequest := testutils.FakeSendRawTransactionRequest()
		txRequest.IdempotencyKey = ""
		requestBytes, _ := json.Marshal(txRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/transactions/%s/send-raw", s.chainUUID),
			bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})
}
