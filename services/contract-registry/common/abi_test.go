package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/abi"
)

var ERC20 = []byte(
	`[{
    "anonymous": false,
    "inputs": [
      {"indexed": true, "name": "account", "type": "address"},
      {"indexed": false, "name": "account2", "type": "address"}
    ],
    "name": "MinterAdded",
    "type": "event"
  },
  {
    "inputs": [
      {"indexed": true, "name": "account", "type": "address"},
      {"indexed": true, "name": "account2", "type": "address"}
    ],
    "name": "MinterAdded2",
    "type": "event"
    }]`)

func TestGetIndexedCount(t *testing.T) {

	var ERC20Contract = &abi.Contract{
		Id: &abi.ContractId{
			Name: "ERC20",
			Tag:  "v1.0.0",
		},
		Abi:              ERC20,
		Bytecode:         []byte{1, 2},
		DeployedBytecode: []byte{1, 2, 3},
	}
	erc20ABI, err := ERC20Contract.ToABI()
	assert.NoError(t, err, "should not error on toABI()")

	expected := map[string]uint{
		"MinterAdded":  1,
		"MinterAdded2": 2,
		"Unknown":      0,
	}
	for i, e := range erc20ABI.Events {
		c := GetIndexedCount(e)
		assert.Equal(t, expected[i], c)
	}
}
