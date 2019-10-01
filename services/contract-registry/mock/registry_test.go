package mock

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/contract-registry/testutils"
)

type MockTestSuite struct {
	testutils.ContractRegistryTestSuite
}

func (s *MockTestSuite) SetupTest() {
	s.R = NewRegistry()
}

func TestMock(t *testing.T) {
	s := new(MockTestSuite)
	suite.Run(t, s)
}
