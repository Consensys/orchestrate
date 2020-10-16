// +build unit

package controllers

import (
	"bytes"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/key-manager/use-cases/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/key-manager/use-cases/ethereum/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/service/formatters"
)

type ethereumCtrlTestSuite struct {
	suite.Suite
	createAccountUC *mocks.MockCreateAccountUseCase
	router          *mux.Router
}

func (s *ethereumCtrlTestSuite) CreateAccount() ethereum.CreateAccountUseCase {
	return s.createAccountUC
}

var _ ethereum.UseCases = &ethereumCtrlTestSuite{}

func TestEthereumController(t *testing.T) {
	s := new(ethereumCtrlTestSuite)
	suite.Run(t, s)
}

func (s *ethereumCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.createAccountUC = mocks.NewMockCreateAccountUseCase(ctrl)
	s.router = mux.NewRouter()

	controller := NewEthereumController(s)
	controller.Append(s.router)
}

func (s *ethereumCtrlTestSuite) TestEthereumController_Create() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		createAccountRequest := testutils.FakeCreateETHAccountRequest()
		requestBytes, _ := json.Marshal(createAccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/ethereum/accounts", bytes.NewReader(requestBytes))

		fakeETHAccount := testutils.FakeETHAccount()

		s.createAccountUC.EXPECT().
			Execute(gomock.Any(), createAccountRequest.Namespace, "").
			Return(fakeETHAccount, nil)

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

		s.createAccountUC.EXPECT().
			Execute(gomock.Any(), gomock.Any(), "").
			Return(nil, errors.HashicorpVaultConnectionError("error"))

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

		s.createAccountUC.EXPECT().
			Execute(gomock.Any(), importAccountRequest.Namespace, importAccountRequest.PrivateKey).
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

		s.createAccountUC.EXPECT().
			Execute(gomock.Any(), importAccountRequest.Namespace, importAccountRequest.PrivateKey).
			Return(nil, errors.InvalidParameterError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})
}
