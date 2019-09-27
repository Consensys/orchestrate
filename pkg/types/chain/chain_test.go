package chain

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromInt(t *testing.T) {
	chain := FromInt(42)
	assert.Equal(t, int64(42), chain.ID().Int64(), "#1: Chain ID should match")

	chain = FromInt(1)
	assert.Equal(t, int64(1), chain.ID().Int64(), "#2: Chain ID should match")
}

func TestFromBigInt(t *testing.T) {
	chain := FromBigInt(big.NewInt(54))
	assert.Equal(t, int64(54), chain.ID().Int64(), "#3: Chain ID should match")

	chain.SetID(big.NewInt(54))
	assert.Equal(t, []byte{0x36}, chain.Id, "#4: Chain ID should have be correct")
}

func TestFromString(t *testing.T) {
	chain := FromString("54")
	assert.Equal(t, int64(54), chain.ID().Int64(), "#5: Chain ID should match")

	assert.Panics(t, func() { FromString("boom") }, "#6: Chain ID shouldn't parse")
}
