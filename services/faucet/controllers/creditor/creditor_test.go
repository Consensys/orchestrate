package creditor

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	faucetMock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet/mocks"

	ethcommon "github.com/ethereum/go-ethereum/common"
	faucettypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types/testutils"
)

func TestCreditor(t *testing.T) {
	testSet := []struct {
		name  string
		input *testutils.TestRequest
	}{
		{
			"no faucet candidate",
			&testutils.TestRequest{
				Req: &faucettypes.Request{
					Beneficiary: ethcommon.HexToAddress("0xab"),
				},
				ExpectedErr: errors.FaucetWarning("no faucet candidates").ExtendComponent(component),
			},
		},
		{
			"faucet candidate",
			&testutils.TestRequest{
				Req: &faucettypes.Request{
					Beneficiary: ethcommon.HexToAddress("0xab"),
					FaucetsCandidates: map[string]faucettypes.Faucet{
						"test": {
							Amount:   big.NewInt(10),
							Creditor: ethcommon.HexToAddress("0xab"),
						},
						"test1": {
							Amount:   big.NewInt(11),
							Creditor: ethcommon.HexToAddress("0xcd"),
						},
					},
				},
				ExpectedAmount: big.NewInt(11),
			},
		},
		{
			"creditor should discard all candidates",
			&testutils.TestRequest{
				Req: &faucettypes.Request{
					Beneficiary: ethcommon.HexToAddress("0xab"),
					FaucetsCandidates: map[string]faucettypes.Faucet{
						"test": {
							Amount:   big.NewInt(10),
							Creditor: ethcommon.HexToAddress("0xab"),
						},
						"test1": {
							Amount:   big.NewInt(11),
							Creditor: ethcommon.HexToAddress("0xab"),
						},
					},
				},
				ExpectedErr: errors.FaucetSelfCreditWarning("attempt to credit the creditor").ExtendComponent(component),
			},
		},
	}

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

	cntrl := NewController()
	credit := cntrl.Control(mockFaucet.Credit)

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			amount, err := credit(context.Background(), test.input.Req)
			test.input.ResultAmount, test.input.ResultErr = amount, err
			testutils.AssertRequest(t, test.input)
		})
	}
}
