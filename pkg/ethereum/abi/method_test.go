// +build unit

package abi

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/assert"
)

func TestSignatureToMethod(t *testing.T) {
	emptyMethod, _ := ParseMethod([]byte(`{
		"inputs": [
		],
		"rawName": "empty"
	}`))

	erc20TransferMethod, _ := ParseMethod([]byte(`{
		"inputs": [
			{
				"name": "",
				"type": "address"
			},
			{
				"name": "",
				"type": "uint256"
			}
		],
		"rawName": "transfer"
	}`))

	tests := []struct {
		sig    string
		result *abi.Method
		err    bool
	}{
		// {sig: "FunctionTest(address[3])", result: abi.Method{}, err: nil},
		{sig: "Malformed", result: nil, err: true},
		{sig: "", result: nil, err: true},
		{sig: "()Malformed", result: nil, err: true},
		{sig: "Malformed(,)", result: nil, err: true},
		{sig: "Malformed(address,uint)", result: nil, err: true},
		{sig: "Malformed(unknown)", result: nil, err: true},
		// {sig: "Malformed(address,", result: nil, err: true}, TODO: this test should error
		{sig: "empty()", result: emptyMethod, err: false},
		{sig: "transfer(address,uint256)", result: erc20TransferMethod, err: false},
	}

	for _, test := range tests {
		result, e := ParseMethodSignature(test.sig)
		assert.Equal(t, test.err, e != nil, test.sig, result, e)
		assert.Equal(t, test.result, result)
	}
}
