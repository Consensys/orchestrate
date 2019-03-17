package handlers

import (
	"fmt"
	"sync"
	"testing"

	log "github.com/sirupsen/logrus"	
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
	abipb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/abi"
	commonpb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/ethereum"
)

type MockABIRegistry struct{}

func (r *MockABIRegistry) GetMethodByID(ID string) (abi.Method, error) {
	if ID == "unknown@" {
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

func (r *MockABIRegistry) RegisterContract(contract *abipb.Contract) error {
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

func makeCrafterContext(i int) *worker.Context {
	ctx := worker.NewContext()
	ctx.Reset()
	ctx.Logger = log.NewEntry(log.StandardLogger())
	ctx.T.Tx = &ethpb.Transaction{TxData: &ethpb.TxData{}}

	switch i % 5 {
	case 0:
		ctx.Keys["errors"] = 0
		ctx.Keys["result"] = ""
	case 1:
		ctx.T.Tx.TxData = (&ethpb.TxData{}).SetData(hexutil.MustDecode("0xa9059cbb"))
		ctx.Keys["errors"] = 0
		ctx.Keys["result"] = "0xa9059cbb"
	case 2:
		ctx.T.Call = &commonpb.Call{
			Method: &abipb.Method{Name: "unknown"},
		}
		ctx.T.Call.Args = []string{"test"}
		ctx.Keys["errors"] = 1
		ctx.Keys["result"] = ""
	case 3:
		ctx.T.Call = &commonpb.Call{
			Method: &abipb.Method{Name: "known"},
		}
		ctx.T.Call.Args = []string{"test"}
		ctx.Keys["errors"] = 0
		ctx.Keys["result"] = payload
	case 4:
		ctx.T.Call = &commonpb.Call{
			Method: &abipb.Method{Name: "known"},
		}
		ctx.Keys["errors"] = 1
		ctx.Keys["result"] = ""
	}
	return ctx
}

func TestCrafter(t *testing.T) {
	// Create crafter handler
	crafter := Crafter(&MockABIRegistry{}, &MockCrafter{})

	rounds := 100
	outs := make(chan *worker.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeCrafterContext(i)
		go func(ctx *worker.Context) {
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

		if out.T.Tx.TxData.GetData() != result {
			t.Errorf("Crafter: expected Data to be %v but got %v", result, out.T.Tx.TxData.GetData())
		}
	}
}
