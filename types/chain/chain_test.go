package chain

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromInt(t *testing.T) {
	chain := FromInt(42)
	assert.Equal(t, int64(42), chain.GetBigChainID().Int64(), "#1: Chain UUID should match")

	chain = FromInt(1)
	assert.Equal(t, int64(1), chain.GetBigChainID().Int64(), "#2: Chain UUID should match")
}

func TestFromBigInt(t *testing.T) {
	chain := FromBigInt(big.NewInt(54))
	assert.Equal(t, int64(54), chain.GetBigChainID().Int64(), "#3: Chain UUID should match")

	chain.SetChainID(big.NewInt(54))
	assert.Equal(t, "54", chain.ChainId, "#4: Chain UUID should have be correct")
}

func TestFromString(t *testing.T) {
	chain := FromString("54")
	assert.Equal(t, int64(54), chain.GetBigChainID().Int64(), "#5: Chain UUID should match")

	assert.Panics(t, func() { FromString("boom") }, "#6: Chain UUID shouldn't parse")
}
