package abi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContract(t *testing.T) {
	c := Contract{}
	assert.Equal(t, "", c.Short(), "Short should be empty")

	c.Name = "ERC1400"
	assert.Equal(t, "ERC1400", c.Short(), "Short should be correct")

	c.Tag = "v1.0.1"
	assert.Equal(t, "ERC1400[v1.0.1]", c.Short(), "Short should be correct")
}

func TestFromShortContract(t *testing.T) {
	c, err := FromShortContract("ERC20")
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, "ERC20", c.Name, "Contract should be correct")
	assert.Equal(t, "", c.Tag, "Tag should be correct")

	c, err = FromShortContract("ERC20[v0.1.2-alpha]")
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, "ERC20", c.Name, "Contract should be correct")
	assert.Equal(t, "v0.1.2-alpha", c.Tag, "Tag should be correct")

	c, err = FromShortContract("ERC20[v0.1.2;alpha]")
	assert.NotNil(t, err, "Expected error")
}
