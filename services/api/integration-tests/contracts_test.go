// +build integration

package integrationtests

import (
	"context"
	"testing"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

type contractsTestSuite struct {
	suite.Suite
	client client.OrchestrateClient
	env    *IntegrationEnvironment
}

func (s *contractsTestSuite) TestContractRegistry_Register() {
	ctx := context.Background()

	s.T().Run("should register a contract with tag", func(t *testing.T) {
		txRequest := testutils.FakeRegisterContractRequest()

		resp, err := s.client.RegisterContract(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, txRequest.Name, resp.ID.Name)
		assert.Equal(t, txRequest.Tag, resp.ID.Tag)
		assert.Equal(t, txRequest.DeployedBytecode, resp.DeployedBytecode)
		assert.Equal(t, txRequest.Bytecode, resp.Bytecode)
		assert.NotEmpty(t, resp.Constructor.Signature)
		assert.NotEmpty(t, resp.Events)
		assert.NotEmpty(t, resp.Methods)

		abi, err := json.Marshal(txRequest.ABI)
		assert.NoError(t, err)
		assert.Equal(t, string(abi), resp.ABI)
	})

	s.T().Run("should register a contract with tag latest", func(t *testing.T) {
		txRequest := testutils.FakeRegisterContractRequest()
		txRequest.Tag = ""

		resp, err := s.client.RegisterContract(ctx, txRequest)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, txRequest.Name, resp.ID.Name)
		assert.Equal(t, entities.DefaultTagValue, resp.ID.Tag)
		assert.Equal(t, txRequest.DeployedBytecode, resp.DeployedBytecode)
		assert.Equal(t, txRequest.Bytecode, resp.Bytecode)
	})

	s.T().Run("should fail with invalidFormatError if payload is invalid", func(t *testing.T) {
		txRequest := testutils.FakeRegisterContractRequest()
		txRequest.Name = ""

		_, err := s.client.RegisterContract(ctx, txRequest)
		assert.Error(t, err)
		assert.True(t, errors.IsInvalidFormatError(err), err.Error())
	})

	s.T().Run("should fail with encodingError if ABI payload is invalid", func(t *testing.T) {
		txRequest := testutils.FakeRegisterContractRequest()
		txRequest.ABI = "{asd}asdasd"

		_, err := s.client.RegisterContract(ctx, txRequest)
		assert.Error(t, err)
		assert.True(t, errors.IsEncodingError(err), err.Error())
	})
}

func (s *contractsTestSuite) TestContractRegistry_Get() {
	contractName := "contract_" + utils.RandomString(5)
	ctx := context.Background()
	txRequest := testutils.FakeRegisterContractRequest()
	txRequest.Name = contractName
	_, err := s.client.RegisterContract(ctx, txRequest)
	if err != nil {
		assert.Fail(s.T(), err.Error())
		return
	}

	s.T().Run("should get all contracts", func(t *testing.T) {
		resp, err := s.client.GetContractsCatalog(ctx)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Contains(t, resp, contractName)
	})

	s.T().Run("should get all tags of a contract", func(t *testing.T) {
		resp, err := s.client.GetContractTags(ctx, contractName)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Len(t, resp, 1)
		assert.Contains(t, resp, txRequest.Tag)
	})

	s.T().Run("should get a contract", func(t *testing.T) {
		resp, err := s.client.GetContract(ctx, txRequest.Name, txRequest.Tag)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, txRequest.Name, resp.ID.Name)
		abi, err := json.Marshal(txRequest.ABI)
		assert.NoError(t, err)
		assert.Equal(t, string(abi), resp.ABI)
	})

	s.T().Run("should get a contract method signatures", func(t *testing.T) {
		resp, err := s.client.GetContractMethodSignatures(ctx, txRequest.Name, txRequest.Tag, "")
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
	
		assert.Contains(t, resp, "transferFrom(address,address,uint256)")
		assert.Contains(t, resp, "totalSupply()")
		assert.Contains(t, resp, "approve(address,uint256)")
		
		resp2, err := s.client.GetContractMethodSignatures(ctx, txRequest.Name, txRequest.Tag, "balanceOf")
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		
		assert.Len(t, resp2, 1)
		assert.Contains(t, resp2, "balanceOf(address)")
	})
}

func (s *contractsTestSuite) TestContractRegistry_CodeHash() {
	ctx := context.Background()
	contractName := "contract_" + utils.RandomString(5)
	txRequest := testutils.FakeRegisterContractRequest()
	txRequest.Name = contractName
	_, err := s.client.RegisterContract(ctx, txRequest)
	if err != nil {
		assert.Fail(s.T(), err.Error())
		return
	}
	
	address := ethcommon.HexToAddress(utils.RandHexString(10))
	address2 := ethcommon.HexToAddress(utils.RandHexString(10))
	codeHash := ethcommon.HexToHash(utils.RandHexString(20))
	codeHash2 := "0xd63259750ca3b56efab25f0646a4d1fb659b6b643474506e1be24d81f9e55fd8"
	chainID := "2017"
	
	s.T().Run("should set contract code hashes successfully", func(t *testing.T) {
		err := s.client.SetContractAddressCodeHash(ctx, &api.SetContractCodeHashRequest{
			Address:  address.String(),
			CodeHash: codeHash.String(),
			ChainID:  chainID,
		})
	
		assert.NoError(t, err)
		
		err = s.client.SetContractAddressCodeHash(ctx, &api.SetContractCodeHashRequest{
			Address:  address2.String(),
			CodeHash: codeHash2,
			ChainID:  chainID,
		})
	
		assert.NoError(t, err)
	})
	
	s.T().Run("should get default contract event by sigHash successfully", func(t *testing.T) {
		resp, err := s.client.GetContractEventsBySigHash(ctx, address.String(), &api.GetContractEventsBySignHashRequest{
			ChainID:           chainID,
			SigHash:           "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			IndexedInputCount: 2,
		})
	
		assert.NoError(t, err)
		if len(resp.DefaultEvents) == 0 {
			assert.Fail(t, "expected some default events")
		}
	
		event := &ethAbi.Event{}
		err = json.Unmarshal([]byte(resp.DefaultEvents[0]), event)
		assert.NoError(t, err)
		assert.Equal(t, "Transfer", event.Name)
	})
	
	s.T().Run("should get contract event by sigHash successfully", func(t *testing.T) {
		resp, err := s.client.GetContractEventsBySigHash(ctx, address2.String(), &api.GetContractEventsBySignHashRequest{
			ChainID:           chainID,
			SigHash:           "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			IndexedInputCount: 2,
		})
	
		assert.NoError(t, err)
		if resp.Event == "" {
			assert.Fail(t, "expected contract events")
		}
	
		event := &ethAbi.Event{}
		err = json.Unmarshal([]byte(resp.Event), event)
		assert.NoError(t, err)
		assert.Equal(t, "Transfer", event.Name)
	})
	
	s.T().Run("should fail to set contract code hashes if invalid address", func(t *testing.T) {
		err := s.client.SetContractAddressCodeHash(ctx, &api.SetContractCodeHashRequest{
			Address:  "InvalidAddress",
			CodeHash: codeHash.String(),
			ChainID:  chainID,
		})

		assert.Error(t, err)
		assert.True(t, errors.IsInvalidFormatError(err), "IsInvalidFormatError")
	})
	
	s.T().Run("should fail to set contract code hashes if invalid codeHash", func(t *testing.T) {
		err := s.client.SetContractAddressCodeHash(ctx, &api.SetContractCodeHashRequest{
			Address:  address.String(),
			CodeHash: "{invalidCodeHash}",
			ChainID:  chainID,
		})
	
		assert.Error(t, err)
		assert.True(t, errors.IsInvalidFormatError(err), "IsInvalidFormatError")
	})
}
