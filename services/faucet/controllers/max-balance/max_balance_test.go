package maxbalance

import (
	"context"
	"math/big"
	"sync"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/faucet/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/types/testutils"
)

var (
	chains = []*big.Int{big.NewInt(10), big.NewInt(2345), big.NewInt(1)}
	values = []*big.Int{big.NewInt(9), big.NewInt(11), big.NewInt(10)}
)

func MockBalanceAt(ctx context.Context, chainID *big.Int, _ ethcommon.Address, _ *big.Int) (*big.Int, error) {
	if chainID.Cmp(chains[2]) == 0 {
		// Simulate error
		return nil, errors.ConnectionError("balanceAtError")
	}
	return big.NewInt(10), nil
}

func TestMaxBalance(t *testing.T) {
	// Create CoolDown controlled credit
	conf := &Config{
		BalanceAt:  MockBalanceAt,
		MaxBalance: big.NewInt(20),
	}
	c := NewController(conf)
	credit := c.Control(mock.Credit)

	// Prepare test data
	rounds := 600
	tests := make([]*testutils.TestRequest, 0)
	for i := 0; i < rounds; i++ {
		var expectedAmount *big.Int
		var expectedErr bool
		switch i % 3 {
		case 0:
			expectedAmount = big.NewInt(9)
		case 1:
			expectedAmount = big.NewInt(0)
			expectedErr = true
		case 2:
			expectedAmount = big.NewInt(0)
			expectedErr = true
		}

		tests = append(
			tests,
			&testutils.TestRequest{
				Req: &types.Request{
					ChainID: chains[i%3],
					Amount:  values[i%3],
				},
				ExpectedOK:     i%3 == 0,
				ExpectedAmount: expectedAmount,
				ExpectedErr:    expectedErr,
			},
		)
	}

	// Apply tests
	wg := &sync.WaitGroup{}
	for _, test := range tests {
		wg.Add(1)
		go func(test *testutils.TestRequest) {
			defer wg.Done()
			amount, ok, err := credit(context.Background(), test.Req)
			test.ResultAmount, test.ResultOK, test.ResultErr = amount, ok, err
		}(test)
	}
	wg.Wait()

	// Ensure results are correct
	for _, test := range tests {
		testutils.AssertRequest(t, test)
		if test.ResultErr != nil {
			assert.True(t, errors.IsFaucetWarning(test.ResultErr) || errors.IsConnectionError(test.ResultErr), "%v should be a faucet warning", test.ResultErr)
			assert.Equal(t, "controller.max-balance", errors.FromError(test.ResultErr).GetComponent(), "Error component should be correct")
		}
	}
}
