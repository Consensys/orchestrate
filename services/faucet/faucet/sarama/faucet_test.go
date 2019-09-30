package sarama

import (
	"context"
	"math/big"
	"sync"
	"testing"

	"github.com/Shopify/sarama/mocks"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/types/testutils"
)

func TestFaucet(t *testing.T) {
	p := mocks.NewSyncProducer(t, nil)
	f := NewFaucet(p)

	rounds := 500
	tests := make([]*testutils.TestRequest, 0)
	for i := 0; i < rounds; i++ {
		r := &types.Request{
			ChainID:     big.NewInt(10),
			Beneficiary: ethcommon.HexToAddress("0xab"),
			Creditor:    ethcommon.HexToAddress("0xcd"),
			Amount:      big.NewInt(20),
		}

		p.ExpectSendMessageAndSucceed()
		tests = append(
			tests,
			&testutils.TestRequest{
				Req:            r,
				ExpectedOK:     true,
				ExpectedAmount: big.NewInt(20),
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
			amount, ok, err := f.Credit(context.Background(), test.Req)
			test.ResultAmount, test.ResultOK, test.ResultErr = amount, ok, err
		}(test)
	}
	wg.Wait()

	// Ensure results are correct
	for _, test := range tests {
		testutils.AssertRequest(t, test)
	}
}
