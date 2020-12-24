// +build unit

package controllers

import (
	"bytes"
	"context"
	"fmt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/service/formatters"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/mocks"
)

const endpoint = "/faucets"

type faucetsCtrlTestSuite struct {
	suite.Suite
	registerFaucetUC *mocks.MockRegisterFaucetUseCase
	getFaucetUC      *mocks.MockGetFaucetUseCase
	searchFaucetUC   *mocks.MockSearchFaucetsUseCase
	updateFaucetUC   *mocks.MockUpdateFaucetUseCase
	deleteFaucetUC   *mocks.MockDeleteFaucetUseCase
	keyManagerClient *mock.MockKeyManagerClient
	ctx              context.Context
	tenants          []string
	router           *mux.Router
}

var _ usecases.FaucetUseCases = &faucetsCtrlTestSuite{}

func (s *faucetsCtrlTestSuite) RegisterFaucet() usecases.RegisterFaucetUseCase {
	return s.registerFaucetUC
}

func (s *faucetsCtrlTestSuite) GetFaucet() usecases.GetFaucetUseCase {
	return s.getFaucetUC
}

func (s *faucetsCtrlTestSuite) SearchFaucets() usecases.SearchFaucetsUseCase {
	return s.searchFaucetUC
}

func (s *faucetsCtrlTestSuite) UpdateFaucet() usecases.UpdateFaucetUseCase {
	return s.updateFaucetUC
}

func (s *faucetsCtrlTestSuite) DeleteFaucet() usecases.DeleteFaucetUseCase {
	return s.deleteFaucetUC
}

func TestFaucetsController(t *testing.T) {
	s := new(faucetsCtrlTestSuite)
	suite.Run(t, s)
}

func (s *faucetsCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.tenants = []string{"tenantID"}
	s.registerFaucetUC = mocks.NewMockRegisterFaucetUseCase(ctrl)
	s.getFaucetUC = mocks.NewMockGetFaucetUseCase(ctrl)
	s.searchFaucetUC = mocks.NewMockSearchFaucetsUseCase(ctrl)
	s.updateFaucetUC = mocks.NewMockUpdateFaucetUseCase(ctrl)
	s.deleteFaucetUC = mocks.NewMockDeleteFaucetUseCase(ctrl)

	s.ctx = context.Background()
	s.ctx = context.WithValue(s.ctx, multitenancy.TenantIDKey, s.tenants[0])
	s.ctx = context.WithValue(s.ctx, multitenancy.AllowedTenantsKey, s.tenants)
	s.router = mux.NewRouter()

	controller := NewFaucetsController(s)
	controller.Append(s.router)
}

func (s *faucetsCtrlTestSuite) TestFaucetsController_Register() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		req := testutils.FakeRegisterFaucetRequest()
		requestBytes, _ := json.Marshal(req)
		faucet := testutils.FakeFaucet()
		rw := httptest.NewRecorder()

		httpRequest := httptest.
			NewRequest(http.MethodPost, endpoint, bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.registerFaucetUC.EXPECT().Execute(gomock.Any(), formatters.FormatRegisterFaucetRequest(req, s.tenants[0])).Return(faucet, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatFaucetResponse(faucet)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with Bad request if invalid format", func(t *testing.T) {
		req := testutils.FakeRegisterFaucetRequest()
		req.ChainRule = ""
		requestBytes, _ := json.Marshal(req)

		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPost, endpoint, bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	s.T().Run("should fail with 500 if use case fails with an unexpected error", func(t *testing.T) {
		req := testutils.FakeRegisterFaucetRequest()
		requestBytes, _ := json.Marshal(req)

		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPost, endpoint, bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.registerFaucetUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func (s *faucetsCtrlTestSuite) TestFaucetsController_GetOne() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		faucet := testutils.FakeFaucet()
		rw := httptest.NewRecorder()

		httpRequest := httptest.
			NewRequest(http.MethodGet, endpoint+"/faucetUUID", nil).
			WithContext(s.ctx)

		s.getFaucetUC.EXPECT().Execute(gomock.Any(), "faucetUUID", s.tenants).Return(faucet, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatFaucetResponse(faucet)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with Internal server error if use case fails", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodGet, endpoint+"/faucetUUID", nil).
			WithContext(s.ctx)

		s.getFaucetUC.EXPECT().Execute(gomock.Any(), "faucetUUID", s.tenants).Return(nil, fmt.Errorf("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func (s *faucetsCtrlTestSuite) TestFaucetsController_Update() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		req := testutils.FakeUpdateFaucetRequest()
		requestBytes, _ := json.Marshal(req)
		faucet := testutils.FakeFaucet()
		rw := httptest.NewRecorder()

		httpRequest := httptest.
			NewRequest(http.MethodPatch, endpoint+"/faucetUUID", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.updateFaucetUC.EXPECT().
			Execute(gomock.Any(), formatters.FormatUpdateFaucetRequest(req, "faucetUUID", s.tenants[0]), s.tenants).
			Return(faucet, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatFaucetResponse(faucet)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with Bad request if invalid format", func(t *testing.T) {
		req := testutils.FakeUpdateFaucetRequest()
		req.CreditorAccount = "notAnAddress"
		requestBytes, _ := json.Marshal(req)

		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPatch, endpoint+"/faucetUUID", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	s.T().Run("should fail with 500 if use case fails with an unexpected error", func(t *testing.T) {
		req := testutils.FakeUpdateFaucetRequest()
		requestBytes, _ := json.Marshal(req)

		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPatch, endpoint+"/faucetUUID", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.updateFaucetUC.EXPECT().Execute(gomock.Any(), gomock.Any(), s.tenants).Return(nil, fmt.Errorf("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func (s *faucetsCtrlTestSuite) TestFaucetsController_Search() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		faucet := testutils.FakeFaucet()
		rw := httptest.NewRecorder()
		names := []string{"name1", "name2"}
		chainRule := "chainRule"

		httpRequest := httptest.
			NewRequest(http.MethodGet, "/faucets?names="+strings.Join(names, ",")+"&chain_rule="+chainRule, nil).
			WithContext(s.ctx)

		expectedFilters := &entities.FaucetFilters{
			Names:     names,
			ChainRule: chainRule,
		}
		s.searchFaucetUC.EXPECT().Execute(gomock.Any(), expectedFilters, s.tenants).Return([]*entities.Faucet{faucet}, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatFaucetResponse(faucet)
		expectedBody, _ := json.Marshal([]*api.FaucetResponse{response})
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})
}

func (s *faucetsCtrlTestSuite) TestFaucetsController_Delete() {
	s.T().Run("should execute verify signature request successfully", func(t *testing.T) {
		acc := testutils.FakeAccount()
		acc.Address = inputTestAddress
		rw := httptest.NewRecorder()
		request := testutils.FakeVerifyPayloadRequest()
		requestBytes, _ := json.Marshal(request)

		httpRequest := httptest.
			NewRequest(http.MethodDelete, endpoint+"/faucetUUID", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.deleteFaucetUC.EXPECT().Execute(gomock.Any(), "faucetUUID", s.tenants).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, http.StatusNoContent, rw.Code)
	})
}
