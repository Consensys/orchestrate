// +build integration

package integrationtests

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

// txSentryTestSuite is a test suite for Transaction Sentry
type txSentryTestSuite struct {
	suite.Suite
	client client.TransactionSchedulerClient
	env    *IntegrationEnvironment
}

func (s *txSentryTestSuite) TestSentry() {
	chain := testutils.FakeChain()
	chainModel := &models.Chain{
		Name:     chain.Name,
		UUID:     chain.UUID,
		TenantID: chain.TenantID,
	}

	s.T().Run("test", func(t *testing.T) {
		defer gock.Off()

		gock.New(ChainRegistryURL).Get("/chains").Reply(200).JSON([]*models.Chain{chainModel})
		gock.New(ChainRegistryURL).Get("/chains/" + chain.UUID).Reply(200).JSON(chainModel)

		// TODO: Implement tests
	})
}
