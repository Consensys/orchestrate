// +build integration

package integrationtests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	integrationtest "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/integration-test"
)

type envelopeStoreTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
	err error
}

func (s *envelopeStoreTestSuite) SetupSuite() {
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

func (s *envelopeStoreTestSuite) TearDownSuite() {
	s.env.Teardown(context.Background())
	if s.err != nil {
		s.Fail(s.err.Error())
	}
}

func TestEnvelopeStore(t *testing.T) {
	s := new(envelopeStoreTestSuite)
	s.env, s.err = NewIntegrationEnvironment(context.Background())
	if s.err != nil {
		t.Fail()
		return
	}
	suite.Run(t, s)
}

func (s *envelopeStoreTestSuite) TestEnvelopeStore_GRPC() {
	if s.err != nil {
		s.env.logger.Warn("skipping test...")
		return
	}

	grpcSuite := new(EnvelopeStoreTestSuite)
	grpcSuite.env = s.env
	grpcSuite.baseURL = s.env.baseURL
	suite.Run(s.T(), grpcSuite)
}
