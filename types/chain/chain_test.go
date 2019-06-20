package chain

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	chain := CreateChainInt(42)
	assert.Equal(t, int64(42), chain.ID().Int64(), "#1: Chain ID should match")

	chain = CreateChainInt(1)
	assert.Equal(t, int64(1), chain.ID().Int64(), "#2: Chain ID should match")

	chain.SetID(big.NewInt(54))
	assert.Equal(t, []byte{0x36}, chain.Id, "#3: Chain ID should have be correct")
}
