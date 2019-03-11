package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	chain := &Chain{
		Id:       "0x2a",
		IsEIP155: true,
	}

	assert.Equal(t, int64(42), chain.ID().Int64(), "Chain ID should match")
}
