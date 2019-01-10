package ethereum

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

var ERC20TransferABI = []byte(`{
	"constant": false,
	"inputs": [
		{
			"name": "_to",
			"type": "address"
		},
		{
			"name": "_value",
			"type": "uint256"
		}
	],
	"name": "transfer",
	"outputs": [
		{
			"name": "",
			"type": "bool"
		}
	],
	"payable": false,
	"stateMutability": "nonpayable",
	"type": "function"
}`)

func TestDummyRegistry(t *testing.T) {
	registry := NewDummyABIRegistry(ERC20TransferABI)

	if m, _ := registry.GetMethodByID("id"); hexutil.Encode(m.Id()) != "0xa9059cbb" {
		t.Errorf("DummyRegistry: expected to retrieve method by ID %q but got %q", "a9059cbb", hexutil.Encode(m.Id()))
	}

	if m, _ := registry.GetMethodBySig("0xa9059cbb"); hexutil.Encode(m.Id()) != "0xa9059cbb" {
		t.Errorf("DummyRegistry: expected to retrieve method by sig %q but got %q", "a9059cbb", hexutil.Encode(m.Id()))
	}
}
