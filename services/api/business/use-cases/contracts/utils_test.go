package contracts

import (
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

//nolint
var contractAddress = "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"
//nolint
var chainID = "chainId"

var ERC20 = `[{
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
    }]`

func TestGetIndexedCount(t *testing.T) {
	var ERC20Contract = &entities.Contract{
		ID: entities.ContractID{
			Name: "ERC20",
			Tag:  "v1.0.0",
		},
		ABI:              ERC20,
		Bytecode:         hexutil.Encode([]byte{1, 2}),
		DeployedBytecode: hexutil.Encode([]byte{1, 2, 3}),
	}
	erc20ABI, err := ERC20Contract.ToABI()
	assert.NoError(t, err, "should not error on toABI()")

	expected := map[string]uint{
		"MinterAdded":  1,
		"MinterAdded2": 2,
		"Unknown":      0,
	}
	for i, e := range erc20ABI.Events {
		c := getIndexedCount(e)
		assert.Equal(t, expected[i], c)
	}
}

func TestSortStrings(t *testing.T) {
	tests := []struct {
		name string
		args []string
		res  []string
	}{
		{"base", []string{"z", "Z", "a", "A"}, []string{"A", "a", "Z", "z"}},
		{"opposite", []string{"Z", "z", "A", "a"}, []string{"A", "a", "Z", "z"}},
		{"bien", []string{"encore du travail", "1", "2", ".", "🛠"}, []string{".", "1", "2", "encore du travail", "🛠"}},
		{"empty", []string{}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(utils.Alphabetic(tt.args))
			assert.Equal(t, tt.res, tt.args)
		})
	}
}
