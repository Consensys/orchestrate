// +build unit

package gprc

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/contract-registry/use-cases/mocks"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
)

type testSuite struct {
	suite.Suite
	controller                  *ContractRegistry
	mockRegisterContractUseCase *mocks.MockRegisterContractUseCase
	mockGetUseCase              *mocks.MockGetContractUseCase
	mockGetMethodsUseCase       *mocks.MockGetMethodsUseCase
	mockGetEventsUseCase        *mocks.MockGetEventsUseCase
	mockGetCatalogUseCase       *mocks.MockGetCatalogUseCase
	mockGetTagsUseCase          *mocks.MockGetTagsUseCase
	mockSetCodeHashUseCase      *mocks.MockSetCodeHashUseCase
}

var errUseCase = fmt.Errorf("error")
var contract = testutils.FakeContract()

func TestContractRegistryController(t *testing.T) {
	s := new(testSuite)
	suite.Run(t, s)
}

func (s *testSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mockRegisterContractUseCase = mocks.NewMockRegisterContractUseCase(ctrl)
	s.mockGetUseCase = mocks.NewMockGetContractUseCase(ctrl)
	s.mockGetMethodsUseCase = mocks.NewMockGetMethodsUseCase(ctrl)
	s.mockGetEventsUseCase = mocks.NewMockGetEventsUseCase(ctrl)
	s.mockGetCatalogUseCase = mocks.NewMockGetCatalogUseCase(ctrl)
	s.mockGetTagsUseCase = mocks.NewMockGetTagsUseCase(ctrl)
	s.mockSetCodeHashUseCase = mocks.NewMockSetCodeHashUseCase(ctrl)

	s.controller = New(
		s.mockRegisterContractUseCase,
		s.mockGetUseCase,
		s.mockGetMethodsUseCase,
		s.mockGetEventsUseCase,
		s.mockGetCatalogUseCase,
		s.mockGetTagsUseCase,
		s.mockSetCodeHashUseCase,
	)
}

func (s *testSuite) TestContractRegistryController_RegisterContract() {
	request := &svc.RegisterContractRequest{Contract: contract}

	s.T().Run("should execute request successfully", func(t *testing.T) {
		s.mockRegisterContractUseCase.EXPECT().Execute(context.Background(), contract).Return(nil)

		response, err := s.controller.RegisterContract(context.Background(), request)

		assert.Equal(t, &svc.RegisterContractResponse{}, response)
		assert.Nil(t, err)
	})

	s.T().Run("should fail if use case fails", func(t *testing.T) {
		s.mockRegisterContractUseCase.EXPECT().Execute(context.Background(), contract).Return(errUseCase)

		response, err := s.controller.RegisterContract(context.Background(), request)

		assert.Nil(t, response, nil)
		assert.Equal(t, errors.FromError(errUseCase).ExtendComponent(component), err)
	})
}

func (s *testSuite) TestContractRegistryController_GetContract() {
	request := &svc.GetContractRequest{ContractId: &abi.ContractId{
		Name: contract.GetName(),
		Tag:  contract.GetTag(),
	}}

	s.T().Run("should execute request successfully", func(t *testing.T) {
		s.mockGetUseCase.EXPECT().Execute(context.Background(), request.GetContractId()).Return(contract, nil)

		response, err := s.controller.GetContract(context.Background(), request)

		assert.Equal(t, &svc.GetContractResponse{Contract: contract}, response)
		assert.Nil(t, err)
	})

	s.T().Run("should fail if use case fails", func(t *testing.T) {
		s.mockGetUseCase.EXPECT().Execute(context.Background(), request.GetContractId()).Return(nil, errUseCase)

		response, err := s.controller.GetContract(context.Background(), request)

		assert.Nil(t, response, nil)
		assert.Equal(t, errors.FromError(errUseCase).ExtendComponent(component), err)
	})
}

func (s *testSuite) TestContractRegistryController_GetContractABI() {
	request := &svc.GetContractRequest{ContractId: &abi.ContractId{
		Name: contract.GetName(),
		Tag:  contract.GetTag(),
	}}

	s.T().Run("should execute request successfully", func(t *testing.T) {
		s.mockGetUseCase.EXPECT().Execute(context.Background(), request.GetContractId()).Return(contract, nil)

		response, err := s.controller.GetContractABI(context.Background(), request)

		assert.Equal(t, &svc.GetContractABIResponse{Abi: contract.GetAbi()}, response)
		assert.Nil(t, err)
	})

	s.T().Run("should fail if use case fails", func(t *testing.T) {
		s.mockGetUseCase.EXPECT().Execute(context.Background(), request.GetContractId()).Return(nil, errUseCase)

		response, err := s.controller.GetContractABI(context.Background(), request)

		assert.Nil(t, response, nil)
		assert.Equal(t, errors.FromError(errUseCase).ExtendComponent(component), err)
	})
}

func (s *testSuite) TestContractRegistryController_GetContractBytecode() {
	request := &svc.GetContractRequest{ContractId: &abi.ContractId{
		Name: contract.GetName(),
		Tag:  contract.GetTag(),
	}}

	s.T().Run("should execute request successfully", func(t *testing.T) {
		s.mockGetUseCase.EXPECT().Execute(context.Background(), request.GetContractId()).Return(contract, nil)

		response, err := s.controller.GetContractBytecode(context.Background(), request)

		assert.Equal(t, &svc.GetContractBytecodeResponse{Bytecode: contract.GetBytecode()}, response)
		assert.Nil(t, err)
	})

	s.T().Run("should fail if use case fails", func(t *testing.T) {
		s.mockGetUseCase.EXPECT().Execute(context.Background(), request.GetContractId()).Return(nil, errUseCase)

		response, err := s.controller.GetContractBytecode(context.Background(), request)

		assert.Nil(t, response, nil)
		assert.Equal(t, errors.FromError(errUseCase).ExtendComponent(component), err)
	})
}

func (s *testSuite) TestContractRegistryController_GetContractDeployedBytecode() {
	request := &svc.GetContractRequest{ContractId: &abi.ContractId{
		Name: contract.GetName(),
		Tag:  contract.GetTag(),
	}}

	s.T().Run("should execute request successfully", func(t *testing.T) {
		s.mockGetUseCase.EXPECT().Execute(context.Background(), request.GetContractId()).Return(contract, nil)

		response, err := s.controller.GetContractDeployedBytecode(context.Background(), request)

		assert.Equal(t, &svc.GetContractDeployedBytecodeResponse{DeployedBytecode: contract.GetDeployedBytecode()}, response)
		assert.Nil(t, err)
	})

	s.T().Run("should fail if use case fails", func(t *testing.T) {
		s.mockGetUseCase.EXPECT().Execute(context.Background(), request.GetContractId()).Return(nil, errUseCase)

		response, err := s.controller.GetContractDeployedBytecode(context.Background(), request)

		assert.Nil(t, response, nil)
		assert.Equal(t, errors.FromError(errUseCase).ExtendComponent(component), err)
	})
}

func (s *testSuite) TestContractRegistryController_GetCatalog() {
	request := &svc.GetCatalogRequest{}

	s.T().Run("should execute request successfully", func(t *testing.T) {
		names := []string{"Contract0", "Contract1"}

		s.mockGetCatalogUseCase.EXPECT().Execute(context.Background()).Return(names, nil)

		response, err := s.controller.GetCatalog(context.Background(), request)

		assert.Equal(t, &svc.GetCatalogResponse{Names: names}, response)
		assert.Nil(t, err)
	})

	s.T().Run("should fail if use case fails", func(t *testing.T) {
		s.mockGetCatalogUseCase.EXPECT().Execute(context.Background()).Return(nil, errUseCase)

		response, err := s.controller.GetCatalog(context.Background(), request)

		assert.Nil(t, response, nil)
		assert.Equal(t, errors.FromError(errUseCase).ExtendComponent(component), err)
	})
}

func (s *testSuite) TestContractRegistryController_GetTags() {
	request := &svc.GetTagsRequest{Name: contract.GetName()}

	s.T().Run("should execute request successfully", func(t *testing.T) {
		tags := []string{"latest", "v1.0.0"}
		s.mockGetTagsUseCase.EXPECT().Execute(context.Background(), contract.GetName()).Return(tags, nil)

		response, err := s.controller.GetTags(context.Background(), request)

		assert.Equal(t, &svc.GetTagsResponse{Tags: tags}, response)
		assert.Nil(t, err)
	})

	s.T().Run("should fail if use case fails", func(t *testing.T) {
		s.mockGetTagsUseCase.EXPECT().Execute(context.Background(), contract.GetName()).Return(nil, errUseCase)

		response, err := s.controller.GetTags(context.Background(), request)

		assert.Nil(t, response, nil)
		assert.Equal(t, errors.FromError(errUseCase).ExtendComponent(component), err)
	})
}

func (s *testSuite) TestContractRegistryController_GetMethodsBySelector() {
	request := &svc.GetMethodsBySelectorRequest{
		AccountInstance: &common.AccountInstance{ChainId: "chainID"},
		Selector:        []byte{58, 58},
	}

	s.T().Run("should execute request successfully if ABI is not empty", func(t *testing.T) {
		responseABI := contract.GetAbi()

		s.mockGetMethodsUseCase.EXPECT().Execute(context.Background(), request.GetAccountInstance(), request.GetSelector()).Return(responseABI, nil, nil)

		response, err := s.controller.GetMethodsBySelector(context.Background(), request)

		assert.Equal(t, &svc.GetMethodsBySelectorResponse{Method: responseABI}, response)
		assert.Nil(t, err)
	})

	s.T().Run("should execute request successfully if ABI is empty", func(t *testing.T) {
		methodsABI := []string{"MethodABI0", "MethodABI1"}

		s.mockGetMethodsUseCase.EXPECT().Execute(context.Background(), request.GetAccountInstance(), request.GetSelector()).Return("", methodsABI, nil)

		response, err := s.controller.GetMethodsBySelector(context.Background(), request)

		assert.Equal(t, &svc.GetMethodsBySelectorResponse{DefaultMethods: methodsABI}, response)
		assert.Nil(t, err)
	})

	s.T().Run("should fail if use case fails", func(t *testing.T) {
		s.mockGetMethodsUseCase.EXPECT().Execute(context.Background(), request.GetAccountInstance(), request.GetSelector()).Return("", nil, errUseCase)

		response, err := s.controller.GetMethodsBySelector(context.Background(), request)

		assert.Nil(t, response, nil)
		assert.Equal(t, errors.FromError(errUseCase).ExtendComponent(component), err)
	})
}

func (s *testSuite) TestContractRegistryController_GetEventsBySigHash() {
	request := &svc.GetEventsBySigHashRequest{
		AccountInstance:   &common.AccountInstance{ChainId: "chainID"},
		SigHash:           "sigHash",
		IndexedInputCount: 35,
	}

	s.T().Run("should execute request successfully if ABI is not empty", func(t *testing.T) {
		responseABI := contract.GetAbi()

		s.mockGetEventsUseCase.EXPECT().Execute(
			context.Background(),
			request.GetAccountInstance(),
			request.GetSigHash(),
			request.GetIndexedInputCount(),
		).Return(responseABI, nil, nil)

		response, err := s.controller.GetEventsBySigHash(context.Background(), request)

		assert.Equal(t, &svc.GetEventsBySigHashResponse{Event: responseABI}, response)
		assert.Nil(t, err)
	})

	s.T().Run("should execute request successfully if ABI is empty", func(t *testing.T) {
		eventsABI := []string{"EventABI0", "EventABI1"}

		s.mockGetEventsUseCase.EXPECT().Execute(
			context.Background(),
			request.GetAccountInstance(),
			request.GetSigHash(),
			request.GetIndexedInputCount(),
		).Return("", eventsABI, nil)

		response, err := s.controller.GetEventsBySigHash(context.Background(), request)

		assert.Equal(t, &svc.GetEventsBySigHashResponse{DefaultEvents: eventsABI}, response)
		assert.Nil(t, err)
	})

	s.T().Run("should fail if use case fails", func(t *testing.T) {
		s.mockGetEventsUseCase.EXPECT().Execute(
			context.Background(),
			request.GetAccountInstance(),
			request.GetSigHash(),
			request.GetIndexedInputCount(),
		).Return("", nil, errUseCase)

		response, err := s.controller.GetEventsBySigHash(context.Background(), request)

		assert.Nil(t, response, nil)
		assert.Equal(t, errors.FromError(errUseCase).ExtendComponent(component), err)
	})
}

func (s *testSuite) TestContractRegistryController_SetAccountCodeHash() {
	request := &svc.SetAccountCodeHashRequest{
		AccountInstance: &common.AccountInstance{ChainId: "chainID"},
		CodeHash:        "codeHash",
	}

	s.T().Run("should execute request successfully", func(t *testing.T) {
		s.mockSetCodeHashUseCase.EXPECT().Execute(context.Background(), request.GetAccountInstance(), request.GetCodeHash()).Return(nil)

		response, err := s.controller.SetAccountCodeHash(context.Background(), request)

		assert.Equal(t, &svc.SetAccountCodeHashResponse{}, response)
		assert.Nil(t, err)
	})

	s.T().Run("should fail if use case fails", func(t *testing.T) {
		s.mockSetCodeHashUseCase.EXPECT().Execute(context.Background(), request.GetAccountInstance(), request.GetCodeHash()).Return(errUseCase)

		response, err := s.controller.SetAccountCodeHash(context.Background(), request)

		assert.Nil(t, response, nil)
		assert.Equal(t, errors.FromError(errUseCase).ExtendComponent(component), err)
	})
}

func (s *testSuite) TestContractRegistryController_DeleteArtifact() {
	s.T().Run("should return not implemented error", func(t *testing.T) {
		_, err := s.controller.DeleteArtifact(context.Background(), &svc.DeleteArtifactRequest{})

		assert.Equal(t, errors.FeatureNotSupportedError("DeleteArtifact not implemented yet").ExtendComponent(component), err)
	})
}

func (s *testSuite) TestContractRegistryController_DeregisterContract() {
	s.T().Run("should return not implemented error", func(t *testing.T) {
		_, err := s.controller.DeregisterContract(context.Background(), &svc.DeregisterContractRequest{})

		assert.Equal(t, errors.FeatureNotSupportedError("DeregisterContract not implemented yet").ExtendComponent(component), err)
	})
}
