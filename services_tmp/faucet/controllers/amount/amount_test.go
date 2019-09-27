package amount

import (
	"context"
	"math/big"
	"sync"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/faucet/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/types/testutils"
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
	conf := &Config{
		Amount: big.NewInt(10),
	}
	cntrl := NewController(conf)
	credit := cntrl.Control(mock.Credit)

	// Prepare test data
	rounds := 600
	tests := make([]*testutils.TestRequest, 0)
	for i := 0; i < rounds; i++ {
		tests = append(
			tests,
			&testutils.TestRequest{
				Req: &types.Request{
					ChainID:     chains[i%2],
					Beneficiary: addresses[i%3],
					Amount:      big.NewInt(20),
				},
				ExpectedOK:     true,
				ExpectedAmount: big.NewInt(10),
				ExpectedErr:    false,
			},
		)
	}

	// Apply tests
	wg := &sync.WaitGroup{}
	for _, test := range tests {
		wg.Add(1)
		go func(test *testutils.TestRequest) {
			defer wg.Done()
			test.ResultAmount, test.ResultOK, test.ResultErr = credit(context.Background(), test.Req)
		}(test)
	}
	wg.Wait()

	// Ensure results are correct
	for _, test := range tests {
		testutils.AssertRequest(t, test)
	}
}
