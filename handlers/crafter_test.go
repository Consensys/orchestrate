package handlers

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

var (
	ERC20TransferABI = []byte(`{
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
	ERC20TransferMethod abi.Method
	_                   = json.Unmarshal(ERC20TransferABI, &ERC20TransferMethod)
	ERC20Getter         = NewDummyABIGetter(&ERC20TransferMethod)
)

var errGetABI = errors.New("Could not retrieve ABI")

type ErrorABIGetter struct{}

func (g *ErrorABIGetter) GetMethodByID(ID string) (*abi.Method, error) {
	return nil, errGetABI
}

func TestCrafter(t *testing.T) {
	// Valid craft
	handler := Crafter(ERC20Getter)

	// Create a context
	ctx := infra.NewContext()
	ctx.Init([]infra.HandlerFunc{handler})
	ctx.T.Call().Args = []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x2386f26fc10000"}

	// Execute handler
	ctx.Next()

	if len(ctx.T.Errors) != 0 {
		t.Errorf("Crafter: could not craft %v", ctx.T.Errors)
	}

	payload := "0xa9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000"
	data := ctx.T.Tx().Data()
	if hexutil.Encode(data) != payload {
		t.Errorf("Crafter: expected payload %q but got %q", payload, hexutil.Encode(data))
	}

	// Arguments missing
	ctx.Init([]infra.HandlerFunc{handler})
	ctx.T.Call().Args = []string{}

	// Execute handler
	ctx.Next()

	if len(ctx.T.Errors) != 1 {
		t.Errorf("Crafter: expected to produce error")
	}

	// ABI Getter return Error
	handler = Crafter(&ErrorABIGetter{})
	ctx.Init([]infra.HandlerFunc{handler})
	ctx.T.Call().Args = []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x2386f26fc10000"}

	// Execute handler
	ctx.Next()

	if len(ctx.T.Errors) != 1 {
		t.Errorf("Crafter: expected to produce error")
	}
}
