package cooldown

import (
	"context"
	"math/big"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	faucetMock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet/mocks"
	faucettypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types/testutils"
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
	// Create Controller and set creditors
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockFaucet := faucetMock.NewMockFaucet(mockCtrl)
	mockFaucet.EXPECT().Credit(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, r *faucettypes.Request) (*big.Int, error) {
		if len(r.FaucetsCandidates) == 0 {
			return nil, errors.FaucetWarning("no faucet request").ExtendComponent(component)
		}
		// Select a first faucet candidate for comparison
		r.ElectedFaucet = reflect.ValueOf(r.FaucetsCandidates).MapKeys()[0].String()
		for key, candidate := range r.FaucetsCandidates {
			if candidate.Amount.Cmp(r.FaucetsCandidates[r.ElectedFaucet].Amount) == -1 {
				r.ElectedFaucet = key
			}
		}
		return r.FaucetsCandidates[r.ElectedFaucet].Amount, nil
	}).AnyTimes()

	// Create CoolDown controlled credit
	cntrl := NewController()
	credit := cntrl.Control(mockFaucet.Credit)

	// Prepare test data
	rounds := 50
	tests := make([]*testutils.TestRequest, 0)
	for i := 0; i < rounds; i++ {
		var expectedAmount *big.Int
		if i%6 < 3 {
			expectedAmount = big.NewInt(10)
		} else {
			expectedAmount = nil
		}
		tests = append(
			tests,
			&testutils.TestRequest{
				Req: &faucettypes.Request{
					ChainID:     chains[i%3],
					Beneficiary: addresses[i%3],
					FaucetsCandidates: map[string]faucettypes.Faucet{
						"faucetID": {
							Amount:     big.NewInt(10),
							MaxBalance: big.NewInt(10),
							Cooldown:   100 * time.Millisecond,
						},
					},
				},
				ExpectedAmount: expectedAmount,
				ExpectedErr: func(i int) error {
					if i%6 >= 3 {
						return errors.FaucetWarning("faucet cooling down").ExtendComponent(component)
					}
					return nil
				}(i),
			},
		)
	}

	// Apply tests
	wg := &sync.WaitGroup{}
	for i, test := range tests {
		wg.Add(1)
		go func(test *testutils.TestRequest) {
			defer wg.Done()
			amount, err := credit(context.Background(), test.Req)
			test.ResultAmount, test.ResultErr = amount, err
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
			assert.Equal(t, "faucet.controllers.cooldown", errors.FromError(test.ResultErr).GetComponent(), "Error component should be correct")
		}
	}
}
