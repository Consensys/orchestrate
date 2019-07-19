package chain

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateChainInt(t *testing.T) {
	chain := CreateChainInt(42)
	assert.Equal(t, int64(42), chain.ID().Int64(), "#1: Chain ID should match")

	chain = CreateChainInt(1)
	assert.Equal(t, int64(1), chain.ID().Int64(), "#2: Chain ID should match")
}

func TestCreateChainBigInt(t *testing.T) {
	chain := CreateChainBigInt(big.NewInt(54))
	assert.Equal(t, int64(54), chain.ID().Int64(), "#4: Chain ID should match")

	chain.SetID(big.NewInt(54))
	assert.Equal(t, []byte{0x36}, chain.Id, "#3: Chain ID should have be correct")
}

func TestCreateChainString(t *testing.T) {
	chain := CreateChainString("54")
	assert.Equal(t, int64(54), chain.ID().Int64(), "#4: Chain ID should match")
}
