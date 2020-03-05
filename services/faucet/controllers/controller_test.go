package controllers

import (
	"context"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	faucetMock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet/mocks"

	"github.com/stretchr/testify/assert"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
)

type MockController struct {
	controls []string
}

func (c *MockController) Control1(f faucet.CreditFunc) faucet.CreditFunc {
	return func(ctx context.Context, r *types.Request) (*big.Int, error) {
		c.controls = append(c.controls, "1")
		// Simulate a valid control
		return f(ctx, r)
	}
}

func (c *MockController) Control2(f faucet.CreditFunc) faucet.CreditFunc {
	return func(ctx context.Context, r *types.Request) (*big.Int, error) {
		c.controls = append(c.controls, "2")
		// Simulate an invalid control
		return big.NewInt(0), errors.FaucetWarning("invalid control")
	}
}

func (c *MockController) Control3(f faucet.CreditFunc) faucet.CreditFunc {
	return func(ctx context.Context, r *types.Request) (*big.Int, error) {
		c.controls = append(c.controls, "3")
		// Simulate a valid control
		return f(ctx, r)
	}
}

func TestCombineControls(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockFaucet := faucetMock.NewMockFaucet(mockCtrl)
	c := MockController{make([]string, 0)}
	creditor := CombineControls(c.Control1, c.Control2, c.Control3)(mockFaucet.Credit)
	amount, err := creditor(context.Background(), &types.Request{})
	assert.Error(t, err, "Expected credited to be invalid")
	assert.Equal(t, 0, amount.Cmp(big.NewInt(0)), "Wrong credit amount")
	assert.Len(t, c.controls, 2, "Wrong controls")
	assert.True(t, c.controls[0] == "1" && c.controls[1] == "2", "Expected controls [\"1\", \"2\"] to have been applied but got %v", c.controls)
}

func TestControlledFaucet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockFaucet := faucetMock.NewMockFaucet(mockCtrl)
	c := MockController{make([]string, 0)}
	f := NewControlledFaucet(mockFaucet, c.Control1, c.Control2, c.Control3)
	amount, err := f.Credit(context.Background(), &types.Request{})
	assert.Error(t, err, "Expected credited to be invalid")
	assert.Equal(t, 0, amount.Cmp(big.NewInt(0)), "Wrong credit amount")
	assert.Len(t, c.controls, 2, "Wrong controls")
	assert.True(t, c.controls[0] == "1" && c.controls[1] == "2", "Expected controls [\"1\", \"2\"] to have been applied but got %v", c.controls)
}
