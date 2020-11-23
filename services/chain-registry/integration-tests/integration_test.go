// +build integration

package integrationtests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	integrationtest "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/integration-test"
)

type transactionSchedulerTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
	err error
}

func (s *transactionSchedulerTestSuite) SetupSuite() {
	err := integrationtest.StartEnvironment(context.Background(), s.env)
	if err != nil {
		s.env.logger.WithError(err).Error()
		if s.err == nil {
			s.err = err
		}
		return
	}

	s.env.logger.Infof("setup test suite has completed")
}

func (s *transactionSchedulerTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())

	if s.err != nil {
		s.Fail(s.err.Error())
	}
}

func TestChainRegistry(t *testing.T) {
	/* Skipping until we have a blockchain node running as tests will fail at the moment
	s := new(transactionSchedulerTestSuite)
	s.env, s.err = NewIntegrationEnvironment(context.Background())
	if s.err != nil {
		t.Fail()
		return
	}
	suite.Run(t, s)*/
}

func (s *transactionSchedulerTestSuite) TestChainRegistry_HTTPChain() {
	httpSuite := new(HttpChainTestSuite)
	httpSuite.env = s.env
	httpSuite.baseURL = s.env.baseURL
	suite.Run(s.T(), httpSuite)
}

func (s *transactionSchedulerTestSuite) TestChainRegistry_HTTPFaucet() {
	httpSuite := new(HttpFaucetTestSuite)
	httpSuite.env = s.env
	httpSuite.baseURL = s.env.baseURL
	suite.Run(s.T(), httpSuite)
}
