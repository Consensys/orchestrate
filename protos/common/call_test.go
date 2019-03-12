package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethod(t *testing.T) {
	m := Method{}
	assert.Equal(t, "", m.Short(), "Short should be empty")

	m.Contract = "ERC1400"
	m.Name = "transferByPartition"
	assert.Equal(t, "transferByPartition@ERC1400", m.Short(), "Short should be correct")

	m.Tag = "v1.0.1"
	assert.Equal(t, "transferByPartition@ERC1400[v1.0.1]", m.Short(), "Short should be correct")
}

func TestFromShortMethod(t *testing.T) {
	m, err := FromShortMethod("transfer@ERC20")
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, "transfer", m.Name, "Name should be correct")
	assert.Equal(t, "ERC20", m.Contract, "Contract should be correct")
	assert.Equal(t, "", m.Tag, "Tag should be correct")
	assert.Equal(t, false, m.IsDeploy(), "Deploy should be correct")

	m, err = FromShortMethod("transfer@ERC20[v1.0.1]")
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, "transfer", m.Name, "Name should be correct")
	assert.Equal(t, "ERC20", m.Contract, "Contract should be correct")
	assert.Equal(t, "v1.0.1", m.Tag, "Tag should be correct")
	assert.Equal(t, false, m.IsDeploy(), "Deploy should be correct")

	m, err = FromShortMethod("constructor@ERC20[v1.0.1]")
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, "constructor", m.Name, "Name should be correct")
	assert.Equal(t, "ERC20", m.Contract, "Contract should be correct")
	assert.Equal(t, "v1.0.1", m.Tag, "Tag should be correct")
	assert.Equal(t, true, m.IsDeploy(), "Deploy should be correct")
}
