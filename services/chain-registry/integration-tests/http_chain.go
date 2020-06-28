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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
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

	chainQuorumPrivTxType := utils.TesseraChainType
	chainQuorumPrivTxURL := "http://tessera1:9080"

	s.T().Run("should fetch env chain geth by name", func(t *testing.T) {
		resp, err := s.client.GetChainByName(ctx, chainNameGeth)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.URLs), "should be one URLs")
		if len(resp.URLs) == 1 {
			assert.Equal(t, chainUrlGeth, resp.URLs[0])
		}
	})

	s.T().Run("should fetch env chain besu by name", func(t *testing.T) {
		resp, err := s.client.GetChainByName(ctx, chainNameBesu)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.URLs), "should be one URLs")
		if len(resp.URLs) == 1 {
			assert.Equal(t, chainUrlBesu, resp.URLs[0])
		}
	})

	s.T().Run("should fetch env chain quorum by name", func(t *testing.T) {
		resp, err := s.client.GetChainByName(ctx, chainNameQuorum)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(resp.URLs), "should be one URLs")
		if len(resp.URLs) == 1 {
			assert.Equal(t, chainUrlQuorum, resp.URLs[0])
		}

		assert.Equal(t, 1, len(resp.PrivateTxManagers), "should be one PrivateTxManagers")
		if len(resp.PrivateTxManagers) == 1 {
			assert.Equal(t, chainQuorumPrivTxURL, resp.PrivateTxManagers[0].URL)
			assert.Equal(t, chainQuorumPrivTxType, resp.PrivateTxManagers[0].Type)
		}
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

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.UUID)
		chainUUID = resp.UUID
	})

	s.T().Run("should fetch registered chain by name", func(t *testing.T) {
		resp, err := s.client.GetChainByName(ctx, chainName)
		assert.NoError(t, err)
		assert.Equal(t, chainURL, resp.URLs[0])
	})

	s.T().Run("should update registered chain by UUID", func(t *testing.T) {
		err := s.client.UpdateChainByUUID(ctx, chainUUID, &models.Chain{
			ListenerCurrentBlock: &curBlockNumber,
		})

		assert.NoError(t, err)
	})

	s.T().Run("should fetch registered chain by UUID", func(t *testing.T) {
		resp, err := s.client.GetChainByUUID(ctx, chainUUID)
		assert.NoError(t, err)
		assert.Equal(t, chainURL, resp.URLs[0])
		assert.Equal(t, curBlockNumber, *resp.ListenerCurrentBlock)
	})

	s.T().Run("should delete registered chain by UUID", func(t *testing.T) {
		err := s.client.DeleteChainByUUID(ctx, chainUUID)
		assert.NoError(t, err)

		_, err = s.client.GetChainByUUID(ctx, chainUUID)
		assert.True(t, errors.IsNotFoundError(err))
	})
}

func (s *HttpChainTestSuite) TestChainRegistry_TesseraChainHappyFlow() {
	ctx := context.Background()
	chainName := fmt.Sprintf("TestTesseraChain%d", rand.Intn(1000))
	chainURL := "http://172.16.239.11:8545"
	privTxManagerURL := "http://172.16.239.11:8545"
	privTxManagerURLTwo := "http://172.16.239.11:9080"
	var chainUUID string

	s.T().Run("should register a new chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                  chainName,
			URLs:                  []string{chainURL},
			ListenerStartingBlock: &(&struct{ x uint64 }{1}).x,
			PrivateTxManagers: []*models.PrivateTxManagerModel{
				{
					URL:  privTxManagerURL,
					Type: utils.TesseraChainType,
				},
			},
		}
		resp, err := s.client.RegisterChain(ctx, &chain)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.UUID)
		chainUUID = resp.UUID
	})

	s.T().Run("should fetch registered chain by name", func(t *testing.T) {
		resp, err := s.client.GetChainByName(ctx, chainName)
		assert.NoError(t, err)
		assert.Equal(t, privTxManagerURL, resp.PrivateTxManagers[0].URL)
	})

	s.T().Run("should update registered chain by UUID", func(t *testing.T) {
		err := s.client.UpdateChainByUUID(ctx, chainUUID, &models.Chain{
			PrivateTxManagers: []*models.PrivateTxManagerModel{
				{
					URL:  privTxManagerURLTwo,
					Type: utils.TesseraChainType,
				},
			},
		})

		assert.NoError(t, err)
	})

	s.T().Run("should fetch registered chain by UUID", func(t *testing.T) {
		resp, err := s.client.GetChainByUUID(ctx, chainUUID)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.PrivateTxManagers))
		if len(resp.PrivateTxManagers) == 1 {
			assert.Equal(t, privTxManagerURLTwo, resp.PrivateTxManagers[0].URL)
		}
	})

	s.T().Run("should deleted registered chain by UUID", func(t *testing.T) {
		err := s.client.DeleteChainByUUID(ctx, chainUUID)
		assert.NoError(t, err)

		_, err = s.client.GetChainByUUID(ctx, chainUUID)
		assert.Error(t, err)
		assert.True(t, errors.IsNotFoundError(err), "should be DataErr, instead "+err.Error())
	})
}

func (s *HttpChainTestSuite) TestChainRegistry_ChainErrors() {
	ctx := context.Background()
	chainName := fmt.Sprintf("TestChainErr%d", rand.Intn(1000))
	var chainUUID string

	s.T().Run("should fail to register a new invalid chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                    chainName,
			URLs:                    []string{"http://test1.com"},
			ListenerBackOffDuration: &(&struct{ x string }{"1000"}).x,
		}
		_, err := s.client.RegisterChain(ctx, &chain)

		assert.Error(t, err)
		assert.True(t, errors.IsDataError(err))
	})

	s.T().Run("should fail to register a new invalid chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                    chainName,
			URLs:                    []string{"$%^^"},
			ListenerBackOffDuration: &(&struct{ x string }{"1s"}).x,
		}
		_, err := s.client.RegisterChain(ctx, &chain)

		assert.Error(t, err)
		assert.True(t, errors.IsDataError(err))
	})

	s.T().Run("should register a new chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                    chainName,
			URLs:                    []string{"http://test1.com"},
			ListenerBackOffDuration: &(&struct{ x string }{"1s"}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
		}
		resp, err := s.client.RegisterChain(ctx, &chain)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.UUID)
		chainUUID = resp.UUID
	})

	s.T().Run("should fail to update chain by UUID with invalid data", func(t *testing.T) {
		err := s.client.UpdateChainByUUID(ctx, chainUUID, &models.Chain{
			ListenerBackOffDuration: &(&struct{ x string }{"1000"}).x,
		})

		assert.Error(t, err)
		assert.True(t, errors.IsDataError(err))
	})

	s.T().Run("should fail to update chain by UUID with invalid data", func(t *testing.T) {
		err := s.client.UpdateChainByUUID(ctx, chainUUID, &models.Chain{
			URLs: []string{"$%^^"},
		})

		assert.Error(t, err)
		assert.True(t, errors.IsDataError(err))
	})

	s.T().Run("should deleted registered chain by UUID", func(t *testing.T) {
		err := s.client.DeleteChainByUUID(ctx, chainUUID)
		assert.NoError(t, err)

		_, err = s.client.GetChainByUUID(ctx, chainUUID)
		assert.Error(t, err)
	})
}

func (s *HttpChainTestSuite) TestChainRegistry_TesseraChainErrs() {
	ctx := context.Background()
	chainName := fmt.Sprintf("TestTesseraChainErr%d", rand.Intn(1000))
	var chainUUID string

	s.T().Run("should fail to register a new invalid tessera chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                  chainName,
			URLs:                  []string{"http://127.0.0.1:8545"},
			ListenerStartingBlock: &(&struct{ x uint64 }{1}).x,
			PrivateTxManagers: []*models.PrivateTxManagerModel{
				{
					URL:  "http://127.0.0.1:9080",
					Type: "InvalidType",
				},
			},
		}
		_, err := s.client.RegisterChain(ctx, &chain)

		assert.Error(t, err)
		assert.True(t, errors.IsDataError(err), "should be DataErr, instead "+err.Error())
	})

	s.T().Run("should fail to register a new invalid tessera chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                  chainName,
			URLs:                  []string{"http://127.0.0.1:8545"},
			ListenerStartingBlock: &(&struct{ x uint64 }{1}).x,
			PrivateTxManagers: []*models.PrivateTxManagerModel{
				{
					Type: utils.TesseraChainType,
				},
			},
		}
		_, err := s.client.RegisterChain(ctx, &chain)

		assert.Error(t, err)
		assert.True(t, errors.IsDataError(err), "should be DataErr, instead "+err.Error())
	})

	s.T().Run("should register a new tessera chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                  chainName,
			URLs:                  []string{"http://127.0.0.1:8545"},
			ListenerStartingBlock: &(&struct{ x uint64 }{1}).x,
			PrivateTxManagers: []*models.PrivateTxManagerModel{
				{
					URL:  "http://127.0.0.1:9080",
					Type: utils.TesseraChainType,
				},
			},
		}
		resp, err := s.client.RegisterChain(ctx, &chain)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.UUID)
		chainUUID = resp.UUID
	})

	s.T().Run("should fail to update tessera chain with invalid data", func(t *testing.T) {
		err := s.client.UpdateChainByUUID(ctx, chainUUID, &models.Chain{
			PrivateTxManagers: []*models.PrivateTxManagerModel{
				{
					URL: "http://127.0.0.1:9080",
				},
			},
		})

		assert.Error(t, err)
		assert.True(t, errors.IsDataError(err), "should be DataErr, instead "+err.Error())
	})

	s.T().Run("should deleted registered chain by UUID", func(t *testing.T) {
		err := s.client.DeleteChainByUUID(ctx, chainUUID)
		assert.NoError(t, err)

		_, err = s.client.GetChainByUUID(ctx, chainUUID)
		assert.Error(t, err)
	})
}
