// +build unit

package controllers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/identitymanager"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/identity-manager/use-cases"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/identity-manager/use-cases/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/service/formatters"
)

type identityCtrlTestSuite struct {
	suite.Suite
	createIdentityUC  *mocks.MockCreateAccountUseCase
	getIdentityUC     *mocks.MockGetAccountUseCase
	searchIdentityUC  *mocks.MockSearchAccountsUseCase
	updateIdentityUC  *mocks.MockUpdateAccountUseCase
	fundingIdentityUC *mocks.MockFundingAccountUseCase
	signPayloadUC     *mocks.MockSignPayloadUseCase
	ctx               context.Context
	tenants           []string
	router            *mux.Router
}

func (s *identityCtrlTestSuite) CreateAccount() usecases.CreateAccountUseCase {
	return s.createIdentityUC
}

func (s *identityCtrlTestSuite) GetAccount() usecases.GetAccountUseCase {
	return s.getIdentityUC
}

func (s *identityCtrlTestSuite) SearchAccounts() usecases.SearchAccountsUseCase {
	return s.searchIdentityUC
}

func (s *identityCtrlTestSuite) UpdateAccount() usecases.UpdateAccountUseCase {
	return s.updateIdentityUC
}

func (s *identityCtrlTestSuite) FundingAccount() usecases.FundingAccountUseCase {
	return s.fundingIdentityUC
}

func (s *identityCtrlTestSuite) SignPayload() usecases.SignPayloadUseCase {
	return s.signPayloadUC
}

var _ usecases.AccountUseCases = &identityCtrlTestSuite{}

const (
	inputTestAddress     = "0x7e654d251da770a068413677967f6d3ea2feA9e4"
	mixedCaseTestAddress = "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4"
)

func TestAccountController(t *testing.T) {
	s := new(identityCtrlTestSuite)
	suite.Run(t, s)
}

func (s *identityCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.tenants = []string{"tenantID"}
	s.createIdentityUC = mocks.NewMockCreateAccountUseCase(ctrl)
	s.getIdentityUC = mocks.NewMockGetAccountUseCase(ctrl)
	s.searchIdentityUC = mocks.NewMockSearchAccountsUseCase(ctrl)
	s.updateIdentityUC = mocks.NewMockUpdateAccountUseCase(ctrl)
	s.signPayloadUC = mocks.NewMockSignPayloadUseCase(ctrl)
	s.ctx = context.Background()
	s.ctx = context.WithValue(s.ctx, multitenancy.TenantIDKey, s.tenants[0])
	s.ctx = context.WithValue(s.ctx, multitenancy.AllowedTenantsKey, s.tenants)
	s.router = mux.NewRouter()

	controller := NewIdentitiesController(s)
	controller.Append(s.router)
}

func (s *identityCtrlTestSuite) TestAccountController_Create() {
	s.T().Run("should execute create account request successfully", func(t *testing.T) {
		req := testutils.FakeCreateAccountRequest()
		req.Chain = "besu"
		requestBytes, _ := json.Marshal(req)
		accResp := testutils.FakeAccount()
		rw := httptest.NewRecorder()

		httpRequest := httptest.
			NewRequest(http.MethodPost, "/accounts", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.createIdentityUC.EXPECT().Execute(gomock.Any(), gomock.Any(), "", req.Chain, s.tenants[0]).Return(accResp, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatAccountResponse(accResp)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with Bad request if invalid format", func(t *testing.T) {
		jobRequest := testutils.FakeImportAccountRequest()
		requestBytes, _ := json.Marshal(jobRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPost, "/accounts", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	s.T().Run("should fail with 500 if use case fails with an unexpected error", func(t *testing.T) {
		rw := httptest.NewRecorder()
		accRequest := testutils.FakeCreateAccountRequest()
		requestBytes, _ := json.Marshal(accRequest)
		httpRequest := httptest.
			NewRequest(http.MethodPost, "/accounts", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.createIdentityUC.EXPECT().Execute(gomock.Any(), gomock.Any(), "", "", s.tenants[0]).
			Return(nil, fmt.Errorf("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func (s *identityCtrlTestSuite) TestAccountController_Import() {
	s.T().Run("should execute import account request successfully", func(t *testing.T) {
		req := testutils.FakeImportAccountRequest()
		req.Chain = "qourum"
		requestBytes, _ := json.Marshal(req)
		accResp := testutils.FakeAccount()
		rw := httptest.NewRecorder()

		httpRequest := httptest.
			NewRequest(http.MethodPost, "/accounts/import", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.createIdentityUC.EXPECT().Execute(gomock.Any(), gomock.Any(), req.PrivateKey, req.Chain, s.tenants[0]).Return(accResp, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatAccountResponse(accResp)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})
}

func (s *identityCtrlTestSuite) TestAccountController_GetAccount() {
	s.T().Run("should execute get identity request successfully", func(t *testing.T) {
		accResp := testutils.FakeAccount()
		rw := httptest.NewRecorder()

		httpRequest := httptest.
			NewRequest(http.MethodGet, "/accounts/"+inputTestAddress, nil).
			WithContext(s.ctx)

		s.getIdentityUC.EXPECT().Execute(gomock.Any(), mixedCaseTestAddress, s.tenants).Return(accResp, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatAccountResponse(accResp)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with Bad request if invalid format", func(t *testing.T) {
		address := ethcommon.HexToAddress("0x123").String()

		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodGet, "/accounts/"+address, nil).
			WithContext(s.ctx)

		s.getIdentityUC.EXPECT().Execute(gomock.Any(), address, s.tenants).Return(nil, fmt.Errorf("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func (s *identityCtrlTestSuite) TestAccountController_UpdateAccount() {
	s.T().Run("should execute update account request successfully", func(t *testing.T) {
		req := testutils.FakeUpdateAccountRequest()
		rw := httptest.NewRecorder()
		requestBytes, _ := json.Marshal(req)

		httpRequest := httptest.
			NewRequest(http.MethodPatch, "/accounts/"+inputTestAddress, bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		acc := &entities.Account{
			Attributes: req.Attributes,
			Alias:      req.Alias,
			Address:    mixedCaseTestAddress,
		}

		s.updateIdentityUC.EXPECT().Execute(gomock.Any(), acc, s.tenants).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatAccountResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail to update account request if invalid request", func(t *testing.T) {
		rw := httptest.NewRecorder()
		address := ethcommon.HexToAddress("0x123").String()

		httpRequest := httptest.
			NewRequest(http.MethodPatch, "/accounts/"+address, nil).
			WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})
}

func (s *identityCtrlTestSuite) TestAccountController_SearchIdentity() {
	s.T().Run("should execute search identity request successfully", func(t *testing.T) {
		accResp := testutils.FakeAccount()
		rw := httptest.NewRecorder()
		aliases := []string{"alias1", "alias2"}

		httpRequest := httptest.
			NewRequest(http.MethodGet, "/accounts?aliases="+strings.Join(aliases, ","), nil).
			WithContext(s.ctx)

		filter := &entities.AccountFilters{
			Aliases: aliases,
		}

		s.searchIdentityUC.EXPECT().Execute(gomock.Any(), filter, s.tenants).Return([]*entities.Account{accResp}, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatAccountResponse(accResp)
		expectedBody, _ := json.Marshal([]*types.AccountResponse{response})
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})
}

func (s *identityCtrlTestSuite) TestAccountController_SignPayload() {
	s.T().Run("should execute sign payload request successfully", func(t *testing.T) {
		acc := testutils.FakeAccount()
		acc.Address = inputTestAddress
		rw := httptest.NewRecorder()
		payload := "payloadMessage"
		signedPayload := ethcommon.HexToHash("0xABCDEF").String()
		req := &types.SignPayloadRequest{
			Data: payload,
		}
		requestBytes, _ := json.Marshal(req)

		httpRequest := httptest.
			NewRequest(http.MethodPost, fmt.Sprintf("/accounts/%v/sign", acc.Address), bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.signPayloadUC.EXPECT().Execute(gomock.Any(), mixedCaseTestAddress, payload, s.tenants[0]).Return(signedPayload, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, signedPayload, rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})
}
