// +build integration

package integrationtests

import (
	"github.com/stretchr/testify/require"
	"github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/types/testutils"
)

type faucetsTestSuite struct {
	suite.Suite
	client client.OrchestrateClient
	env    *IntegrationEnvironment
}

func (s *faucetsTestSuite) TestRegister() {
	ctx := s.env.ctx

	s.T().Run("should register faucet successfully", func(t *testing.T) {
		req := testutils.FakeRegisterFaucetRequest()

		resp, err := s.client.RegisterFaucet(ctx, req)
		require.NoError(t, err)

		assert.Equal(t, req.CreditorAccount, resp.CreditorAccount)
		assert.Equal(t, req.ChainRule, resp.ChainRule)
		assert.Equal(t, req.MaxBalance, resp.MaxBalance)
		assert.Equal(t, req.Amount, resp.Amount)
		assert.Equal(t, req.Name, resp.Name)
		assert.Equal(t, req.Cooldown, resp.Cooldown)
		assert.NotEmpty(t, resp.UUID)
		assert.NotEmpty(t, resp.CreatedAt)
		assert.NotEmpty(t, resp.UpdatedAt)
		assert.Equal(t, resp.CreatedAt, resp.UpdatedAt)

		err = s.client.DeleteFaucet(ctx, resp.UUID)
		assert.NoError(t, err)
	})

	s.T().Run("should fail to register faucet with same name and tenant", func(t *testing.T) {
		req := testutils.FakeRegisterFaucetRequest()

		resp, err := s.client.RegisterFaucet(ctx, req)
		require.NoError(t, err)

		_, err = s.client.RegisterFaucet(ctx, req)
		assert.True(t, errors.IsAlreadyExistsError(err))

		err = s.client.DeleteFaucet(ctx, resp.UUID)
		assert.NoError(t, err)
	})

	s.T().Run("should fail to register faucet if postgres is down", func(t *testing.T) {
		req := testutils.FakeRegisterFaucetRequest()

		err := s.env.client.Stop(ctx, postgresContainerID)
		assert.NoError(t, err)

		_, err = s.client.RegisterFaucet(ctx, req)
		assert.Error(t, err)

		err = s.env.client.StartServiceAndWait(ctx, postgresContainerID, 10*time.Second)
		assert.NoError(t, err)
	})
}

func (s *faucetsTestSuite) TestSearch() {
	ctx := s.env.ctx
	req := testutils.FakeRegisterFaucetRequest()
	faucet, err := s.client.RegisterFaucet(ctx, req)
	require.NoError(s.T(), err)

	s.T().Run("should search faucet by name successfully", func(t *testing.T) {
		resp, err := s.client.SearchFaucets(ctx, &entities.FaucetFilters{
			Names: []string{faucet.Name},
		})
		require.NoError(t, err)

		assert.Len(t, resp, 1)
		assert.Equal(t, faucet.UUID, resp[0].UUID)
	})

	s.T().Run("should search faucet by chain_rule successfully", func(t *testing.T) {
		resp, err := s.client.SearchFaucets(ctx, &entities.FaucetFilters{
			ChainRule: faucet.ChainRule,
		})
		require.NoError(t, err)

		assert.Len(t, resp, 1)
		assert.Equal(t, faucet.UUID, resp[0].UUID)
	})

	err = s.client.DeleteFaucet(ctx, faucet.UUID)
	require.NoError(s.T(), err)
}

func (s *faucetsTestSuite) TestGetOne() {
	ctx := s.env.ctx
	req := testutils.FakeRegisterFaucetRequest()
	faucet, err := s.client.RegisterFaucet(ctx, req)
	require.NoError(s.T(), err)

	s.T().Run("should get faucet successfully", func(t *testing.T) {
		resp, err := s.client.GetFaucet(ctx, faucet.UUID)
		require.NoError(t, err)
		assert.Equal(t, faucet.UUID, resp.UUID)
	})

	err = s.client.DeleteFaucet(ctx, faucet.UUID)
	require.NoError(s.T(), err)
}

func (s *faucetsTestSuite) TestUpdate() {
	ctx := s.env.ctx
	req := testutils.FakeRegisterFaucetRequest()
	faucet, err := s.client.RegisterFaucet(ctx, req)
	require.NoError(s.T(), err)

	s.T().Run("should update faucet successfully", func(t *testing.T) {
		req := testutils.FakeUpdateFaucetRequest()

		resp, err := s.client.UpdateFaucet(ctx, faucet.UUID, req)
		require.NoError(t, err)

		assert.Equal(t, req.CreditorAccount, resp.CreditorAccount)
		assert.Equal(t, req.ChainRule, resp.ChainRule)
		assert.Equal(t, req.MaxBalance, resp.MaxBalance)
		assert.Equal(t, req.Amount, resp.Amount)
		assert.Equal(t, req.Name, resp.Name)
		assert.Equal(t, req.Cooldown, resp.Cooldown)
		assert.NotEmpty(t, resp.UUID)
		assert.NotEmpty(t, resp.CreatedAt)
		assert.NotEmpty(t, resp.UpdatedAt)
		assert.NotEqual(t, resp.CreatedAt, resp.UpdatedAt)

		err = s.client.DeleteFaucet(ctx, resp.UUID)
		assert.NoError(t, err)
	})
}
