package memory

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/testutils"
)

type MemoryEnvelopeStoreTestSuite struct {
	testutils.EnvelopeStoreTestSuite
}

func (s *MemoryEnvelopeStoreTestSuite) SetupTest() {
	s.Store = NewEnvelopeStore()
}

func TestMemory(t *testing.T) {
	s := new(MemoryEnvelopeStoreTestSuite)
	suite.Run(t, s)
}
