// +build integration

package integrationtests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type transactionSchedulerTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
}

func TestChainRegistry(t *testing.T) {
	s := new(transactionSchedulerTestSuite)
	s.env = NewIntegrationEnvironment(context.Background())
	suite.Run(t, s)
}

func (s *transactionSchedulerTestSuite) SetupSuite() {
	err := s.env.Start()
	if err != nil {
		s.T().Error(err)
	}
}

func (s *transactionSchedulerTestSuite) TearDownSuite() {
	s.env.Teardown()
}

func (s *transactionSchedulerTestSuite) TestChainRegistry_HTTPChain() {
	httpSuite := new(HttpChainTestSuite)
	httpSuite.env = s.env
	httpSuite.baseURL = "http://localhost:8081"
	suite.Run(s.T(), httpSuite)
}

func (s *transactionSchedulerTestSuite) TestChainRegistry_HTTPFaucet() {
	httpSuite := new(HttpFaucetTestSuite)
	httpSuite.env = s.env
	httpSuite.baseURL = "http://localhost:8081"
	suite.Run(s.T(), httpSuite)
}
