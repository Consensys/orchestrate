package memory

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/testutils"
)

type ModelsTestSuite struct {
	testutils.ChainRegistryTestSuite
}

func (s *ModelsTestSuite) SetupTest() {
	s.Store = NewChainRegistry()
}

func TestModels(t *testing.T) {
	s := new(ModelsTestSuite)
	suite.Run(t, s)
}
