// +build integration

package integrationtests

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
)

type transactionSchedulerTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
}

func TestTransactionScheduler(t *testing.T) {
	s := new(transactionSchedulerTestSuite)
	s.env = NewIntegrationEnvironment(context.Background())
	suite.Run(t, s)
}

func (s *transactionSchedulerTestSuite) SetupSuite() {
	s.env.Start()
}

func (s *transactionSchedulerTestSuite) TearDownSuite() {
	s.env.Teardown()
}

func (s *transactionSchedulerTestSuite) TestTransactionScheduler_Jobs() {
	grpcSuite := new(JobsTestSuite)
	grpcSuite.env = s.env
	grpcSuite.baseURL = "localhost:8081/jobs"
	suite.Run(s.T(), grpcSuite)
}
