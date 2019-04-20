package sarama

import (
	"context"
	"math/big"
	"sync"
	"testing"

	"github.com/Shopify/sarama/mocks"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/types/testutils"
)

func TestFaucet(t *testing.T) {
	p := mocks.NewSyncProducer(t, nil)
	f := NewFaucet(p)

	rounds := 500
	tests := make([]*testutils.TestCreditData, 0)
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
			&testutils.TestCreditData{
				Req:            r,
				ExpectedOK:     true,
				ExpectedAmount: big.NewInt(20),
				ExpectedErr:    nil,
			},
		)

	}

	// Apply tests
	wg := &sync.WaitGroup{}
	for _, test := range tests {
		wg.Add(1)
		go func(test *testutils.TestCreditData) {
			defer wg.Done()
			amount, ok, err := f.Credit(context.Background(), test.Req)
			test.ResultAmount, test.ResultOK, test.ResultErr = amount, ok, err
		}(test)
	}
	wg.Wait()

	// Ensure results are correct
	for _, test := range tests {
		testutils.AssertCreditData(t, test)
	}
}
