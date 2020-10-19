// +build unit

package controllers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/use-cases"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/use-cases/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/service/formatters"
)

type identityCtrlTestSuite struct {
	suite.Suite
	createIdentityUC  *mocks.MockCreateIdentityUseCase
	searchIdentityUC  *mocks.MockSearchIdentitiesUseCase
	fundingIdentityUC *mocks.MockFundingIdentityUseCase
	ctx               context.Context
	tenants           []string
	router            *mux.Router
}

func (s *identityCtrlTestSuite) CreateIdentity() usecases.CreateIdentityUseCase {
	return s.createIdentityUC
}

func (s *identityCtrlTestSuite) SearchIdentity() usecases.SearchIdentitiesUseCase {
	return s.searchIdentityUC
}
func (s *identityCtrlTestSuite) FundingIdentity() usecases.FundingIdentityUseCase {
	return s.fundingIdentityUC
}

var _ usecases.IdentityUseCases = &identityCtrlTestSuite{}

func TestJobsController(t *testing.T) {
	s := new(identityCtrlTestSuite)
	suite.Run(t, s)
}

func (s *identityCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.tenants = []string{"tenantID"}
	s.createIdentityUC = mocks.NewMockCreateIdentityUseCase(ctrl)
	s.searchIdentityUC = mocks.NewMockSearchIdentitiesUseCase(ctrl)
	s.ctx = context.Background()
	s.ctx = context.WithValue(s.ctx, multitenancy.TenantIDKey, s.tenants[0])
	s.ctx = context.WithValue(s.ctx, multitenancy.AllowedTenantsKey, s.tenants)
	s.router = mux.NewRouter()

	controller := NewIdentitiesController(s)
	controller.Append(s.router)
}

func (s *identityCtrlTestSuite) TestJobsController_Create() {
	s.T().Run("should execute create identity request successfully", func(t *testing.T) {
		req := testutils.FakeCreateIdentityRequest()
		req.Chain = "besu"
		requestBytes, _ := json.Marshal(req)
		idenResp := testutils.FakeIdentity()
		rw := httptest.NewRecorder()

		httpRequest := httptest.
			NewRequest(http.MethodPost, "/identities", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.createIdentityUC.EXPECT().Execute(gomock.Any(), gomock.Any(), "", req.Chain, s.tenants[0]).Return(idenResp, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatIdentityResponse(idenResp)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should execute import identity request successfully", func(t *testing.T) {
		req := testutils.FakeImportIdentityRequest()
		req.Chain = "qourum"
		requestBytes, _ := json.Marshal(req)
		idenResp := testutils.FakeIdentity()
		rw := httptest.NewRecorder()

		httpRequest := httptest.
			NewRequest(http.MethodPost, "/identities/import", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.createIdentityUC.EXPECT().Execute(gomock.Any(), gomock.Any(), req.PrivateKey, req.Chain, s.tenants[0]).Return(idenResp, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatIdentityResponse(idenResp)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with Bad request if invalid format", func(t *testing.T) {
		jobRequest := testutils.FakeCreateIdentityRequest()
		jobRequest.Alias = ""
		requestBytes, _ := json.Marshal(jobRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPost, "/identities", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	s.T().Run("should fail with 422 if use case fails with InvalidParameterError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		jobRequest := testutils.FakeCreateIdentityRequest()
		requestBytes, _ := json.Marshal(jobRequest)
		httpRequest := httptest.
			NewRequest(http.MethodPost, "/identities", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.createIdentityUC.EXPECT().Execute(gomock.Any(), gomock.Any(), "", "", s.tenants[0]).Return(nil, errors.InvalidParameterError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})
}
