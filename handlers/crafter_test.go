package handlers

import (
	"fmt"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

type MockABIRegistry struct{}

func (r *MockABIRegistry) GetMethodByID(ID string) (abi.Method, error) {
	if ID == "unknown" {
		return abi.Method{}, fmt.Errorf("Could not retrieve ABI")
	}
	return abi.Method{}, nil
}

func (r *MockABIRegistry) GetMethodBySig(sig string) (abi.Method, error) {
	return abi.Method{}, nil
}

func (r *MockABIRegistry) GetEventByID(ID string) (abi.Event, error) {
	if ID == "unknown" {
		return abi.Event{}, fmt.Errorf("Could not retrieve ABI")
	}
	return abi.Event{}, nil
}

func (r *MockABIRegistry) GetEventBySig(sig string) (abi.Event, error) {
	return abi.Event{}, nil
}

func (r *MockABIRegistry) RegisterContract(sig string, abi []byte) error {
	return nil
}

type MockCrafter struct{}

var payload = "0xa9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000"

func (c *MockCrafter) Craft(method abi.Method, args ...string) ([]byte, error) {
	if len(args) != 1 {
		return []byte(``), fmt.Errorf("Could not craft")
	}
	return hexutil.MustDecode(payload), nil
}

func makeCrafterContext(i int) *types.Context {
	ctx := types.NewContext()
	ctx.Reset()
	switch i % 5 {
	case 0:
		ctx.Keys["errors"] = 0
		ctx.Keys["result"] = "0x"
	case 1:
		ctx.T.Tx().SetData(hexutil.MustDecode("0xa9059cbb"))
		ctx.Keys["errors"] = 0
		ctx.Keys["result"] = "0xa9059cbb"
	case 2:
		ctx.T.Call().MethodID = "unknown"
		ctx.T.Call().Args = []string{"test"}
		ctx.Keys["errors"] = 1
		ctx.Keys["result"] = "0x"
	case 3:
		ctx.T.Call().MethodID = "known"
		ctx.T.Call().Args = []string{"test"}
		ctx.Keys["errors"] = 0
		ctx.Keys["result"] = payload
	case 4:
		ctx.T.Call().MethodID = "known"
		ctx.Keys["errors"] = 1
		ctx.Keys["result"] = "0x"
	}
	return ctx
}

func TestCrafter(t *testing.T) {
	// Create crafter handler
	crafter := Crafter(&MockABIRegistry{}, &MockCrafter{})

	rounds := 100
	outs := make(chan *types.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeCrafterContext(i)
		go func(ctx *types.Context) {
			defer wg.Done()
			crafter(ctx)
			outs <- ctx
		}(ctx)
	}
	wg.Wait()
	close(outs)

	if len(outs) != rounds {
		t.Errorf("Crafter: expected %v outs but got %v", rounds, len(outs))
	}

	for out := range outs {
		errCount, result := out.Keys["errors"].(int), out.Keys["result"].(string)
		if len(out.T.Errors) != errCount {
			t.Errorf("Crafter: expected %v errors but got %v", errCount, out.T.Errors)
		}

		if hexutil.Encode(out.T.Tx().Data()) != result {
			t.Errorf("Crafter: expected Dara to be %v but got %v", result, hexutil.Encode(out.T.Tx().Data()))
		}
	}
}
