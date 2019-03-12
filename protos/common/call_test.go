package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethod(t *testing.T) {
	m := Method{}
	assert.Equal(t, "", m.StringShort(), "StringShort should be empty")

	m.Contract = "ERC1400"
	m.Name = "transferByPartition"
	assert.Equal(t, "transferByPartition@ERC1400", m.StringShort(), "StringShort should be correct")

	m.Deploy = true
	assert.Equal(t, "deploy(ERC1400)", m.StringShort(), "StringShort should be correct")

	m.Version = "v1.0.1"
	assert.Equal(t, "deploy(ERC1400[v1.0.1])", m.StringShort(), "StringShort should be correct")

	m.Deploy = false
	assert.Equal(t, "transferByPartition@ERC1400[v1.0.1]", m.StringShort(), "StringShort should be correct")
}
