package creditor

import (
	"context"
	"math/big"
	"sync"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types/testutils"
)

var (
	chains    = []*big.Int{big.NewInt(10), big.NewInt(2345)}
	addresses = []ethcommon.Address{
		ethcommon.HexToAddress("0xab"),
		ethcommon.HexToAddress("0xcd"),
		ethcommon.HexToAddress("0xef"),
	}
)

func TestCreditor(t *testing.T) {
	// Create Controller and set creditors
	cntrl := NewController()
	for i := range chains {
		cntrl.SetCreditor(chains[i], addresses[i])
	}
	credit := cntrl.Control(mock.Credit)

	// Prepare test data
	rounds := 50
	tests := make([]*testutils.TestRequest, 0)
	for i := 0; i < rounds; i++ {
		tests = append(
			tests,
			&testutils.TestRequest{
				Req: &types.Request{
					ChainID:     chains[i%2],
					Beneficiary: addresses[i%3],
					Amount:      big.NewInt(0),
				},
				ExpectedAmount: big.NewInt(0),
				ExpectedErr:    i%6 == 0 || i%6 == 1,
			},
		)
	}

	// Apply tests
	wg := &sync.WaitGroup{}
	for _, test := range tests {
		wg.Add(1)
		go func(test *testutils.TestRequest) {
			defer wg.Done()
			test.ResultAmount, test.ResultErr = credit(context.Background(), test.Req)
		}(test)
	}
	wg.Wait()

	// Ensure results are correct
	for _, test := range tests {
		testutils.AssertRequest(t, test)
	}
}

func TestCreditorNoCreditor(t *testing.T) {
	// Create Controller and set creditors
	cntrl := NewController()
	credit := cntrl.Control(mock.Credit)

	test := &testutils.TestRequest{
		Req: &types.Request{
			ChainID:     chains[0],
			Beneficiary: addresses[0],
			Amount:      big.NewInt(10),
		},
		ExpectedAmount: big.NewInt(0),
		ExpectedErr:    true,
	}
	test.ResultAmount, test.ResultErr = credit(context.Background(), test.Req)

	testutils.AssertRequest(t, test)
}
