// +build integration

package integrationtests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client/dialer"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
)

type contractRegistryGRPCTestSuite struct {
	suite.Suite
	baseURL    string
	grpcClient registry.ContractRegistryClient
	env        *IntegrationEnvironment
}

func (s *contractRegistryGRPCTestSuite) SetupSuite() {
	client, err := dialer.DialContextWithDefaultOptions(context.Background(), s.baseURL)
	if err != nil {
		panic(err)
	}
	s.grpcClient = client
}

func (s *contractRegistryGRPCTestSuite) TestContractRegistry_Validation() {
	s.T().Run("should fail with X if payload is invalid", func(t *testing.T) {
		contract := testutils.FakeContract()
		contract.SetName("")

		_, err := s.registerContract(contract)
		assert.True(t, errors.IsDataError(err))
		assert.Equal(t, errors.FromError(err).GetCode(), errors.InvalidArg)
		assert.Equal(t, errors.FromError(err).GetMessage(), "No name provided in request")
	})

	s.T().Run("should not fail if contract registered twice", func(t *testing.T) {
		contract := testutils.FakeContract()

		_, err := s.registerContract(contract)
		assert.NoError(t, err)

		_, err = s.registerContract(contract)
		assert.NoError(t, err)
	})
}

func (s *contractRegistryGRPCTestSuite) TestContractRegistry_Register() {
	ctx := context.Background()

	s.T().Run("should register a contract with tag", func(t *testing.T) {
		contract := testutils.FakeContract()
		contract.SetTag("tag")

		_, _ = s.registerContract(contract)
		resp, err := s.grpcClient.GetContract(ctx, &registry.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: contract.GetName(),
				Tag:  contract.GetTag(),
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, resp.GetContract().GetName(), contract.GetName())
		assert.Equal(t, resp.GetContract().GetTag(), contract.GetTag())
	})

	s.T().Run("should register a contract with tag latest", func(t *testing.T) {
		contract := testutils.FakeContract()
		contract.SetTag("")

		_, _ = s.registerContract(contract)
		resp, err := s.grpcClient.GetContract(ctx, &registry.GetContractRequest{
			ContractId: &abi.ContractId{
				Name: contract.GetName(),
				Tag:  contract.GetTag(),
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, resp.GetContract().GetName(), contract.GetName())
		assert.Equal(t, resp.GetContract().GetTag(), "latest")
	})
}

func (s *contractRegistryGRPCTestSuite) TestContractRegistry_Get() {
	ctx := context.Background()
	contract0 := testutils.FakeContract()
	_ = contract0.CompactABI()
	_, _ = s.registerContract(contract0)

	contract1 := testutils.FakeContract()
	_, _ = s.registerContract(contract1)

	s.T().Run("should get all contracts", func(t *testing.T) {
		_, err := s.grpcClient.GetCatalog(ctx, &registry.GetCatalogRequest{})
		assert.NoError(t, err)
	})

	s.T().Run("should get all tags of a contract", func(t *testing.T) {
		resp, err := s.grpcClient.GetTags(ctx, &registry.GetTagsRequest{
			Name: contract0.GetName(),
		})
		assert.NoError(t, err)
		assert.Equal(t, resp.GetTags(), []string{"v1.0.0"})
	})

	s.T().Run("should get a contract", func(t *testing.T) {
		resp, err := s.grpcClient.GetContract(ctx, s.getContractRequest(contract0.GetName(), contract0.GetTag()))

		assert.NoError(t, err)
		assert.Equal(t, resp.GetContract().GetName(), contract0.GetName())
		assert.Equal(t, resp.GetContract().GetTag(), contract0.GetTag())
		assert.Equal(t, resp.GetContract().GetBytecode(), contract0.GetBytecode())
		assert.Equal(t, resp.GetContract().GetDeployedBytecode(), contract0.GetDeployedBytecode())
		assert.Equal(t, resp.GetContract().GetAbi(), contract0.GetAbi())
	})

	s.T().Run("should get a contract abi", func(t *testing.T) {
		resp, err := s.grpcClient.GetContractABI(ctx, s.getContractRequest(contract0.GetName(), contract0.GetTag()))
		assert.NoError(t, err)
		assert.Equal(t, resp.GetAbi(), contract0.GetAbi())
	})

	s.T().Run("should get a contract bytecode", func(t *testing.T) {
		resp, err := s.grpcClient.GetContractBytecode(ctx, s.getContractRequest(contract0.GetName(), contract0.GetTag()))
		assert.NoError(t, err)
		assert.Equal(t, resp.GetBytecode(), contract0.GetBytecode())
	})

	s.T().Run("should get a contract deployed bytecode", func(t *testing.T) {
		resp, err := s.grpcClient.GetContractDeployedBytecode(ctx, s.getContractRequest(contract0.GetName(), contract0.GetTag()))
		assert.NoError(t, err)
		assert.Equal(t, resp.GetDeployedBytecode(), contract0.GetDeployedBytecode())
	})

	s.T().Run("should get a contract constructor and contract methods", func(t *testing.T) {
		resp0, err := s.grpcClient.GetMethodSignatures(ctx, &registry.GetMethodSignaturesRequest{
			ContractId: &abi.ContractId{
				Name: contract0.GetName(),
				Tag:  contract0.GetTag(),
			},
			MethodName: "constructor",
		})

		assert.NoError(t, err)
		assert.Equal(t, resp0.GetSignatures()[0], "constructor(uint256)")

		resp1, err := s.grpcClient.GetMethodSignatures(ctx, &registry.GetMethodSignaturesRequest{
			ContractId: &abi.ContractId{
				Name: contract0.GetName(),
				Tag:  contract0.GetTag(),
			},
			MethodName: "transfer",
		})

		assert.NoError(t, err)
		assert.Equal(t, resp1.GetSignatures()[0], "transfer(address,uint256)")
	})
}

func (s *contractRegistryGRPCTestSuite) registerContract(contract *abi.Contract) (*registry.RegisterContractResponse, error) {
	return s.grpcClient.RegisterContract(context.Background(), &registry.RegisterContractRequest{
		Contract: contract,
	})
}

func (s *contractRegistryGRPCTestSuite) getContractRequest(name, tag string) *registry.GetContractRequest {
	return &registry.GetContractRequest{
		ContractId: &abi.ContractId{
			Name: name,
			Tag:  tag,
		},
	}
}
