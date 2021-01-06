package controllers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/service/formatters"
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
		s.getContract.EXPECT().Execute(gomock.Any(), &entities.ContractID{
			Name: req.Name,
			Tag:  req.Tag,
		}).Return(contract, nil)

		s.router.ServeHTTP(rw, httpRequest)
		expectedBody, _ := json.Marshal(api.ContractResponse{Contract: contract})
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

		s.getContract.EXPECT().Execute(gomock.Any(), &entities.ContractID{
			Name: req.Name,
			Tag:  req.Tag,
		}).Return(nil, fmt.Errorf("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func (s *contractsCtrlTestSuite) TestContractsController_CodeHash() {
	ctx := context.Background()

	s.T().Run("should execute set contract codeHash successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := testutils.FakeSetContractCodeHashRequest()
		requestBytes, _ := json.Marshal(req)
		httpRequest := httptest.
			NewRequest(http.MethodPatch, "/contracts", bytes.NewReader(requestBytes)).
			WithContext(ctx)

		s.setContractCodeHash.EXPECT().
			Execute(gomock.Any(), req.ChainID, req.Address, req.CodeHash).
			Return(nil)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, "OK", rw.Body.String())
	})

	s.T().Run("should fail set contract codeHash if address is not valid", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := testutils.FakeSetContractCodeHashRequest()
		req.Address = "invalid_address_2"
		requestBytes, _ := json.Marshal(req)
		httpRequest := httptest.
			NewRequest(http.MethodPatch, "/contracts", bytes.NewReader(requestBytes)).
			WithContext(ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	s.T().Run("should fail set contract codeHash if set contract usecase fails", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := testutils.FakeSetContractCodeHashRequest()
		requestBytes, _ := json.Marshal(req)
		httpRequest := httptest.
			NewRequest(http.MethodPatch, "/contracts", bytes.NewReader(requestBytes)).
			WithContext(ctx)

		s.setContractCodeHash.EXPECT().
			Execute(gomock.Any(), req.ChainID, req.Address, req.CodeHash).
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
			NewRequest(http.MethodGet, fmt.Sprintf("/contracts/%s/%s", contract.ID.Name, contract.ID.Tag), nil).
			WithContext(ctx)

		s.getContract.EXPECT().
			Execute(gomock.Any(), &contract.ID).
			Return(contract, nil)

		s.router.ServeHTTP(rw, httpRequest)
		expectedBody, _ := json.Marshal(contract)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
	})

	s.T().Run("should fail to get contract if usecase fails", func(t *testing.T) {
		rw := httptest.NewRecorder()
		contract := testutils.FakeContract()
		httpRequest := httptest.
			NewRequest(http.MethodGet, fmt.Sprintf("/contracts/%s/%s", contract.ID.Name, contract.ID.Tag), nil).
			WithContext(ctx)

		s.getContract.EXPECT().
			Execute(gomock.Any(), &contract.ID).
			Return(nil, fmt.Errorf("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func (s *contractsCtrlTestSuite) TestContractsController_GetContractEvents() {
	ctx := context.Background()
	address := ethcommon.HexToAddress(utils.RandHexString(10)).String()
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
			NewRequest(http.MethodGet, fmt.Sprintf("/contracts/%s/events?sig_hash=%s&chain_id=%s&indexed_input_count=%d",
				address, sigHash, chainID, indexInput), nil).
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
			NewRequest(http.MethodGet, fmt.Sprintf("/contracts/%s/events?sig_hash=%s&chain_id=%s&indexed_input_count=%d",
				invalidAddr, sigHash, chainID, indexInput), nil).
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
			NewRequest(http.MethodGet, fmt.Sprintf("/contracts/%s/%s/method-signatures", contract.ID.Name, contract.ID.Tag), nil).
			WithContext(ctx)

		methodSignatures := []string{"method1()", "method2()"}
		s.getContractMethodSignatures.EXPECT().
			Execute(gomock.Any(), &contract.ID, "").
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
			NewRequest(http.MethodGet, fmt.Sprintf("/contracts/%s/%s/method-signatures?method=%s", contract.ID.Name, contract.ID.Tag, method), nil).
			WithContext(ctx)

		methodSignatures := []string{"method1()"}
		s.getContractMethodSignatures.EXPECT().
			Execute(gomock.Any(), &contract.ID, method).
			Return(methodSignatures, nil)

		s.router.ServeHTTP(rw, httpRequest)
		expectedBody, _ := json.Marshal(methodSignatures)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
	})
}
