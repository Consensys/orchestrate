package memory

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/offset/testutils"
)

type ManagerTestSuite struct {
	testutils.OffsetManagerTestSuite
}

func (s *ManagerTestSuite) SetupTest() {
	s.Manager = NewManager()
}

func TestMemory(t *testing.T) {
	s := new(ManagerTestSuite)
	suite.Run(t, s)
}
