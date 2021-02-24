package parsers

import (
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/types/testutils"
	"github.com/stretchr/testify/assert"
)

func TestChainsParser(t *testing.T) {
	chain := testutils.FakeChain()
	chainModel := NewChainModelFromEntity(chain)
	finalChain := NewChainFromModel(chainModel)

	assert.Equal(t, chain, finalChain)
}
