package abi

import (
	"testing"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/assert"
)

var raw = []byte(`
[
	{
		"constant": false,
		"inputs": [
			{
				"name": "to",
				"type": "address"
			}
		],
		"name": "delegate",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "winningProposal",
		"outputs": [
			{
				"name": "_winningProposal",
				"type": "uint8"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "toVoter",
				"type": "address"
			}
		],
		"name": "giveRightToVote",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "toProposal",
				"type": "uint8"
			}
		],
		"name": "vote",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"name": "_numProposals",
				"type": "uint8"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "fallback"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "voter",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "vote",
				"type": "uint8"
			}
		],
		"name": "Voted",
		"type": "event"
	}
]`)

func TestMarshalABI(t *testing.T) {
	// Prepare reference ABI object
	abi1 := &ethabi.ABI{}
	_ = abi1.UnmarshalJSON(raw)

	// Marshal ABI object into byte
	b, err := MarshalABI(abi1)
	assert.Nil(t, err, "MarshalABI should not error")

	// Unmarshal result byte into a second ABI object
	abi2 := &ethabi.ABI{}
	err = abi2.UnmarshalJSON(b)
	assert.Nil(t, err, "Unmarshal should not error")

	// Check constructor is equal
	assertEqualMethod(t, abi1.Constructor, abi2.Constructor)

	// Check methods are the same
	assert.Equal(t, len(abi1.Methods), len(abi2.Methods), "Same count of methods expected")
	for name, m1 := range abi1.Methods {
		m2, ok := abi2.Methods[name]
		assert.True(t, ok, "Should have a method names %q", name)
		assertEqualMethod(t, m1, m2)
	}

	// Check events are the same
	assert.Equal(t, len(abi1.Events), len(abi2.Events), "Same count of events expected")
	for name, e1 := range abi1.Events {
		e2, ok := abi2.Events[name]
		assert.True(t, ok, "Should have a method names %q", name)
		assertEqualEvent(t, e1, e2)
	}
}

func assertEqualMethod(t *testing.T, m1, m2 ethabi.Method) {
	assert.Equal(t, m1.Name, m2.Name, "Method Name should be equal")
	assert.Equal(t, m1.Const, m2.Const, "Method constant should be equal")
	assert.Equal(t, len(m1.Inputs), len(m2.Inputs), "Method constant should have same count of inputs")
	assert.Equal(t, len(m1.Outputs), len(m2.Outputs), "Method constant should have same count of outputs")
	assert.Equal(t, m1.String(), m2.String(), "Method String should be equal")
	assert.Equal(t, m1.Sig(), m2.Sig(), "Method Sig should be equal")
}

func assertEqualEvent(t *testing.T, e1, e2 ethabi.Event) {
	assert.Equal(t, e1.Name, e2.Name, "Method Name should be equal")
	assert.Equal(t, e1.Anonymous, e2.Anonymous, "Method anonymous should be equal")
	assert.Equal(t, len(e1.Inputs), len(e2.Inputs), "Method constant should have same count of inputs")
	assert.Equal(t, e1.String(), e2.String(), "Event String should be equal")
	assert.Equal(t, e1.ID().Hex(), e2.ID().Hex(), "Event Sig should be equal")
}
