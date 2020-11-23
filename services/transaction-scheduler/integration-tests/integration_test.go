// +build integration

package integrationtests

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	integrationtest "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/integration-test"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

type txSchedulerTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
	err error
}

func (s *txSchedulerTestSuite) SetupSuite() {
	err := integrationtest.StartEnvironment(s.env.ctx, s.env)
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
	ctx, cancel := context.WithCancel(context.Background())

	s.env, s.err = NewIntegrationEnvironment(ctx)
	if s.err != nil {
		t.Errorf(s.err.Error())
		return
	}

	sig := utils.NewSignalListener(func(signal os.Signal) {
		cancel()
	})
	defer sig.Close()

	suite.Run(t, s)
}

func (s *txSchedulerTestSuite) TestTxScheduler_Transactions() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(txSchedulerTransactionTestSuite)
	testSuite.env = s.env
	testSuite.baseURL = s.env.baseURL
	suite.Run(s.T(), testSuite)
}

func (s *txSchedulerTestSuite) TestTxScheduler_Metrics() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(txSchedulerMetricsTestSuite)
	testSuite.env = s.env
	testSuite.baseURL = s.env.baseURL
	testSuite.metricsURL = s.env.metricsURL
	suite.Run(s.T(), testSuite)
}
