package abi

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

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

func shouldMatch(f interface{}, input string, expect string, t *testing.T) {
	var s string
	switch f := f.(type) {
	case func(string) (abi.Method, error):
		method, _ := f(input)
		s = method.String()
	case func(string) (abi.Event, error):
		event, _ := f(input)
		s = event.String()
	}
	if s != expect {
		t.Errorf("%v: wrong abi returned, expected %v, got %v", input, expect, s)
	}
}

func shouldError(f interface{}, input string, t *testing.T) {
	var err error
	switch f := f.(type) {
	case func(string) (abi.Method, error):
		_, err = f(input)
	case func(string) (abi.Event, error):
		_, err = f(input)
	}

	if err == nil {
		t.Errorf("%v id tag should have returned an error", input)
	}
}

func TestRegisterContract(t *testing.T) {
	r := NewContractABIRegistry()
	err := r.RegisterContract("ERC1400", []byte(""))
	if err == nil {
		t.Error("Should have returned an error because wrong ABI structure")
	}

	err = r.RegisterContract("ERC1400", ERC1400)
	if err != nil {
		t.Errorf("Got error %v when registering contract", err)
	}
}

func TestContractRegistryByID(t *testing.T) {
	r := NewContractABIRegistry()
	r.RegisterContract("ERC1400", ERC1400)

	method := r.GetMethodByID
	shouldMatch(method, "isMinter@ERC1400", "function isMinter(address account) constant returns(bool)", t)
	shouldError(method, "isMinters@ERC1400", t)
	shouldError(method, "is@Minters@ERC1400", t)
	shouldError(method, "isMinter@ERC1401", t)
	shouldError(method, "isMinters@ERC1401", t)

	event := r.GetEventByID
	shouldMatch(event, "MinterAdded@ERC1400", "event MinterAdded(address indexed account)", t)
	shouldError(event, "MinterAdd@ERC1400", t)
	shouldError(event, "is@MinterAdded@ERC1400", t)
	shouldError(event, "MinterAdded@ERC1401", t)
	shouldError(event, "MinterAdded@ERC1401", t)
}

func TestContractRegistryBySig(t *testing.T) {
	r := NewContractABIRegistry()
	r.RegisterContract("ERC1400", ERC1400)

	method := r.GetMethodBySig
	shouldMatch(method, "0xaa271e1a", "function isMinter(address account) constant returns(bool)", t)
	shouldMatch(method, "aa271e1a", "function isMinter(address account) constant returns(bool)", t)
	shouldError(method, "0xaa271e1ab", t)
	shouldError(method, "0xaa271e1b", t)
	shouldError(method, "wrong", t)

	event := r.GetEventBySig
	sig := "6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6"
	sig0x := "0x" + sig
	shouldMatch(event, sig, "event MinterAdded(address indexed account)", t)
	shouldMatch(event, sig0x, "event MinterAdded(address indexed account)", t)
	shouldError(event, "6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f", t)
	shouldError(event, sig[:63]+"a", t)
	shouldError(event, "wrong", t)
}
