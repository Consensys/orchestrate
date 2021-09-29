// +build integration

package integrationtests

import (
	"context"
	"os"
	"testing"
	"time"

	integrationtest "github.com/consensys/orchestrate/pkg/toolkit/integration-test"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/stretchr/testify/suite"
)

type txSenderTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
	err error
}

func (s *txSenderTestSuite) SetupSuite() {
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

func (s *txSenderTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())

	if s.err != nil {
		s.Fail(s.err.Error())
	}
}

func TestTxSender(t *testing.T) {
	s := new(txSenderTestSuite)
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

func (s *txSenderTestSuite) TestTxSender_Ethereum() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	testSuite := new(txSenderEthereumTestSuite)
	testSuite.env = s.env

	time.Sleep(3 * time.Second)
	suite.Run(s.T(), testSuite)
}
