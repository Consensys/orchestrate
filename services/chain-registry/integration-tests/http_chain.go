// +build integration

package integrationtests

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

// JobsTestSuite is a test suite for Transaction API jobs controller
type HttpChainTestSuite struct {
	suite.Suite
	baseURL string
	client  client.ChainClient
	env     *IntegrationEnvironment
}

func (s *HttpChainTestSuite) SetupSuite() {
	s.client = client.DialWithDefaultOptions(s.baseURL)
}

func (s *HttpChainTestSuite) TestChainRegistry_EnvChainImport() {
	ctx := context.Background()
	chainNameGeth := "geth"
	chainUrlGeth := "http://geth:8545"
	
	chainNameBesu := "besu"
	chainUrlBesu := "http://validator2:8545"
	
	chainNameQuorum := "quorum"
	chainUrlQuorum := "http://172.16.239.11:8545"
	
	s.T().Run("should fetch env chain geth by name", func(t *testing.T) {
		resp, err := s.client.GetChainByName(ctx, chainNameGeth)
		assert.Nil(t, err)
		assert.Equal(t, resp.URLs[0], chainUrlGeth)
	})
	
	s.T().Run("should fetch env chain besu by name", func(t *testing.T) {
		resp, err := s.client.GetChainByName(ctx, chainNameBesu)
		assert.Nil(t, err)
		assert.Equal(t, resp.URLs[0], chainUrlBesu)
	})
	
	s.T().Run("should fetch env chain quorum by name", func(t *testing.T) {
		resp, err := s.client.GetChainByName(ctx, chainNameQuorum)
		assert.Nil(t, err)
		assert.Equal(t, resp.URLs[0], chainUrlQuorum)
	})
}
 
func (s *HttpChainTestSuite) TestChainRegistry_ChainHappyFlow() {
	ctx := context.Background()
	chainName := fmt.Sprintf("TestChain%d", rand.Intn(1000))
	chainURL := "http://test1.com"
	var curBlockNumber uint64 = 666
	var chainUUID string

	s.T().Run("should register a new chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                    chainName,
			URLs:                    []string{chainURL},
			ListenerBackOffDuration: &(&struct{ x string }{"1s"}).x,
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
		}
		resp, err := s.client.RegisterChain(ctx, &chain)

		assert.Nil(t, err)
		assert.NotEmpty(t, resp.UUID)
		chainUUID = resp.UUID
	})

	s.T().Run("should fetch registered chain by name", func(t *testing.T) {
		resp, err := s.client.GetChainByName(ctx, chainName)
		assert.Nil(t, err)
		assert.Equal(t, resp.URLs[0], chainURL)
	})

	s.T().Run("should update registered chain by UUID", func(t *testing.T) {
		err := s.client.UpdateChainByUUID(ctx, chainUUID, &models.Chain{
			ListenerCurrentBlock: &curBlockNumber,
		})

		assert.Nil(t, err)
	})

	s.T().Run("should fetch registered chain by UUID", func(t *testing.T) {
		resp, err := s.client.GetChainByUUID(ctx, chainUUID)
		assert.Nil(t, err)
		assert.Equal(t, resp.URLs[0], chainURL)
		assert.Equal(t, *resp.ListenerCurrentBlock, curBlockNumber)
	})

	s.T().Run("should deleted registered chain by UUID", func(t *testing.T) {
		err := s.client.DeleteChainByUUID(ctx, chainUUID)
		assert.Nil(t, err)

		_, err = s.client.GetChainByUUID(ctx, chainUUID)
		assert.NotNil(t, err)
		assert.True(t, errors.IsNotFoundError(err), "should be DataErr, instead "+err.Error())
	})
}

func (s *HttpChainTestSuite) TestChainRegistry_ChainErrors() {
	ctx := context.Background()
	chainName := fmt.Sprintf("TestChain%d", rand.Intn(1000))
	var chainUUID string

	s.T().Run("should fail to register a new invalid chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                    chainName,
			URLs:                    []string{"http://test1.com"},
			ListenerBackOffDuration: &(&struct{ x string }{"1000"}).x,
		}
		_, err := s.client.RegisterChain(ctx, &chain)

		assert.NotNil(t, err)
		assert.True(t, errors.IsDataError(err), "should be DataErr, instead "+err.Error())
	})

	s.T().Run("should fail to register a new invalid chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                    chainName,
			URLs:                    []string{"$%^^"},
			ListenerBackOffDuration: &(&struct{ x string }{"1s"}).x,
		}
		_, err := s.client.RegisterChain(ctx, &chain)

		assert.NotNil(t, err)
		assert.True(t, errors.IsDataError(err), "should be DataErr, instead "+err.Error())
	})

	s.T().Run("should register a new chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                    chainName,
			URLs:                    []string{"http://test1.com"},
			ListenerBackOffDuration: &(&struct{ x string }{"1s"}).x,
		}
		resp, err := s.client.RegisterChain(ctx, &chain)

		assert.Nil(t, err)
		assert.NotEmpty(t, resp.UUID)
		chainUUID = resp.UUID
	})

	s.T().Run("should fail to update chain by UUID with invalid data", func(t *testing.T) {
		err := s.client.UpdateChainByUUID(ctx, chainUUID, &models.Chain{
			ListenerBackOffDuration: &(&struct{ x string }{"1000"}).x,
		})

		assert.NotNil(t, err)
		assert.True(t, errors.IsDataError(err), "should be DataErr, instead "+err.Error())
	})

	s.T().Run("should fail to update chain by UUID with invalid data", func(t *testing.T) {
		err := s.client.UpdateChainByUUID(ctx, chainUUID, &models.Chain{
			URLs: []string{"$%^^"},
		})

		assert.NotNil(t, err)
		assert.True(t, errors.IsDataError(err), "should be DataErr, instead "+err.Error())
	})

	s.T().Run("should deleted registered chain by UUID", func(t *testing.T) {
		err := s.client.DeleteChainByUUID(ctx, chainUUID)
		assert.Nil(t, err)

		_, err = s.client.GetChainByUUID(ctx, chainUUID)
		assert.NotNil(t, err)
	})
}
