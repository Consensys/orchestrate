// +build unit

package memory

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/testutils"
)

type MemoryEnvelopeStoreTestSuite struct {
	testutils.EnvelopeStoreTestSuite
}

func (s *MemoryEnvelopeStoreTestSuite) SetupTest() {
	s.Store = New()
}

func TestMemory(t *testing.T) {
	s := new(MemoryEnvelopeStoreTestSuite)
	suite.Run(t, s)
}
