// +build unit

package contracts

import (
	"sort"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

//nolint
var contractAddress = ethcommon.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18")

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
	parsedABI, _ := abi.JSON(strings.NewReader(ERC20))
	ERC20Contract := &entities.Contract{
		Name:             "ERC20",
		Tag:              "v1.0.0",
		RawABI:           ERC20,
		ABI:              parsedABI,
		Bytecode:         hexutil.MustDecode(hexutil.Encode([]byte{1, 2})),
		DeployedBytecode: hexutil.MustDecode(hexutil.Encode([]byte{1, 2, 3})),
	}

	expected := map[string]uint{
		"MinterAdded":  1,
		"MinterAdded2": 2,
		"Unknown":      0,
	}
	for i, e := range ERC20Contract.ABI.Events {
		c := getIndexedCount(&e)
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
