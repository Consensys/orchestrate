// +build integration

package integrationtests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/consensys/orchestrate/pkg/multitenancy"
	"github.com/consensys/orchestrate/pkg/sdk/client"
	"testing"
)

// schedulesTestSuite is a test suite for Schedules controller
type schedulesTestSuite struct {
	suite.Suite
	client client.OrchestrateClient
	env    *IntegrationEnvironment
}

func (s *schedulesTestSuite) TestSuccess() {
	ctx := s.env.ctx

	s.T().Run("should create a schedule successfully", func(t *testing.T) {
		schedule, err := s.client.CreateSchedule(ctx, nil)

		require.NoError(t, err)
		assert.NotEmpty(t, schedule.UUID)
		assert.NotEmpty(t, schedule.CreatedAt)
		assert.Equal(t, multitenancy.DefaultTenant, schedule.TenantID)
	})
}

