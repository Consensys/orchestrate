package args

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromShortCall(t *testing.T) {
	c, err := SignatureToCall("transfer()")
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, "transfer", c.GetMethod().GetName(), "Name should be correct")
	assert.False(t, c.GetMethod().IsConstructor(), "Deploy should be correct")

	c, err = SignatureToCall("transfer()")
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, "transfer", c.GetMethod().GetName(), "Name should be correct")
	assert.False(t, c.GetMethod().IsConstructor(), "Deploy should be correct")

	c, err = SignatureToCall("constructor()")
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, "constructor", c.GetMethod().GetName(), "Name should be correct")
	assert.Equal(t, "constructor", c.Short(), "Name should be correct")
	assert.True(t, c.IsConstructor(), "Deploy should be correct")
	assert.True(t, c.GetMethod().IsConstructor(), "Deploy should be correct")

	_, err = SignatureToCall("transfer")
	assert.Error(t, err, "Error expected")

	_, err = SignatureToCall("transfer(toto)")
	assert.Error(t, err, "Error expected")
}
