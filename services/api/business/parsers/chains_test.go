package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
)

func TestChainsParser(t *testing.T) {
	chain := testutils.FakeChain()
	chainModel := NewChainModelFromEntity(chain)
	finalChain := NewChainFromModel(chainModel)

	assert.Equal(t, chain, finalChain)
}
