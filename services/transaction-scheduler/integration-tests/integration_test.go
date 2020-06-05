// +build integration

package integrationtests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	integrationtest "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/integration-test"
)

type txSchedulerTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
	err error
}

func (s *txSchedulerTestSuite) SetupSuite() {
	err := integrationtest.StartEnvironment(s.env)
	if err != nil {
		s.env.logger.WithError(err).Error()
		if s.err == nil {
			s.err = err
		}
		return
	}

	s.env.logger.Infof("setup test suite has completed")
}

func (s *txSchedulerTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())

	if s.err != nil {
		s.Fail(s.err.Error())
	}
}

func TestTxScheduler(t *testing.T) {
	s := new(txSchedulerTestSuite)
	s.env, s.err = NewIntegrationEnvironment(context.Background())
	if s.err != nil {
		t.Fail()
		return
	}

	suite.Run(t, s)
}

func (s *txSchedulerTestSuite) TestTxScheduler_Transaction() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(txSchedulerTransactionTestSuite)
	testSuite.env = s.env
	testSuite.baseURL = s.env.baseURL
	suite.Run(s.T(), testSuite)
}
