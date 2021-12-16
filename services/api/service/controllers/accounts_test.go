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

	qkm "github.com/consensys/orchestrate/pkg/quorum-key-manager"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/service/formatters"
	qkmmock "github.com/consensys/quorum-key-manager/pkg/client/mock"
	qkmtypes "github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/consensys/orchestrate/pkg/encoding/json"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/services/api/business/use-cases/mocks"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type accountsCtrlTestSuite struct {
	suite.Suite
	createAccountUC  *mocks.MockCreateAccountUseCase
	getAccountUC     *mocks.MockGetAccountUseCase
	searchAccountUC  *mocks.MockSearchAccountsUseCase
	updateAccountUC  *mocks.MockUpdateAccountUseCase
	fundAccountUC    *mocks.MockFundAccountUseCase
	keyManagerClient *qkmmock.MockKeyManagerClient
	ctx              context.Context
	userInfo         *multitenancy.UserInfo
	router           *mux.Router
}

var _ usecases.AccountUseCases = &accountsCtrlTestSuite{}

func (s *accountsCtrlTestSuite) CreateAccount() usecases.CreateAccountUseCase {
	return s.createAccountUC
}

func (s *accountsCtrlTestSuite) GetAccount() usecases.GetAccountUseCase {
	return s.getAccountUC
}

func (s *accountsCtrlTestSuite) SearchAccounts() usecases.SearchAccountsUseCase {
	return s.searchAccountUC
}

func (s *accountsCtrlTestSuite) UpdateAccount() usecases.UpdateAccountUseCase {
	return s.updateAccountUC
}

func (s *accountsCtrlTestSuite) FundAccount() usecases.FundAccountUseCase {
	return s.fundAccountUC
}

const (
	inputTestAddress     = "0x7e654d251da770a068413677967f6d3ea2feA9e4"
	mixedCaseTestAddress = "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4"
	qkmStoreName         = "test-store-name"
)

func TestAccountController(t *testing.T) {
	s := new(accountsCtrlTestSuite)
	suite.Run(t, s)
}

func (s *accountsCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.createAccountUC = mocks.NewMockCreateAccountUseCase(ctrl)
	s.getAccountUC = mocks.NewMockGetAccountUseCase(ctrl)
	s.searchAccountUC = mocks.NewMockSearchAccountsUseCase(ctrl)
	s.updateAccountUC = mocks.NewMockUpdateAccountUseCase(ctrl)
	s.keyManagerClient = qkmmock.NewMockKeyManagerClient(ctrl)
	s.userInfo = multitenancy.NewUserInfo("tenantOne", "username")
	s.ctx = multitenancy.WithUserInfo(context.Background(), s.userInfo)
	s.router = mux.NewRouter()

	controller := NewAccountsController(s, s.keyManagerClient, qkmStoreName)
	controller.Append(s.router)
}

func (s *accountsCtrlTestSuite) TestAccountController_Create() {
	s.T().Run("should execute create account request successfully", func(t *testing.T) {
		req := testutils.FakeCreateAccountRequest()
		req.Chain = "besu"
		requestBytes, _ := json.Marshal(req)
		accResp := testutils.FakeAccount()
		rw := httptest.NewRecorder()

		httpRequest := httptest.
			NewRequest(http.MethodPost, "/accounts", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.createAccountUC.EXPECT().Execute(gomock.Any(), gomock.Any(), nil, req.Chain, s.userInfo).Return(accResp, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatAccountResponse(accResp)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with Bad request if invalid format", func(t *testing.T) {
		req := testutils.FakeImportAccountRequest()
		requestBytes, _ := json.Marshal(req)

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

		s.createAccountUC.EXPECT().Execute(gomock.Any(), gomock.Any(), nil, "", s.userInfo).
			Return(nil, fmt.Errorf("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func (s *accountsCtrlTestSuite) TestAccountController_Import() {
	s.T().Run("should execute import account request successfully", func(t *testing.T) {
		req := testutils.FakeImportAccountRequest()
		req.Chain = "qourum"
		requestBytes, _ := json.Marshal(req)
		accResp := testutils.FakeAccount()
		rw := httptest.NewRecorder()

		httpRequest := httptest.
			NewRequest(http.MethodPost, "/accounts/import", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.createAccountUC.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any(), req.Chain, s.userInfo).Return(accResp, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatAccountResponse(accResp)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})
}

func (s *accountsCtrlTestSuite) TestAccountController_GetOne() {
	s.T().Run("should execute get account request successfully", func(t *testing.T) {
		accResp := testutils.FakeAccount()
		rw := httptest.NewRecorder()

		httpRequest := httptest.
			NewRequest(http.MethodGet, "/accounts/"+inputTestAddress, nil).
			WithContext(s.ctx)

		s.getAccountUC.EXPECT().Execute(gomock.Any(), ethcommon.HexToAddress(mixedCaseTestAddress), s.userInfo).Return(accResp, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatAccountResponse(accResp)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with Internal server error if use case fails", func(t *testing.T) {
		address := ethcommon.HexToAddress("0x123")

		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodGet, "/accounts/"+address.Hex(), nil).
			WithContext(s.ctx)

		s.getAccountUC.EXPECT().Execute(gomock.Any(), address, s.userInfo).Return(nil, fmt.Errorf("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func (s *accountsCtrlTestSuite) TestAccountController_UpdateAccount() {
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
			StoreID:    req.StoreID,
			Address:    ethcommon.HexToAddress(mixedCaseTestAddress),
		}

		s.updateAccountUC.EXPECT().Execute(gomock.Any(), acc, s.userInfo).Return(acc, nil)

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

func (s *accountsCtrlTestSuite) TestAccountController_SearchIdentity() {
	s.T().Run("should execute search account request successfully", func(t *testing.T) {
		accResp := testutils.FakeAccount()
		rw := httptest.NewRecorder()
		aliases := []string{"alias1", "alias2"}

		httpRequest := httptest.
			NewRequest(http.MethodGet, "/accounts?aliases="+strings.Join(aliases, ","), nil).
			WithContext(s.ctx)

		filter := &entities.AccountFilters{
			Aliases: aliases,
		}

		s.searchAccountUC.EXPECT().Execute(gomock.Any(), filter, s.userInfo).Return([]*entities.Account{accResp}, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatAccountResponse(accResp)
		expectedBody, _ := json.Marshal([]*api.AccountResponse{response})
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})
}

func (s *accountsCtrlTestSuite) TestAccountController_SignPayload() {
	s.T().Run("should execute sign payload request successfully", func(t *testing.T) {
		acc := testutils.FakeAccount()
		acc.Address = ethcommon.HexToAddress(inputTestAddress)
		rw := httptest.NewRecorder()
		payload := hexutil.MustDecode("0x1234")
		signature := "0xsignature"
		requestBytes, _ := json.Marshal(&api.SignMessageRequest{qkmtypes.SignMessageRequest{Message: payload}, acc.StoreID})

		httpRequest := httptest.
			NewRequest(http.MethodPost, fmt.Sprintf("/accounts/%v/sign-message", acc.Address), bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.getAccountUC.EXPECT().Execute(gomock.Any(), ethcommon.HexToAddress(mixedCaseTestAddress), s.userInfo).Return(acc, nil)
		s.keyManagerClient.EXPECT().SignMessage(gomock.Any(), acc.StoreID, mixedCaseTestAddress, &qkmtypes.SignMessageRequest{
			Message: payload,
		}).Return(signature, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, signature, rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})
}

func (s *accountsCtrlTestSuite) TestAccountController_VerifySignature() {
	s.T().Run("should execute verify signature request successfully", func(t *testing.T) {
		acc := testutils.FakeAccount()
		acc.Address = ethcommon.HexToAddress(inputTestAddress)
		rw := httptest.NewRecorder()
		request := qkm.FakeVerifyPayloadRequest()
		requestBytes, _ := json.Marshal(request)

		httpRequest := httptest.
			NewRequest(http.MethodPost, "/accounts/verify-message", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.keyManagerClient.EXPECT().VerifyMessage(gomock.Any(), request).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, http.StatusNoContent, rw.Code)
	})
}

func (s *accountsCtrlTestSuite) TestAccountController_VerifyTypedDataSignature() {
	s.T().Run("should execute verify typed data signature request successfully", func(t *testing.T) {
		acc := testutils.FakeAccount()
		acc.Address = ethcommon.HexToAddress(inputTestAddress)
		rw := httptest.NewRecorder()
		request := qkm.FakeVerifyTypedDataPayloadRequest()
		requestBytes, _ := json.Marshal(request)

		httpRequest := httptest.
			NewRequest(http.MethodPost, "/accounts/verify-typed-data", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.keyManagerClient.EXPECT().VerifyTypedData(gomock.Any(), request).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, http.StatusNoContent, rw.Code)
	})
}
