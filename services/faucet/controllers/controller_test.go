package controllers

import (
	"context"
	"math/big"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/faucet"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/faucet/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/types"
)

type MockController struct {
	controls []string
}

func (c *MockController) Control1(f faucet.CreditFunc) faucet.CreditFunc {
	return func(ctx context.Context, r *types.Request) (*big.Int, bool, error) {
		c.controls = append(c.controls, "1")
		// Simulate a valid control
		return f(ctx, r)
	}
}

func (c *MockController) Control2(f faucet.CreditFunc) faucet.CreditFunc {
	return func(ctx context.Context, r *types.Request) (*big.Int, bool, error) {
		c.controls = append(c.controls, "2")
		// Simulate an invalid control
		return big.NewInt(0), false, nil
	}
}

func (c *MockController) Control3(f faucet.CreditFunc) faucet.CreditFunc {
	return func(ctx context.Context, r *types.Request) (*big.Int, bool, error) {
		c.controls = append(c.controls, "3")
		// Simulate a valid control
		return f(ctx, r)
	}
}

func TestCombineControls(t *testing.T) {
	c := MockController{make([]string, 0)}
	crediter := CombineControls(c.Control1, c.Control2, c.Control3)(mock.Credit)
	amount, ok, _ := crediter(context.Background(), &types.Request{})

	if amount.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected amount to be 0 but got %v", amount)
	}

	if ok != false {
		t.Errorf("Expected credited to be invalid")
	}

	if len(c.controls) != 2 {
		t.Errorf("Expected %v controls but got %v", 2, len(c.controls))
	}

	if c.controls[0] != "1" || c.controls[1] != "2" {
		t.Errorf("Expected controls [\"1\", \"2\"] to have been applied but got %v", c.controls)
	}
}

func TestControlledFaucet(t *testing.T) {
	c := MockController{make([]string, 0)}
	f := NewControlledFaucet(&mock.Faucet{}, c.Control1, c.Control2, c.Control3)
	amount, ok, _ := f.Credit(context.Background(), &types.Request{})

	if amount.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected amount to be nil but got %v", amount)
	}

	if ok != false {
		t.Errorf("Expected credited to be invalid")
	}

	if len(c.controls) != 2 {
		t.Errorf("Expected %v controls but got %v", 2, len(c.controls))
	}

	if c.controls[0] != "1" || c.controls[1] != "2" {
		t.Errorf("Expected controls [\"1\", \"2\"] to have been applied but got %v", c.controls)
	}
}
