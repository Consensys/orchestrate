// +build integration

package integrationtests

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	integrationtest "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/integration-test"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type identityManagerTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
	err error
}

func (s *identityManagerTestSuite) SetupSuite() {
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

func (s *identityManagerTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())

	if s.err != nil {
		s.Fail(s.err.Error())
	}
}

func TestIdentityManager(t *testing.T) {
	s := new(identityManagerTestSuite)
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

func (s *identityManagerTestSuite) TestIdentityManager_Accounts() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(identityManagerTransactionTestSuite)
	testSuite.env = s.env
	testSuite.baseURL = s.env.baseURL
	suite.Run(s.T(), testSuite)
}
