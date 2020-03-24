package maxbalance

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	ethClientMock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/mocks"
	faucetMock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet/mocks"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	faucettypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types/testutils"
)

const (
	chainURLBalanceAtError = "error"
	chainURLBalanceAt0     = "0"
	chainURLBalanceAt10    = "10"
)

var (
	errBalanceAt = fmt.Errorf("error")
)

func TestMaxBalance(t *testing.T) {
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
				ExpectedAmount: nil,
				ExpectedErr:    errors.FaucetWarning("no faucet candidates").ExtendComponent(component),
			},
		},
		{
			"BalanceAt error",
			&testutils.TestRequest{
				Req: &faucettypes.Request{
					Beneficiary: ethcommon.HexToAddress("0xab"),
					ChainURL:    chainURLBalanceAtError,
					FaucetsCandidates: map[string]faucettypes.Faucet{
						"test": {
							Amount:     big.NewInt(10),
							MaxBalance: big.NewInt(10),
						},
						"test1": {
							Amount:     big.NewInt(1),
							MaxBalance: big.NewInt(10),
						},
					},
				},
				ExpectedErr: errors.FromError(errBalanceAt).ExtendComponent(component),
			},
		},
		{
			"credit after max balance control",
			&testutils.TestRequest{
				Req: &faucettypes.Request{
					Beneficiary: ethcommon.HexToAddress("0xab"),
					ChainURL:    chainURLBalanceAt0,
					FaucetsCandidates: map[string]faucettypes.Faucet{
						"test": {
							Amount:     big.NewInt(10),
							MaxBalance: big.NewInt(10),
						},
						"test1": {
							Amount:     big.NewInt(1),
							MaxBalance: big.NewInt(10),
						},
					},
				},
				ExpectedAmount: big.NewInt(10),
			},
		},
		{
			"remove all candidates after max balance control",
			&testutils.TestRequest{
				Req: &faucettypes.Request{
					Beneficiary: ethcommon.HexToAddress("0xab"),
					ChainURL:    chainURLBalanceAt10,
					FaucetsCandidates: map[string]faucettypes.Faucet{
						"test": {
							Amount:     big.NewInt(10),
							MaxBalance: big.NewInt(10),
						},
						"test1": {
							Amount:     big.NewInt(1),
							MaxBalance: big.NewInt(10),
						},
					},
				},
				ExpectedErr: errors.FaucetWarning("account balance too high").ExtendComponent(component),
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

	mockEthClient := ethClientMock.NewMockChainStateReader(mockCtrl)
	mockEthClient.EXPECT().BalanceAt(gomock.Any(), gomock.Eq(chainURLBalanceAtError), gomock.Any(), gomock.Any()).Return(nil, errBalanceAt).AnyTimes()
	mockEthClient.EXPECT().BalanceAt(gomock.Any(), gomock.Eq(chainURLBalanceAt0), gomock.Any(), gomock.Any()).Return(big.NewInt(0), nil).AnyTimes()
	mockEthClient.EXPECT().BalanceAt(gomock.Any(), gomock.Eq(chainURLBalanceAt10), gomock.Any(), gomock.Any()).Return(big.NewInt(10), nil).AnyTimes()

	cntrl := NewController(mockEthClient)
	credit := cntrl.Control(mockFaucet.Credit)

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			amount, err := credit(context.Background(), test.input.Req)
			test.input.ResultAmount, test.input.ResultErr = amount, err
			testutils.AssertRequest(t, test.input)
		})
	}
}
