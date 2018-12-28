package handlers

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

type MockEthCrediter struct {
	count int
	t     *testing.T
}

func (c *MockEthCrediter) Credit(chainID *big.Int, a common.Address, value *big.Int) error {
	if chainID.Text(10) == "0" {
		return fmt.Errorf("Could not credit")
	}
	c.count++
	return nil
}

type MockEthCreditController struct {
	t *testing.T
}

var blackAddress = "0xdbb881a51CD4023E4400CEF3ef73046743f08da3"

func (c *MockEthCreditController) ShouldCredit(chainID *big.Int, a common.Address, value *big.Int) (*big.Int, bool) {
	if a.Hex() == blackAddress {
		return nil, false
	}
	return big.NewInt(100), true
}

func TestFaucet(t *testing.T) {
	// Create Faucet handler
	c := &MockEthCrediter{t: t}
	faucet := Faucet(c, &MockEthCreditController{t: t})

	ctx := infra.NewContext()
	ctx.Reset()
	ctx.T.Chain().ID = big.NewInt(0)
	faucet(ctx)
	if len(ctx.T.Errors) != 1 {
		t.Errorf("Faucet 1: expected 1 error but got %v", ctx.T.Errors)
	}

	if c.count != 0 {
		t.Errorf("Faucet 1: expected credit count to be 0 but got %v", c.count)
	}

	ctx = infra.NewContext()
	ctx.Reset()
	ctx.T.Chain().ID = big.NewInt(0)
	*ctx.T.Sender().Address = common.HexToAddress(blackAddress)
	faucet(ctx)
	if len(ctx.T.Errors) != 0 {
		t.Errorf("Faucet 2: expected 0 error but got %v", ctx.T.Errors)
	}
	if c.count != 0 {
		t.Errorf("Faucet 2: expected credit count to be 0 but got %v", c.count)
	}

	ctx = infra.NewContext()
	ctx.Reset()
	ctx.T.Chain().ID = big.NewInt(1)
	faucet(ctx)
	if len(ctx.T.Errors) != 0 {
		t.Errorf("Faucet 3: expected no error but got %v", ctx.T.Errors)
	}
	if c.count != 1 {
		t.Errorf("Faucet 3: expected credit count to be 1 but got %v", c.count)
	}
}
