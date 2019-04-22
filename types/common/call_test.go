package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromShortCall(t *testing.T) {
	c, err := StringToCall("transfer@ERC20")
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, "transfer", c.GetMethod().GetName(), "Name should be correct")
	assert.Equal(t, "ERC20", c.GetContract().GetName(), "Contract should be correct")
	assert.Equal(t, "", c.GetContract().GetTag(), "Tag should be correct")
	assert.Equal(t, false, c.GetMethod().IsDeploy(), "Deploy should be correct")

	c, err = StringToCall("transfer@ERC20[v1.0.1]")
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, "transfer", c.GetMethod().GetName(), "Name should be correct")
	assert.Equal(t, "ERC20", c.GetContract().GetName(), "Contract should be correct")
	assert.Equal(t, "v1.0.1", c.GetContract().GetTag(), "Tag should be correct")
	assert.Equal(t, false, c.GetMethod().IsDeploy(), "Deploy should be correct")

	c, err = StringToCall("constructor@ERC20[v1.0.1]")
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, "constructor", c.GetMethod().GetName(), "Name should be correct")
	assert.Equal(t, "ERC20", c.GetContract().GetName(), "Contract should be correct")
	assert.Equal(t, "v1.0.1", c.GetContract().GetTag(), "Tag should be correct")
	assert.Equal(t, true, c.GetMethod().IsDeploy(), "Deploy should be correct")

	_, err = StringToCall("transfer")
	assert.NotNil(t, err, "No error expected")

	_, err = StringToCall("transfer@ERC20[v1;0;1]")
	assert.NotNil(t, err, "No error expected")
}
