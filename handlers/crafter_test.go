package handlers

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/ethereum"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
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

func newCrafterTestMessage(i int) *tracepb.Trace {
	var pb tracepb.Trace
	switch i % 4 {
	case 0:
		// Valid args
		pb.Call = &tracepb.Call{MethodId: "abcde", Args: []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x2386f26fc10000"}}
	case 1:
		// Nothing to craft
		pb.Call = &tracepb.Call{}
	case 2:
		// Already crafted
		pb.Call = &tracepb.Call{MethodId: "abcde", Args: []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x2386f26fc10000"}}
		pb.Transaction = &ethpb.Transaction{Raw: "0xa9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000"}
	case 3:
		// Invalid args
		pb.Call = &tracepb.Call{MethodId: "abcde", Args: []string{}}
	}
	return &pb

}

func TestCrafter(t *testing.T) {
	// Create worker
	w := infra.NewWorker(100)
	w.Use(Loader(&TraceProtoUnmarshaller{}))

	// Create and register crafter handler
	h := Crafter(ERC20Getter)
	w.Use(h)

	mockH := NewMockHandler(50)
	w.Use(mockH.Handler())

	// Create input channel
	in := make(chan interface{})

	// Run worker
	go w.Run(in)

	// Feed input channel and then close it
	rounds := 1000
	for i := 1; i <= rounds; i++ {
		in <- newCrafterTestMessage(i)
	}
	close(in)

	// Wait for worker to be done
	<-w.Done()

	// We expected 1/4 rounds to have aborted
	if len(mockH.handled) != 3*rounds/4 {
		t.Errorf("Crafter: expected %v rounds but got %v", rounds, len(mockH.handled))
	}

	expected := "0xa9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000"
	for _, ctx := range mockH.handled {
		if len(ctx.T.Errors) != 0 && hexutil.Encode(ctx.T.Tx().Data()) != "0x" {
			t.Errorf("Crafter: expected no raw tx on error but got %q", hexutil.Encode(ctx.T.Tx().Data()))
		}
		if len(ctx.T.Errors) == 0 && ctx.T.Call().MethodID != "" && hexutil.Encode(ctx.T.Tx().Data()) != expected {
			t.Errorf("Crafter: expected raw %q but got %q", expected, hexutil.Encode(ctx.T.Tx().Data()))
		}
		if len(ctx.T.Errors) == 0 && ctx.T.Call().MethodID == "" && len(ctx.T.Tx().Data()) != 0 {
			t.Errorf("Crafter: expected raw %q but got %q", expected, hexutil.Encode(ctx.T.Tx().Data()))
		}
	}
}
