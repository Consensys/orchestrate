package redis

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry/testutils"
)

type MockTestSuite struct {
	testutils.ContractRegistryTestSuite
}

func (s *MockTestSuite) SetupTest() {
	s.R = NewRegistry(NewPool(Config(), DialMock))
}

func TestMock(t *testing.T) {
	s := new(MockTestSuite)
	suite.Run(t, s)
}
