package static

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
)

type sigTest struct {
	contract string
	sig      string
	result   string
	err      error
}

type selectorTest struct {
	selector string
	result   string
	err      error
}

var ERC1400 = []byte(
	`[{
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "name": "account",
        "type": "address"
      }
    ],
    "name": "MinterAdded",
    "type": "event"
  },
  {
    "constant": true,
    "inputs": [
      {
        "name": "account",
        "type": "address"
      }
    ],
    "name": "isMinter",
    "outputs": [
      {
        "name": "",
        "type": "bool"
      }
    ],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
    }]`)

var ERC1400Contract = &abi.Contract{Name: "ERC1400", Tag: "", Abi: ERC1400, Bytecode: []byte{}}

func TestRegisterContract(t *testing.T) {
	r := NewRegistry()
	err := r.RegisterContract(&abi.Contract{Name: "ERC1400", Tag: "", Abi: []byte{}, Bytecode: []byte{}})
	assert.NoError(t, err, "Should not error on empty ABI")

	err = r.RegisterContract(ERC1400Contract)
	assert.NoError(t, err, "Got error %v when registering contract")
}

func TestContractRegistryBySig(t *testing.T) {
	r := NewRegistry()
	err := r.RegisterContract(ERC1400Contract)
	assert.NoError(t, err)

	tests := []sigTest{
		{contract: "ERC1400", sig: "isMinter(address)", result: "function isMinter(address account) constant returns(bool)", err: nil},
		{contract: "ERC1400", sig: "constructor()", result: "function () returns()", err: nil},
		{contract: "ERC1400", sig: "isMinters(address)", result: "", err: errors.New("")},
		{contract: "ERC1400", sig: "is()Minters()", result: "", err: errors.New("")},
		{contract: "ERC1401", sig: "isMinter(address)", result: "", err: errors.New("")},
		{contract: "ERC1401", sig: "isMinters(address)", result: "", err: errors.New("")},
	}
	for _, test := range tests {
		result, resErr := r.GetMethodBySig(test.contract, test.sig)
		assert.IsType(t, test.err, resErr)
		if resErr == nil {
			assert.Equal(t, test.result, result.String())
		}
	}

	tests = []sigTest{
		{contract: "ERC1400", sig: "MinterAdded(address)", result: "event MinterAdded(address indexed account)", err: nil},
		{contract: "ERC1400", sig: "MinterAdd(address)", result: "", err: errors.New("")},
		{contract: "ERC1400", sig: "is()MinterAdded", result: "", err: errors.New("")},
		{contract: "ERC1401", sig: "MinterAdded(address)", result: "", err: errors.New("")},
	}
	for _, test := range tests {
		result, resErr := r.GetEventBySig(test.contract, test.sig)
		assert.IsType(t, test.err, resErr)
		if resErr == nil {
			assert.Equal(t, test.result, result.String())
		}
	}

	r = NewRegistry()
	err = r.RegisterContract(&abi.Contract{Name: "ERC1400", Tag: "v0.1.1", Abi: ERC1400, Bytecode: []byte{}})
	assert.NoError(t, err)

	tests = []sigTest{
		{contract: "ERC1400[v0.1.1]", sig: "isMinter(address)", result: "function isMinter(address account) constant returns(bool)", err: nil},
		{contract: "ERC1400[v0.1.1]", sig: "constructor()", result: "function () returns()", err: nil},
		{contract: "ERC1400[v0.1.1]", sig: "isMinters()", result: "", err: errors.New("")},
		{contract: "ERC1400[v0.1.1]", sig: "is()Minters()", result: "", err: errors.New("")},
		{contract: "ERC1401[v0.1.1]", sig: "isMinter())", result: "", err: errors.New("")},
		{contract: "ERC1401[v0.1.1]", sig: "isMinters())", result: "", err: errors.New("")},
	}
	for _, test := range tests {
		result, resErr := r.GetMethodBySig(test.contract, test.sig)
		assert.IsType(t, test.err, resErr)
		if resErr == nil {
			assert.Equal(t, test.result, result.String())
		}
	}

	tests = []sigTest{
		{contract: "ERC1400[v0.1.1]", sig: "MinterAdded(address)", result: "event MinterAdded(address indexed account)", err: nil},
		{contract: "ERC1400[v0.1.1]", sig: "MinterAdd(address)", result: "", err: errors.New("")},
		{contract: "ERC1400[v0.1.1]", sig: "is()MinterAdded", result: "", err: errors.New("")},
		{contract: "ERC1401[v0.1.1]", sig: "MinterAdded(address)", result: "", err: errors.New("")},
	}
	for _, test := range tests {
		result, resErr := r.GetEventBySig(test.contract, test.sig)
		assert.IsType(t, test.err, resErr)
		if resErr == nil {
			assert.Equal(t, test.result, result.String())
		}
	}
}

func TestContractRegistryBySelector(t *testing.T) {
	r := NewRegistry()
	err := r.RegisterContract(ERC1400Contract)
	assert.NoError(t, err)

	tests := []selectorTest{
		{selector: "0xaa271e1a", result: "function isMinter(address account) constant returns(bool)", err: nil},
		{selector: "aa271e1a", result: "function isMinter(address account) constant returns(bool)", err: nil},
		{selector: "0xaa271e1ab", result: "", err: hexutil.ErrSyntax},
		{selector: "0xaa271e1b", result: "", err: errors.New("")},
		{selector: "wrong", result: "", err: hexutil.ErrSyntax},
	}
	for _, test := range tests {
		result, resErr := r.GetMethodBySelector(test.selector)
		assert.IsType(t, test.err, resErr)
		if resErr == nil {
			assert.Equal(t, test.result, result.String())
		}
	}

	sig := "6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6"
	sig0x := "0x" + sig

	tests = []selectorTest{
		{selector: sig, result: "event MinterAdded(address indexed account)", err: nil},
		{selector: sig0x, result: "event MinterAdded(address indexed account)", err: nil},
		{selector: sig[:63], result: "", err: hexutil.ErrSyntax},
		{selector: sig[:63] + "a", result: "", err: errors.New("")},
		{selector: "wrong", result: "", err: hexutil.ErrSyntax},
	}
	for _, test := range tests {
		result, resErr := r.GetEventBySelector(test.selector)
		assert.IsType(t, test.err, resErr)
		if resErr == nil {
			assert.Equal(t, test.result, result.String())
		}
	}
}
