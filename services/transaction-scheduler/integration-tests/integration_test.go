// +build integration

package integrationtests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type txSchedulerTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
}

func (s *txSchedulerTestSuite) SetupSuite() {
	// s.env.Start()
}

func (s *txSchedulerTestSuite) TearDownSuite() {
	// s.env.Teardown()
}

func TestTxScheduler(t *testing.T) {
	s := new(txSchedulerTestSuite)
	s.env = NewIntegrationEnvironment(context.Background())
	suite.Run(t, s)
}

func (s *txSchedulerTestSuite) TestTxScheduler_Transaction() {
	testSuite := new(txSchedulerTransactionTestSuite)
	testSuite.env = s.env
	testSuite.baseURL = "http://localhost:8081"
	suite.Run(s.T(), testSuite)
}
