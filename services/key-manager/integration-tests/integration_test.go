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

type keyManagerTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
	err error
}

func (s *keyManagerTestSuite) SetupSuite() {
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

func (s *keyManagerTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())

	if s.err != nil {
		s.Fail(s.err.Error())
	}
}

func TestKeyManager(t *testing.T) {
	s := new(keyManagerTestSuite)
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

func (s *keyManagerTestSuite) TestTxScheduler_Ethereum() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(keyManagerEthereumTestSuite)
	testSuite.env = s.env
	testSuite.baseURL = s.env.baseURL
	suite.Run(s.T(), testSuite)
}
