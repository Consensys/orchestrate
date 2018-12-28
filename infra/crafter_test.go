package infra

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func testBindArg(stringKind string, arg string, t *testing.T) interface{} {
	boundArg, err := bindArg(stringKind, arg)
	if err != nil {
		t.Errorf("%q expected to be compatible with type %q", arg, stringKind)
	}
	return boundArg
}

func TestBindArg(t *testing.T) {
	a := "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"
	addr := testBindArg("address", a, t).(common.Address)
	if addr.Hex() != a {
		t.Errorf("Expect bind %q but got %q", a, addr.Hex())
	}

	dec := testBindArg("int", "0x400", t).(*big.Int)
	if dec.Int64() != 1024 {
		t.Errorf("Expect bind to %v but got %v", 1024, dec.Int64())
	}

	b := testBindArg("bool", "0x1", t).(bool)
	if !b {
		t.Errorf("Expect bind to %v but got %v", true, false)
	}
}

func newMethod(methodABI []byte) *abi.Method {
	var method abi.Method
	json.Unmarshal(methodABI, &method)
	return &method
}

var ERC20TransferMethod = newMethod([]byte(`{
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
}`))

var CustomMethod = newMethod([]byte(`{
	"constant": false,
	"inputs": [
		{
			"name": "_address",
			"type": "address"
		},
		{
			"name": "_bytesA",
			"type": "bytes"
		},
		{
			"name": "_uint256",
			"type": "uint256"
		},
		{
			"name": "_uint17",
			"type": "uint17"
		},
		{
			"name": "_bool",
			"type": "bool"
		},
		{
			"name": "_bytesB",
			"type": "bytes"
		}
	],
	"name": "custom",
	"outputs": [
		{
			"name": "",
			"type": "bool"
		}
	],
	"payable": true,
	"stateMutability": "nonpayable",
	"type": "function"
}`))

func TestBindArgs(t *testing.T) {
	var (
		_to    = "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"
		_value = "0x2386f26fc10000"
	)
	_, err := bindArgs(ERC20TransferMethod, _to, _value)

	if err != nil {
		t.Errorf("Prepare Args: should prepare args")
	}

	var (
		_address = "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"
		_bytesA  = "0x2386f26fc10000"
		_uint256 = "0x6009608a02a7a15fd6689d6dad560c44e9ab61ff"
		_uint17  = "0xdd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775af"
		_bool    = "0x1"
		_bytesB  = "0xa1a45fabb381e6ab02448013f651fa0792c3fa05b38771f161cb8f7ebdbee973b5"
	)
	_, err = bindArgs(CustomMethod, _address, _bytesA, _uint256, _uint17, _bool, _bytesB)

	if err != nil {
		t.Errorf("Prepare Args: should prepare args")
	}
}

func TestPayloadCraft(t *testing.T) {
	c := PayloadCrafter{}
	var (
		_to     = "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"
		_value  = "0x2386f26fc10000"
		payload = "0xa9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000"
	)
	data, err := c.Craft(ERC20TransferMethod, _to, _value)

	if err != nil {
		t.Errorf("Craft: received error %q ", err)
	}

	if hexutil.Encode(data) != payload {
		t.Errorf("Craft: expected payload %q but got %q", payload, hexutil.Encode(data))
	}

	var (
		_address = "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"
		_bytesA  = "0x2386f26fc10000"
		_uint256 = "0x6009608a02a7a15fd6689d6dad560c44e9ab61ff"
		_uint17  = "0xdd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775af"
		_bool    = "0x1"
		_bytesB  = "0xa1a45fabb381e6ab02448013f651fa0792c3fa05b38771f161cb8f7ebdbee973b5"
	)

	payload = "0xa8817683000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb68400000000000000000000000000000000000000000000000000000000000000c00000000000000000000000006009608a02a7a15fd6689d6dad560c44e9ab61ff000000000000000000000dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775af0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000072386f26fc10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000021a1a45fabb381e6ab02448013f651fa0792c3fa05b38771f161cb8f7ebdbee973b500000000000000000000000000000000000000000000000000000000000000"
	data, err = c.Craft(CustomMethod, _address, _bytesA, _uint256, _uint17, _bool, _bytesB)

	if err != nil {
		t.Errorf("Craft: received error %q ", err)
	}

	if hexutil.Encode(data) != payload {
		t.Errorf("Craft: expected payload %q but got %q", payload, hexutil.Encode(data))
	}

}
