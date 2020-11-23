// +build integration

package integrationtests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	integrationtest "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/integration-test"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

type txSignerTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
	err error
}

func (s *txSignerTestSuite) SetupSuite() {
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

func (s *txSignerTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())

	if s.err != nil {
		s.Fail(s.err.Error())
	}
}

func TestTxSigner(t *testing.T) {
	s := new(txSignerTestSuite)
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

func (s *txSignerTestSuite) TestTxSigner_Ethereum() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(txSignerEthereumTestSuite)
	testSuite.env = s.env

	time.Sleep(3 * time.Second)
	suite.Run(s.T(), testSuite)
}
