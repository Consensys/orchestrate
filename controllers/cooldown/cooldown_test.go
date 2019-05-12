package cooldown

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/faucet/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/types/testutils"
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

func TestCoolDown(t *testing.T) {
	// Create CoolDown controlled credit
	conf := &Config{
		Delay:   time.Duration(10 * time.Millisecond),
		Stripes: 2,
	}
	ctrl := NewController(conf)
	credit := ctrl.Control(mock.Credit)

	// Prepare test data
	rounds := 600
	tests := make([]*testutils.TestRequest, 0)
	for i := 0; i < rounds; i++ {
		var expectedAmount *big.Int
		if i%6 < 3 {
			expectedAmount = big.NewInt(10)
		} else {
			expectedAmount = big.NewInt(0)
		}
		tests = append(
			tests,
			&testutils.TestRequest{
				Req: &types.Request{
					ChainID:     chains[i%3],
					Beneficiary: addresses[i%3],
					Amount:      big.NewInt(10),
				},
				ExpectedOK:     i%6 < 3,
				ExpectedAmount: expectedAmount,
				ExpectedErr:    nil,
			},
		)
	}

	// Apply tests
	wg := &sync.WaitGroup{}
	for i, test := range tests {
		wg.Add(1)
		go func(test *testutils.TestRequest) {
			defer wg.Done()
			test.ResultAmount, test.ResultOK, test.ResultErr = credit(context.Background(), test.Req)
		}(test)
		switch i % 6 {
		case 2:
			// Sleeps half cooldown time
			time.Sleep(5 * time.Millisecond)
		case 5:
			// Sleep to cooldown delay on controller
			time.Sleep(10 * time.Millisecond)
		}
	}
	wg.Wait()

	// Ensure results are correct
	for _, test := range tests {
		testutils.AssertRequest(t, test)
	}
}
