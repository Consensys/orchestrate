package controllers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/consensys/orchestrate/pkg/encoding/json"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/pkg/utils"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/business/use-cases/mocks"
	"github.com/consensys/orchestrate/services/api/service/formatters"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type contractsCtrlTestSuite struct {
	suite.Suite
	getContractsCatalog         *mocks.MockGetContractsCatalogUseCase
	getContract                 *mocks.MockGetContractUseCase
	getContractEvents           *mocks.MockGetContractEventsUseCase
	getContractMethodSignatures *mocks.MockGetContractMethodSignaturesUseCase
	getContractMethods          *mocks.MockGetContractMethodsUseCase
	getContractTags             *mocks.MockGetContractTagsUseCase
	setContractCodeHash         *mocks.MockSetContractCodeHashUseCase
	registerContract            *mocks.MockRegisterContractUseCase
	router                      *mux.Router
}

var _ usecases.ContractUseCases = &contractsCtrlTestSuite{}

func (s *contractsCtrlTestSuite) GetContractsCatalog() usecases.GetContractsCatalogUseCase {
	return s.getContractsCatalog
}
func (s *contractsCtrlTestSuite) GetContract() usecases.GetContractUseCase {
	return s.getContract
}
func (s *contractsCtrlTestSuite) GetContractEvents() usecases.GetContractEventsUseCase {
	return s.getContractEvents
}
func (s *contractsCtrlTestSuite) GetContractMethodSignatures() usecases.GetContractMethodSignaturesUseCase {
	return s.getContractMethodSignatures
}
func (s *contractsCtrlTestSuite) GetContractMethods() usecases.GetContractMethodsUseCase {
	return s.getContractMethods
}
func (s *contractsCtrlTestSuite) GetContractTags() usecases.GetContractTagsUseCase {
	return s.getContractTags
}
func (s *contractsCtrlTestSuite) SetContractCodeHash() usecases.SetContractCodeHashUseCase {
	return s.setContractCodeHash
}
func (s *contractsCtrlTestSuite) RegisterContract() usecases.RegisterContractUseCase {
	return s.registerContract
}

func TestContractController(t *testing.T) {
	s := new(contractsCtrlTestSuite)
	suite.Run(t, s)
}

func (s *contractsCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.getContractsCatalog = mocks.NewMockGetContractsCatalogUseCase(ctrl)
	s.getContract = mocks.NewMockGetContractUseCase(ctrl)
	s.getContractEvents = mocks.NewMockGetContractEventsUseCase(ctrl)
	s.getContractMethodSignatures = mocks.NewMockGetContractMethodSignaturesUseCase(ctrl)
	s.getContractMethods = mocks.NewMockGetContractMethodsUseCase(ctrl)
	s.getContractTags = mocks.NewMockGetContractTagsUseCase(ctrl)
	s.setContractCodeHash = mocks.NewMockSetContractCodeHashUseCase(ctrl)
	s.registerContract = mocks.NewMockRegisterContractUseCase(ctrl)
	s.router = mux.NewRouter()

	controller := NewContractsController(s)
	controller.Append(s.router)
}

func (s *contractsCtrlTestSuite) TestContractsController_Register() {
	ctx := context.Background()
	s.T().Run("should execute register contract request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := testutils.FakeRegisterContractRequest()
		requestBytes, _ := json.Marshal(req)
		httpRequest := httptest.
			NewRequest(http.MethodPost, "/contracts", bytes.NewReader(requestBytes)).
			WithContext(ctx)

		expectedContract, _ := formatters.FormatRegisterContractRequest(req)
		s.registerContract.EXPECT().Execute(gomock.Any(), expectedContract).Return(nil)

		contract := testutils.FakeContract()
		s.getContract.EXPECT().Execute(gomock.Any(), req.Name, req.Tag).Return(contract, nil)

		s.router.ServeHTTP(rw, httpRequest)
		expectedBody, _ := json.Marshal(formatters.FormatContractResponse(contract))
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
	})

	s.T().Run("should fail to register contract request if invalid format", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := testutils.FakeRegisterContractRequest()
		req.Name = ""
		requestBytes, _ := json.Marshal(req)
		httpRequest := httptest.
			NewRequest(http.MethodPost, "/contracts", bytes.NewReader(requestBytes)).
			WithContext(ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	s.T().Run("should fail to register contract if register contract fails", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := testutils.FakeRegisterContractRequest()
		requestBytes, _ := json.Marshal(req)
		httpRequest := httptest.
			NewRequest(http.MethodPost, "/contracts", bytes.NewReader(requestBytes)).
			WithContext(ctx)

		expectedContract, _ := formatters.FormatRegisterContractRequest(req)
		s.registerContract.EXPECT().Execute(gomock.Any(), expectedContract).Return(fmt.Errorf("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})

	s.T().Run("should fail to register contract if get contract fails", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := testutils.FakeRegisterContractRequest()
		requestBytes, _ := json.Marshal(req)
		httpRequest := httptest.
			NewRequest(http.MethodPost, "/contracts", bytes.NewReader(requestBytes)).
			WithContext(ctx)

		expectedContract, _ := formatters.FormatRegisterContractRequest(req)
		s.registerContract.EXPECT().Execute(gomock.Any(), expectedContract).Return(nil)

		s.getContract.EXPECT().Execute(gomock.Any(), req.Name, req.Tag).Return(nil, fmt.Errorf("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func (s *contractsCtrlTestSuite) TestContractsController_CodeHash() {
	ctx := context.Background()
	chainID := "2017"
	address := testutils.FakeAddress()

	s.T().Run("should execute set contract codeHash successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := testutils.FakeSetContractCodeHashRequest()
		requestBytes, _ := json.Marshal(req)
		httpRequest := httptest.
			NewRequest(http.MethodPost, fmt.Sprintf("/contracts/accounts/%s/%s", chainID, address.Hex()), bytes.NewReader(requestBytes)).
			WithContext(ctx)

		s.setContractCodeHash.EXPECT().
			Execute(gomock.Any(), chainID, address, req.CodeHash).
			Return(nil)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, "OK", rw.Body.String())
	})

	s.T().Run("should fail set contract codeHash if address is not valid", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := testutils.FakeSetContractCodeHashRequest()
		address2 := "invalid_address_2"
		requestBytes, _ := json.Marshal(req)
		httpRequest := httptest.
			NewRequest(http.MethodPost, fmt.Sprintf("/contracts/accounts/%s/%s", chainID, address2), bytes.NewReader(requestBytes)).
			WithContext(ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	s.T().Run("should fail set contract codeHash if set contract usecase fails", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := testutils.FakeSetContractCodeHashRequest()
		requestBytes, _ := json.Marshal(req)
		httpRequest := httptest.
			NewRequest(http.MethodPost, fmt.Sprintf("/contracts/accounts/%s/%s", chainID, address), bytes.NewReader(requestBytes)).
			WithContext(ctx)

		s.setContractCodeHash.EXPECT().
			Execute(gomock.Any(), chainID, address, req.CodeHash).
			Return(fmt.Errorf("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func (s *contractsCtrlTestSuite) TestContractsController_GetContract() {
	ctx := context.Background()

	s.T().Run("should execute get contract successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		contract := testutils.FakeContract()
		httpRequest := httptest.
			NewRequest(http.MethodGet, fmt.Sprintf("/contracts/%s/%s", contract.Name, contract.Tag), nil).
			WithContext(ctx)

		s.getContract.EXPECT().Execute(gomock.Any(), contract.Name, contract.Tag).Return(contract, nil)

		s.router.ServeHTTP(rw, httpRequest)
		expectedBody, _ := json.Marshal(formatters.FormatContractResponse(contract))
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
	})

	s.T().Run("should fail to get contract if usecase fails", func(t *testing.T) {
		rw := httptest.NewRecorder()
		contract := testutils.FakeContract()
		httpRequest := httptest.
			NewRequest(http.MethodGet, fmt.Sprintf("/contracts/%s/%s", contract.Name, contract.Tag), nil).
			WithContext(ctx)

		s.getContract.EXPECT().Execute(gomock.Any(), contract.Name, contract.Tag).Return(nil, fmt.Errorf("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func (s *contractsCtrlTestSuite) TestContractsController_GetContractEvents() {
	ctx := context.Background()
	address := ethcommon.HexToAddress(utils.RandHexString(10))
	sigHash := ethcommon.HexToHash(utils.RandHexString(10)).String()
	indexInput := uint32(2)
	chainID := "2017"

	event := testutils.FakeEventABI()
	defaultEvent := testutils.FakeEventABI()
	rawEvent, _ := json.Marshal(event)
	rawDefaultEvent, _ := json.Marshal(defaultEvent)

	s.T().Run("should execute get contract events successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodGet, fmt.Sprintf("/contracts/accounts/%s/%s/events?sig_hash=%s&indexed_input_count=%d",
				chainID, address.Hex(), sigHash, indexInput), nil).
			WithContext(ctx)

		s.getContractEvents.EXPECT().Execute(gomock.Any(), chainID, address, sigHash, indexInput).Return(string(rawEvent), []string{string(rawDefaultEvent)}, nil)

		s.router.ServeHTTP(rw, httpRequest)
		expectedBody, _ := json.Marshal(api.GetContractEventsBySignHashResponse{Event: string(rawEvent), DefaultEvents: []string{string(rawDefaultEvent)}})
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
	})

	s.T().Run("should failt get contract events if address is invalid", func(t *testing.T) {
		invalidAddr := "invalid_address"
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodGet, fmt.Sprintf("/contracts/accounts/%s/%s/events?sig_hash=%s&indexed_input_count=%d",
				chainID, invalidAddr, sigHash, indexInput), nil).
			WithContext(ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})
}

func (s *contractsCtrlTestSuite) TestContractsController_GetContractsCatalog() {
	ctx := context.Background()

	s.T().Run("should execute get catalog successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodGet, "/contracts", nil).
			WithContext(ctx)

		catalog := []string{"contractOne", "contractTwo"}
		s.getContractsCatalog.EXPECT().Execute(gomock.Any()).Return(catalog, nil)

		s.router.ServeHTTP(rw, httpRequest)
		expectedBody, _ := json.Marshal(catalog)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
	})
}

func (s *contractsCtrlTestSuite) TestContractsController_GetContractMethodSignatures() {
	ctx := context.Background()

	s.T().Run("should execute get contract method signatures successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		contract := testutils.FakeContract()
		httpRequest := httptest.
			NewRequest(http.MethodGet, fmt.Sprintf("/contracts/%s/%s/method-signatures", contract.Name, contract.Tag), nil).
			WithContext(ctx)

		methodSignatures := []string{"method1()", "method2()"}
		s.getContractMethodSignatures.EXPECT().
			Execute(gomock.Any(), contract.Name, contract.Tag, "").
			Return(methodSignatures, nil)

		s.router.ServeHTTP(rw, httpRequest)
		expectedBody, _ := json.Marshal(methodSignatures)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
	})

	s.T().Run("should execute get contract method signatures with filter successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		contract := testutils.FakeContract()
		method := "method1"
		httpRequest := httptest.
			NewRequest(http.MethodGet, fmt.Sprintf("/contracts/%s/%s/method-signatures?method=%s", contract.Name, contract.Tag, method), nil).
			WithContext(ctx)

		methodSignatures := []string{"method1()"}
		s.getContractMethodSignatures.EXPECT().
			Execute(gomock.Any(), contract.Name, contract.Tag, method).
			Return(methodSignatures, nil)

		s.router.ServeHTTP(rw, httpRequest)
		expectedBody, _ := json.Marshal(methodSignatures)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
	})
}
