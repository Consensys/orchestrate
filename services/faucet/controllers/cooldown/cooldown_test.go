package cooldown

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/faucet/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/types/testutils"
)

var (
	chains    = []*big.Int{big.NewInt(10), big.NewInt(2345), big.NewInt(1)}
	addresses = []ethcommon.Address{
		ethcommon.HexToAddress("0xab"),
		ethcommon.HexToAddress("0xcd"),
		ethcommon.HexToAddress("0xef"),
	}
)

func TestCoolDown(t *testing.T) {
	// Create CoolDown controlled credit
	conf := &Config{
		Delay:   100 * time.Millisecond,
		Stripes: 2,
	}
	cntrl := NewController(conf)
	credit := cntrl.Control(mock.Credit)

	// Prepare test data
	rounds := 60
	tests := make([]*testutils.TestRequest, 0)
	for i := 0; i < rounds; i++ {
		var expectedAmount *big.Int
		var expectedErr bool
		if i%6 < 3 {
			expectedAmount = big.NewInt(10)
			expectedErr = false
		} else {
			expectedAmount = big.NewInt(0)
			expectedErr = true
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
				ExpectedErr:    expectedErr,
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
			time.Sleep(50 * time.Millisecond)
		case 5:
			// Sleep to cooldown delay on controller
			time.Sleep(100 * time.Millisecond)
		}
	}
	wg.Wait()

	// Ensure results are correct
	for _, test := range tests {
		testutils.AssertRequest(t, test)
		if test.ResultErr != nil {
			assert.True(t, errors.IsFaucetWarning(test.ResultErr), "%v should be a faucet warning", test.ResultErr)
			assert.Equal(t, "controller.cooldown", errors.FromError(test.ResultErr).GetComponent(), "Error component should be correct")
		}
	}
}
