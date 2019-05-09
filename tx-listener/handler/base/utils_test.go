package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePosition(t *testing.T) {
	chain, position, err := ParsePosition("42:genesis")
	assert.Nil(t, err, "#1: No error expected")
	assert.Equal(t, "42", chain, "#1: Correct chain ID expected")
	assert.Equal(t, int64(0), position.BlockNumber, "#1: Correct blockNumber expected")
	assert.Equal(t, int64(0), position.TxIndex, "#1: Correct txIndex expected")

	chain, position, err = ParsePosition("11:24-124")
	assert.Nil(t, err, "#2: No error expected")
	assert.Equal(t, "11", chain, "#2: Correct chain ID expected")
	assert.Equal(t, int64(24), position.BlockNumber, "#2: Correct blockNumber expected")
	assert.Equal(t, int64(124), position.TxIndex, "#2: Correct txIndex expected")

	chain, position, err = ParsePosition("3:latest-0")
	assert.Nil(t, err, "#3: No error expected")
	assert.Equal(t, "3", chain, "#3: Correct chain ID expected")
	assert.Equal(t, int64(-1), position.BlockNumber, "#3: Correct blockNumber expected")
	assert.Equal(t, int64(0), position.TxIndex, "#3: Correct txIndex expected")

	chain, position, err = ParsePosition("0x3:latest-0")
	assert.Equal(t, "", chain, "#4: Correct chain ID expected")
	assert.NotNil(t, err, "#4: Error expected")
}
