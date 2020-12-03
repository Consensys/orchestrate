// +build integration

package integrationtests

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
	"math/rand"
	http2 "net/http"
	"testing"
	"time"
)

// chainsTestSuite is a test suite for Chains API
type chainsTestSuite struct {
	suite.Suite
	baseURL string
	client  client.ChainClient
	env     *IntegrationEnvironment
}

func (s *chainsTestSuite) SetupSuite() {
	s.client = client.DialWithDefaultOptions(s.baseURL)
}

func (s *chainsTestSuite) TestChainRegistry_ChainHappyFlow() {
	ctx := context.Background()
	chainName := fmt.Sprintf("TestChain%d", rand.Intn(1000))
	var curBlockNumber uint64 = 666
	var chainUUID string

	s.T().Run("should fetch imported chain by name", func(t *testing.T) {
		resp, err := s.client.GetChainByName(ctx, "ganache")
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.NotEmpty(t, resp.UUID)
	})

	s.T().Run("should register a new chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                    chainName,
			URLs:                    []string{s.env.blockchainNodeURL},
			ListenerBackOffDuration: &(&struct{ x string }{"1s"}).x,
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
		}
		resp, err := s.client.RegisterChain(ctx, &chain)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.NotEmpty(t, resp.UUID)
		chainUUID = resp.UUID
	})

	s.T().Run("should fetch registered chain by name", func(t *testing.T) {
		resp, err := s.client.GetChainByName(ctx, chainName)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, s.env.blockchainNodeURL, resp.URLs[0])
	})

	s.T().Run("should update registered chain by UUID", func(t *testing.T) {
		err := s.client.UpdateChainByUUID(ctx, chainUUID, &models.Chain{
			ListenerCurrentBlock: &curBlockNumber,
		})

		assert.NoError(t, err)
	})

	s.T().Run("should fetch registered chain by UUID", func(t *testing.T) {
		resp, err := s.client.GetChainByUUID(ctx, chainUUID)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, s.env.blockchainNodeURL, resp.URLs[0])
		assert.Equal(t, curBlockNumber, *resp.ListenerCurrentBlock)
	})

	s.T().Run("should delete registered chain by UUID", func(t *testing.T) {
		err := s.client.DeleteChainByUUID(ctx, chainUUID)
		assert.NoError(t, err)

		_, err = s.client.GetChainByUUID(ctx, chainUUID)
		assert.True(t, errors.IsNotFoundError(err))
	})
}

func (s *chainsTestSuite) TestChainRegistry_TesseraChainHappyFlow() {
	ctx := context.Background()
	chainName := fmt.Sprintf("TestTesseraChain%d", rand.Intn(1000))
	privTxManagerURL := "http://172.16.239.11:8545"
	privTxManagerURLTwo := "http://172.16.239.11:9080"
	var chainUUID string

	s.T().Run("should register a new chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                  chainName,
			URLs:                  []string{s.env.blockchainNodeURL},
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
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.NotEmpty(t, resp.UUID)
		chainUUID = resp.UUID
	})

	s.T().Run("should fetch registered chain by name", func(t *testing.T) {
		resp, err := s.client.GetChainByName(ctx, chainName)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

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
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

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

func (s *chainsTestSuite) TestChainRegistry_ChainErrors() {
	ctx := context.Background()
	chainName := fmt.Sprintf("TestChainErr%d", rand.Intn(1000))
	var chainUUID string

	s.T().Run("should fail to register a new invalid chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                    chainName,
			URLs:                    []string{"http://invalid:8545"},
			ListenerBackOffDuration: &(&struct{ x string }{"1000"}).x,
		}
		_, err := s.client.RegisterChain(ctx, &chain)

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail to register a new invalid chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                    chainName,
			URLs:                    []string{"$%^^"},
			ListenerBackOffDuration: &(&struct{ x string }{"1s"}).x,
		}
		_, err := s.client.RegisterChain(ctx, &chain)

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail to update chain by UUID with invalid data", func(t *testing.T) {
		chain := models.Chain{
			Name:                    chainName,
			URLs:                    []string{s.env.blockchainNodeURL},
			ListenerBackOffDuration: &(&struct{ x string }{"1s"}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
		}
		resp, err := s.client.RegisterChain(ctx, &chain)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.NotEmpty(t, resp.UUID)
		chainUUID = resp.UUID

		err = s.client.UpdateChainByUUID(ctx, chainUUID, &models.Chain{
			ListenerBackOffDuration: &(&struct{ x string }{"1000"}).x,
		})
		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail to update chain by UUID with invalid data", func(t *testing.T) {
		err := s.client.UpdateChainByUUID(ctx, chainUUID, &models.Chain{
			URLs: []string{"$%^^"},
		})

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should delete registered chain by UUID", func(t *testing.T) {
		err := s.client.DeleteChainByUUID(ctx, chainUUID)
		assert.NoError(t, err)

		_, err = s.client.GetChainByUUID(ctx, chainUUID)
		assert.Error(t, err)
	})
}

func (s *chainsTestSuite) TestChainRegistry_TesseraChainErrs() {
	ctx := context.Background()
	chainName := fmt.Sprintf("TestTesseraChainErr%d", rand.Intn(1000))
	var chainUUID string

	s.T().Run("should fail to register a new invalid tessera chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                  chainName,
			URLs:                  []string{s.env.blockchainNodeURL},
			ListenerStartingBlock: &(&struct{ x uint64 }{1}).x,
			PrivateTxManagers: []*models.PrivateTxManagerModel{
				{
					URL:  "http://127.0.0.1:9080",
					Type: "InvalidType",
				},
			},
		}
		_, err := s.client.RegisterChain(ctx, &chain)

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail to register a new invalid tessera chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                  chainName,
			URLs:                  []string{s.env.blockchainNodeURL},
			ListenerStartingBlock: &(&struct{ x uint64 }{1}).x,
			PrivateTxManagers:     []*models.PrivateTxManagerModel{{Type: utils.TesseraChainType}},
		}
		_, err := s.client.RegisterChain(ctx, &chain)

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should register a new tessera chain", func(t *testing.T) {
		chain := models.Chain{
			Name:                  chainName,
			URLs:                  []string{s.env.blockchainNodeURL},
			ListenerStartingBlock: &(&struct{ x uint64 }{1}).x,
			PrivateTxManagers: []*models.PrivateTxManagerModel{
				{
					URL:  "http://127.0.0.1:9080",
					Type: utils.TesseraChainType,
				},
			},
		}
		resp, err := s.client.RegisterChain(ctx, &chain)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

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

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should deleted registered chain by UUID", func(t *testing.T) {
		err := s.client.DeleteChainByUUID(ctx, chainUUID)
		assert.NoError(t, err)

		_, err = s.client.GetChainByUUID(ctx, chainUUID)
		assert.Error(t, err)
	})
}

func (s *chainsTestSuite) TestChainRegistry_ZHealthCheck() {
	type healthRes struct {
		Database string `json:"Database,omitempty"`
	}

	httpClient := http.NewClient(http.NewDefaultConfig())
	ctx := context.Background()
	s.T().Run("should retrieve positive health check over service dependencies", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		resp, err := httpClient.Do(req)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		assert.Equal(s.T(), 200, resp.StatusCode)
		status := healthRes{}
		err = json.UnmarshalBody(resp.Body, &status)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), "OK", status.Database)
	})

	s.T().Run("should retrieve a negative health check over postgres service", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		// Kill Kafka on first call so data is added in DB and status is CREATED but does not get updated to STARTED
		err = s.env.client.Stop(ctx, postgresContainerID)
		assert.NoError(t, err)

		resp, err := httpClient.Do(req)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		err = s.env.client.StartServiceAndWait(ctx, postgresContainerID, 10*time.Second)
		assert.NoError(s.T(), err)

		assert.Equal(s.T(), 503, resp.StatusCode)
		status := healthRes{}
		err = json.UnmarshalBody(resp.Body, &status)
		assert.NoError(s.T(), err)
		assert.NotEqual(s.T(), "OK", status.Database)
	})
}
