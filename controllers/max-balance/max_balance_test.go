package maxbalance

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/faucet"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/faucet/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/faucet/testutils"
)

var (
	chains    = []*big.Int{big.NewInt(10), big.NewInt(2345), big.NewInt(1)}
	addresses = []ethcommon.Address{
		ethcommon.HexToAddress("0xab"),
		ethcommon.HexToAddress("0xcd"),
		ethcommon.HexToAddress("0xef"),
	}
	values = []*big.Int{big.NewInt(9), big.NewInt(11), big.NewInt(10)}
)

var errTest = fmt.Errorf("Could not connect")

func MockBalanceAt(ctx context.Context, chainID *big.Int, a ethcommon.Address, blocknumber *big.Int) (*big.Int, error) {
	if chainID.Cmp(chains[2]) == 0 {
		// Simulate error
		return nil, errTest
	}
	return big.NewInt(10), nil
}

func TestMaxBalance(t *testing.T) {
	// Create CoolDown controlled credit
	conf := &Config{
		BalanceAt:  MockBalanceAt,
		MaxBalance: big.NewInt(20),
	}
	ctrl := NewController(conf)
	credit := ctrl.Control(mock.Credit)

	// Prepare test data
	rounds := 600
	tests := make([]*testutils.TestCreditData, 0)
	for i := 0; i < rounds; i++ {
		var expectedAmount *big.Int
		var expectedErr error
		switch i % 3 {
		case 0:
			expectedAmount = big.NewInt(9)
		case 1:
			expectedAmount = big.NewInt(0)
		case 2:
			expectedAmount = big.NewInt(0)
			expectedErr = errTest
		}

		tests = append(
			tests,
			&testutils.TestCreditData{
				Req: &faucet.Request{
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
		go func(test *testutils.TestCreditData) {
			defer wg.Done()
			amount, ok, err := credit(context.Background(), test.Req)
			test.ResultAmount, test.ResultOK, test.ResultErr = amount, ok, err
		}(test)
	}
	wg.Wait()

	// Ensure results are correct
	for _, test := range tests {
		testutils.AssertCreditData(t, test)
	}
}
