// +build integration

package integrationtests

import (
	"testing"
	"time"

	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/testutils"
)

type chainsTestSuite struct {
	suite.Suite
	client client.OrchestrateClient
	env    *IntegrationEnvironment
}

func (s *chainsTestSuite) TestRegister() {
	ctx := s.env.ctx

	s.T().Run("should register chain successfully from latest", func(t *testing.T) {
		req := testutils.FakeRegisterChainRequest()
		req.Listener.FromBlock = "latest"
		req.URLs = []string{s.env.blockchainNodeURL}

		resp, err := s.client.RegisterChain(ctx, req)
		require.NoError(t, err)

		assert.Equal(t, req.Name, resp.Name)
		assert.Equal(t, req.URLs, resp.URLs)
		assert.Equal(t, multitenancy.DefaultTenant, resp.TenantID)
		assert.Equal(t, req.Listener.ExternalTxEnabled, resp.ListenerExternalTxEnabled)
		assert.Equal(t, "5s", resp.ListenerBackOffDuration)
		assert.Equal(t, req.Listener.Depth, resp.ListenerDepth)
		assert.Equal(t, req.Labels, resp.Labels)
		assert.NotEmpty(t, resp.UUID)
		assert.Greater(t, resp.ListenerStartingBlock, uint64(0))
		assert.NotEmpty(t, resp.CreatedAt)
		assert.NotEmpty(t, resp.UpdatedAt)
		assert.Equal(t, resp.CreatedAt, resp.UpdatedAt)

		err = s.client.DeleteChain(ctx, resp.UUID)
		assert.NoError(t, err)
	})

	s.T().Run("should register chain successfully from latest if fromBlock is empty", func(t *testing.T) {
		req := testutils.FakeRegisterChainRequest()
		req.Listener.FromBlock = ""
		req.URLs = []string{s.env.blockchainNodeURL}

		resp, err := s.client.RegisterChain(ctx, req)
		require.NoError(t, err)

		assert.Equal(t, req.Name, resp.Name)
		assert.Equal(t, req.URLs, resp.URLs)
		assert.Equal(t, multitenancy.DefaultTenant, resp.TenantID)
		assert.Equal(t, req.Listener.ExternalTxEnabled, resp.ListenerExternalTxEnabled)
		assert.Equal(t, "5s", resp.ListenerBackOffDuration)
		assert.Equal(t, req.Listener.Depth, resp.ListenerDepth)
		assert.NotEmpty(t, resp.UUID)
		assert.Greater(t, resp.ListenerStartingBlock, uint64(0))
		assert.NotEmpty(t, resp.CreatedAt)
		assert.NotEmpty(t, resp.UpdatedAt)
		assert.Equal(t, resp.CreatedAt, resp.UpdatedAt)

		err = s.client.DeleteChain(ctx, resp.UUID)
		assert.NoError(t, err)
	})

	s.T().Run("should register chain successfully from 0", func(t *testing.T) {
		req := testutils.FakeRegisterChainRequest()
		req.Listener.FromBlock = "0"
		req.URLs = []string{s.env.blockchainNodeURL}

		resp, err := s.client.RegisterChain(ctx, req)
		require.NoError(t, err)

		assert.Equal(t, req.Name, resp.Name)
		assert.Equal(t, req.URLs, resp.URLs)
		assert.Equal(t, multitenancy.DefaultTenant, resp.TenantID)
		assert.Equal(t, req.Listener.ExternalTxEnabled, resp.ListenerExternalTxEnabled)
		assert.Equal(t, "5s", resp.ListenerBackOffDuration)
		assert.Equal(t, req.Listener.Depth, resp.ListenerDepth)
		assert.NotEmpty(t, resp.UUID)
		assert.Equal(t, uint64(0), resp.ListenerStartingBlock)
		assert.NotEmpty(t, resp.CreatedAt)
		assert.NotEmpty(t, resp.UpdatedAt)
		assert.Equal(t, resp.CreatedAt, resp.UpdatedAt)

		err = s.client.DeleteChain(ctx, resp.UUID)
		assert.NoError(t, err)
	})

	s.T().Run("should register chain successfully from 666", func(t *testing.T) {
		req := testutils.FakeRegisterChainRequest()
		req.Listener.FromBlock = "666"
		req.URLs = []string{s.env.blockchainNodeURL}

		resp, err := s.client.RegisterChain(ctx, req)
		require.NoError(t, err)

		assert.Equal(t, req.Name, resp.Name)
		assert.Equal(t, req.URLs, resp.URLs)
		assert.Equal(t, multitenancy.DefaultTenant, resp.TenantID)
		assert.Equal(t, req.Listener.ExternalTxEnabled, resp.ListenerExternalTxEnabled)
		assert.Equal(t, "5s", resp.ListenerBackOffDuration)
		assert.Equal(t, req.Listener.Depth, resp.ListenerDepth)
		assert.NotEmpty(t, resp.UUID)
		assert.Equal(t, uint64(666), resp.ListenerStartingBlock)
		assert.NotEmpty(t, resp.CreatedAt)
		assert.NotEmpty(t, resp.UpdatedAt)
		assert.Equal(t, resp.CreatedAt, resp.UpdatedAt)

		err = s.client.DeleteChain(ctx, resp.UUID)
		assert.NoError(t, err)
	})

	s.T().Run("should fail with 400 if payload is invalid", func(t *testing.T) {
		req := testutils.FakeRegisterChainRequest()
		req.URLs = nil

		_, err := s.client.RegisterChain(ctx, req)
		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail with 400 if invalid backoff duration", func(t *testing.T) {
		req := testutils.FakeRegisterChainRequest()
		req.Listener.BackOffDuration = "invalidDuration"

		_, err := s.client.RegisterChain(ctx, req)
		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail with 400 if invalid urls", func(t *testing.T) {
		req := testutils.FakeRegisterChainRequest()
		req.URLs = []string{"invalidURL"}

		_, err := s.client.RegisterChain(ctx, req)
		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail with 400 if invalid private tx manager type", func(t *testing.T) {
		req := testutils.FakeRegisterChainRequest()
		req.PrivateTxManager.Type = "invalidType"

		_, err := s.client.RegisterChain(ctx, req)
		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail with 400 if invalid private tx manager url", func(t *testing.T) {
		req := testutils.FakeRegisterChainRequest()
		req.PrivateTxManager.URL = "invalidURL"

		_, err := s.client.RegisterChain(ctx, req)
		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail with 422 if URL is not reachable", func(t *testing.T) {
		_, err := s.client.RegisterChain(ctx, testutils.FakeRegisterChainRequest())
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	s.T().Run("should fail with 409 if chain with same name and tenant already exists", func(t *testing.T) {
		req := testutils.FakeRegisterChainRequest()
		req.URLs = []string{s.env.blockchainNodeURL}

		resp, err := s.client.RegisterChain(ctx, req)
		require.NoError(t, err)

		_, err = s.client.RegisterChain(ctx, req)
		assert.True(t, errors.IsAlreadyExistsError(err))

		err = s.client.DeleteChain(ctx, resp.UUID)
		assert.NoError(t, err)
	})

	s.T().Run("should fail with 500 to register chain if postgres is down", func(t *testing.T) {
		req := testutils.FakeRegisterChainRequest()

		err := s.env.client.Stop(ctx, postgresContainerID)
		assert.NoError(t, err)

		_, err = s.client.RegisterChain(ctx, req)
		assert.Error(t, err)

		err = s.env.client.StartServiceAndWait(ctx, postgresContainerID, 10*time.Second)
		assert.NoError(t, err)
	})
}

func (s *chainsTestSuite) TestSearch() {
	ctx := s.env.ctx
	req := testutils.FakeRegisterChainRequest()
	req.URLs = []string{s.env.blockchainNodeURL}
	chain, err := s.client.RegisterChain(ctx, req)
	require.NoError(s.T(), err)

	s.T().Run("should search chain by name successfully", func(t *testing.T) {
		resp, err := s.client.SearchChains(ctx, &entities.ChainFilters{
			Names: []string{chain.Name},
		})
		require.NoError(t, err)

		assert.Len(t, resp, 1)
		assert.Equal(t, chain.UUID, resp[0].UUID)
	})

	s.T().Run("should return empty array if nothing is found", func(t *testing.T) {
		resp, err := s.client.SearchChains(ctx, &entities.ChainFilters{
			Names: []string{"inexistentName"},
		})
		require.NoError(t, err)
		assert.Empty(t, resp)
	})

	err = s.client.DeleteChain(ctx, chain.UUID)
	assert.NoError(s.T(), err)
}

func (s *chainsTestSuite) TestGetOne() {
	ctx := s.env.ctx

	s.T().Run("should get chain successfully", func(t *testing.T) {
		req := testutils.FakeRegisterChainRequest()
		req.URLs = []string{s.env.blockchainNodeURL}
		chain, err := s.client.RegisterChain(ctx, req)
		require.NoError(s.T(), err)

		resp, err := s.client.GetChain(ctx, chain.UUID)
		require.NoError(t, err)
		assert.Equal(t, chain.UUID, resp.UUID)

		err = s.client.DeleteChain(ctx, chain.UUID)
		require.NoError(s.T(), err)
	})
}

func (s *chainsTestSuite) TestUpdate() {
	ctx := s.env.ctx
	req := testutils.FakeRegisterChainRequest()
	req.URLs = []string{s.env.blockchainNodeURL}
	chain, err := s.client.RegisterChain(ctx, req)
	require.NoError(s.T(), err)

	s.T().Run("should update chain name successfully", func(t *testing.T) {
		req := testutils.FakeUpdateChainRequest()
		req.Name = "newName"

		resp, err := s.client.UpdateChain(ctx, chain.UUID, req)
		require.NoError(t, err)

		assert.Equal(t, req.Name, resp.Name)
		assert.Equal(t, req.Labels, resp.Labels)
		assert.NotEqual(t, resp.CreatedAt, resp.UpdatedAt)
	})

	s.T().Run("should update chain listener successfully", func(t *testing.T) {
		req := testutils.FakeUpdateChainRequest()
		req.Listener.CurrentBlock = 666
		req.Listener.Depth = 2
		req.Listener.ExternalTxEnabled = utils.ToPtr(true).(*bool)
		req.Listener.BackOffDuration = "10s"

		resp, err := s.client.UpdateChain(ctx, chain.UUID, req)
		require.NoError(t, err)

		assert.Equal(t, req.Listener.CurrentBlock, resp.ListenerCurrentBlock)
		assert.Equal(t, req.Listener.Depth, resp.ListenerDepth)
		assert.Equal(t, req.Listener.BackOffDuration, resp.ListenerBackOffDuration)
		assert.Equal(t, req.Listener.ExternalTxEnabled, resp.ListenerExternalTxEnabled)
		assert.NotEqual(t, resp.CreatedAt, resp.UpdatedAt)
	})

	s.T().Run("should update chain private tx manager successfully", func(t *testing.T) {
		req := testutils.FakeUpdateChainRequest()
		req.PrivateTxManager = &api.PrivateTxManagerRequest{
			URL:  "http://myURLUpdated:8545",
			Type: entities.TesseraChainType,
		}

		resp, err := s.client.UpdateChain(ctx, chain.UUID, req)
		require.NoError(t, err)

		assert.Equal(t, req.PrivateTxManager.URL, resp.PrivateTxManager.URL)
		assert.Equal(t, req.PrivateTxManager.Type, resp.PrivateTxManager.Type)
		assert.NotEqual(t, resp.CreatedAt, resp.UpdatedAt)
	})

	err = s.client.DeleteChain(ctx, chain.UUID)
	assert.NoError(s.T(), err)
}
