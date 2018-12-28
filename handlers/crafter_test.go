package handlers

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

type MockABIRegistry struct{}

func (r *MockABIRegistry) GetMethodByID(ID string) (*abi.Method, error) {
	if ID == "unknown" {
		return nil, fmt.Errorf("Could not retrieve ABI")
	}
	return &abi.Method{}, nil
}

func (r *MockABIRegistry) GetMethodBySig(sig string) (*abi.Method, error) {
	return &abi.Method{}, nil
}

type MockCrafter struct{}

var payload = []byte(`0xa9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000`)

func (c *MockCrafter) Craft(method *abi.Method, args ...string) ([]byte, error) {
	if len(args) != 1 {
		return []byte(``), fmt.Errorf("Could not craft")
	}
	return payload, nil
}

func TestCrafter(t *testing.T) {
	// Create crafter handler
	c := Crafter(&MockABIRegistry{}, &MockCrafter{})

	ctx := infra.NewContext()
	ctx.Reset()
	c(ctx)
	if len(ctx.T.Tx().Data()) != 0 {
		t.Errorf("Crafter: expected data no be empty but got %q", hexutil.Encode(ctx.T.Tx().Data()))
	}

	ctx = infra.NewContext()
	ctx.Reset()
	ctx.T.Tx().SetData([]byte(`a9059cbb`))
	c(ctx)
	if len(ctx.T.Tx().Data()) != 8 {
		t.Errorf("Crafter: expected data to be %q but got %q", "0xa9059cbb", hexutil.Encode(ctx.T.Tx().Data()))
	}

	ctx = infra.NewContext()
	ctx.Reset()
	ctx.T.Call().MethodID = "unknown"
	ctx.T.Call().Args = []string{"test"}
	c(ctx)
	if len(ctx.T.Tx().Data()) != 0 {
		t.Errorf("Crafter: expected data no be empty but got %q", hexutil.Encode(ctx.T.Tx().Data()))
	}
	if len(ctx.T.Errors) != 1 {
		t.Errorf("Crafter: expected an error")
	}

	ctx = infra.NewContext()
	ctx.Reset()
	ctx.T.Call().MethodID = "known"
	ctx.T.Call().Args = []string{"test"}
	c(ctx)
	if hexutil.Encode(ctx.T.Tx().Data()) != hexutil.Encode(payload) {
		t.Errorf("Crafter: expected data no be %v but got %q", hexutil.Encode(payload), hexutil.Encode(ctx.T.Tx().Data()))
	}

	ctx = infra.NewContext()
	ctx.Reset()
	ctx.T.Call().MethodID = "known"
	c(ctx)
	if len(ctx.T.Tx().Data()) != 0 {
		t.Errorf("Crafter: expected data no be empty but got %q", hexutil.Encode(ctx.T.Tx().Data()))
	}
	if len(ctx.T.Errors) != 1 {
		t.Errorf("Crafter: expected an error")
	}
}
