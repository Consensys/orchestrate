package abi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContract(t *testing.T) {
	c := Contract{
		Id: &ContractId{},
	}
	assert.Equal(t, "", c.Short(), "Short should be empty")

	c.SetName("ERC1400")
	assert.Equal(t, "ERC1400", c.Short(), "Short should be correct")

	c.SetTag("v1.0.1")
	assert.Equal(t, "ERC1400[v1.0.1]", c.Short(), "Short should be correct")
}

func TestFromStringd(t *testing.T) {
	c, err := StringToContract("ERC20")
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, "ERC20", c.GetName(), "Contract should be correct")
	assert.Equal(t, "", c.GetTag(), "Tag should be correct")
	assert.Equal(t, []byte{}, c.Abi, "ABI should be correct")
	assert.Equal(t, []byte{}, c.Bytecode, "Bytecode should be correct")

	c, err = StringToContract("ERC20[v0.1.2-alpha]")
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, "ERC20", c.GetName(), "Contract should be correct")
	assert.Equal(t, "v0.1.2-alpha", c.GetTag(), "Tag should be correct")
	assert.Equal(t, []byte{}, c.Abi, "ABI should be correct")
	assert.Equal(t, []byte{}, c.Bytecode, "Bytecode should be correct")

	c, err = StringToContract("ERC20[v0.1.2-alpha]:[{\"constant\":true,\"inputs\":[],\"name\":\"testMethod\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]")
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, "ERC20", c.GetName(), "Contract should be correct")
	assert.Equal(t, "v0.1.2-alpha", c.GetTag(), "Tag should be correct")
	assert.NotEmpty(t, c.Abi, "ABI should have been registered")
	assert.Equal(t, []byte{}, c.Bytecode, "Bytecode should be correct")

	c, err = StringToContract("ERC20[v0.1.2-alpha]:[{\"constant\":true,\"inputs\":[],\"name\":\"testMethod\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]:0xabcd1234ef")
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, "ERC20", c.GetName(), "Contract should be correct")
	assert.Equal(t, "v0.1.2-alpha", c.GetTag(), "Tag should be correct")
	assert.NotEmpty(t, c.Abi, "ABI should have been registered")
	assert.Equal(t, []byte{0xab, 0xcd, 0x12, 0x34, 0xef}, c.Bytecode, "Bytecode should be correct")

	c, err = StringToContract("ERC20[v0.1.2-alpha]:[{\"constant\":true,\"inputs\":[],\"name\":\"testMethod\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]:0xabcd1234ef:0xabcd1234ef")
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, "ERC20", c.GetName(), "Contract should be correct")
	assert.Equal(t, "v0.1.2-alpha", c.GetTag(), "Tag should be correct")
	assert.NotEmpty(t, c.Abi, "ABI should have been registered")
	assert.Equal(t, []byte{0xab, 0xcd, 0x12, 0x34, 0xef}, c.Bytecode, "Bytecode should be correct")
	assert.Equal(t, []byte{0xab, 0xcd, 0x12, 0x34, 0xef}, c.DeployedBytecode, "DeployedBytecode should be correct")

	gethABI, err := c.ToABI()
	assert.NoError(t, err, "ABI has been properly parsed")
	assert.Len(t, gethABI.Methods, 1, "Method has been registered")
	assert.Equal(t, "testMethod", gethABI.Methods["testMethod"].Name, "method name should match")

	_, err = StringToContract("ERC20[v0.1.2;alpha]")
	assert.Error(t, err, "Expected error")
}
