package memory

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/testutils"
)

type MemoryTestSuite struct {
	testutils.ContractRegistryTestSuite
}

func (s *MemoryTestSuite) SetupTest() {
	s.R = NewRegistry()
}

func TestMemory(t *testing.T) {
	s := new(MemoryTestSuite)
	suite.Run(t, s)
}
