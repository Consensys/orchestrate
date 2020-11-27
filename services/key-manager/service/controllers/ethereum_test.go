// +build unit

package controllers

import (
	"bytes"
	"fmt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/service/formatters"
)

type ethereumCtrlTestSuite struct {
	suite.Suite
	vault  *mocks.MockVault
	router *mux.Router
}

const (
	inputTestAddress     = "0x7e654d251da770a068413677967f6d3ea2feA9e4"
	mixedCaseTestAddress = "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4"
)

func TestEthereumController(t *testing.T) {
	s := new(ethereumCtrlTestSuite)
	suite.Run(t, s)
}

func (s *ethereumCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.vault = mocks.NewMockVault(ctrl)
	s.router = mux.NewRouter()

	controller := NewEthereumController(s.vault)
	controller.Append(s.router)
}

func (s *ethereumCtrlTestSuite) TestEthereumController_Create() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		createAccountRequest := testutils.FakeCreateETHAccountRequest()
		requestBytes, _ := json.Marshal(createAccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/ethereum/accounts", bytes.NewReader(requestBytes))

		fakeETHAccount := testutils.FakeETHAccount()

		s.vault.EXPECT().ETHCreateAccount(createAccountRequest.Namespace).Return(fakeETHAccount, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatETHAccountResponse(fakeETHAccount)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		createAccountRequest := testutils.FakeCreateETHAccountRequest()
		requestBytes, _ := json.Marshal(createAccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/ethereum/accounts", bytes.NewReader(requestBytes))

		s.vault.EXPECT().ETHCreateAccount(gomock.Any()).Return(nil, errors.HashicorpVaultConnectionError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func (s *ethereumCtrlTestSuite) TestEthereumController_Import() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		importAccountRequest := testutils.FakeImportETHAccountRequest()
		requestBytes, _ := json.Marshal(importAccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/ethereum/accounts/import", bytes.NewReader(requestBytes))

		fakeETHAccount := testutils.FakeETHAccount()

		s.vault.EXPECT().
			ETHImportAccount(importAccountRequest.Namespace, importAccountRequest.PrivateKey).
			Return(fakeETHAccount, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatETHAccountResponse(fakeETHAccount)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		importAccountRequest := testutils.FakeImportETHAccountRequest()
		requestBytes, _ := json.Marshal(importAccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/ethereum/accounts/import", bytes.NewReader(requestBytes))

		s.vault.EXPECT().
			ETHImportAccount(importAccountRequest.Namespace, importAccountRequest.PrivateKey).
			Return(nil, errors.InvalidParameterError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})
}

func (s *ethereumCtrlTestSuite) TestEthereumController_Sign() {
	url := fmt.Sprintf("/ethereum/accounts/%v/sign", inputTestAddress)

	s.T().Run("should execute request successfully", func(t *testing.T) {
		signature := "0xsignature"
		payloadRequest := &keymanager.PayloadRequest{
			Data:      "my data to sign",
			Namespace: "namespace",
		}
		requestBytes, _ := json.Marshal(payloadRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(requestBytes))

		s.vault.EXPECT().
			ETHSign(mixedCaseTestAddress, payloadRequest.Namespace, payloadRequest.Data).
			Return(signature, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, signature, rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		payloadRequest := &keymanager.PayloadRequest{
			Data:      "my data to sign",
			Namespace: "namespace",
		}
		requestBytes, _ := json.Marshal(payloadRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(requestBytes))

		s.vault.EXPECT().
			ETHSign(mixedCaseTestAddress, payloadRequest.Namespace, payloadRequest.Data).
			Return("", errors.InvalidParameterError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})
}

func (s *ethereumCtrlTestSuite) TestEthereumController_SignTransaction() {
	url := fmt.Sprintf("/ethereum/accounts/%v/sign-transaction", inputTestAddress)

	s.T().Run("should execute request successfully", func(t *testing.T) {
		signature := "0xsignature"
		signRequest := testutils.FakeSignETHTransactionRequest()
		requestBytes, _ := json.Marshal(signRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(requestBytes))

		s.vault.EXPECT().ETHSignTransaction(mixedCaseTestAddress, signRequest).Return(signature, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, signature, rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		signRequest := testutils.FakeSignETHTransactionRequest()
		requestBytes, _ := json.Marshal(signRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(requestBytes))

		s.vault.EXPECT().
			ETHSignTransaction(mixedCaseTestAddress, signRequest).
			Return("", errors.ServiceConnectionError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func (s *ethereumCtrlTestSuite) TestEthereumController_SignQuorumPrivateTransaction() {
	url := fmt.Sprintf("/ethereum/accounts/%v/sign-quorum-private-transaction", inputTestAddress)

	s.T().Run("should execute request successfully", func(t *testing.T) {
		signature := "0xsignature"
		signRequest := testutils.FakeSignQuorumPrivateTransactionRequest()
		requestBytes, _ := json.Marshal(signRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(requestBytes))

		s.vault.EXPECT().
			ETHSignQuorumPrivateTransaction(mixedCaseTestAddress, signRequest).
			Return(signature, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, signature, rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		signRequest := testutils.FakeSignQuorumPrivateTransactionRequest()
		requestBytes, _ := json.Marshal(signRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(requestBytes))

		s.vault.EXPECT().
			ETHSignQuorumPrivateTransaction(mixedCaseTestAddress, signRequest).
			Return("", errors.ServiceConnectionError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func (s *ethereumCtrlTestSuite) TestEthereumController_SignEEATransaction() {
	url := fmt.Sprintf("/ethereum/accounts/%v/sign-eea-transaction", inputTestAddress)

	s.T().Run("should execute request successfully", func(t *testing.T) {
		signature := "0xsignature"
		signRequest := testutils.FakeSignEEATransactionRequest()
		requestBytes, _ := json.Marshal(signRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(requestBytes))

		s.vault.EXPECT().ETHSignEEATransaction(mixedCaseTestAddress, signRequest).Return(signature, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, signature, rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		signRequest := testutils.FakeSignEEATransactionRequest()
		requestBytes, _ := json.Marshal(signRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(requestBytes))

		s.vault.EXPECT().
			ETHSignEEATransaction(mixedCaseTestAddress, signRequest).
			Return("", errors.InvalidParameterError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})
}
