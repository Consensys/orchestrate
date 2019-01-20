package faucet

import (
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
)

type testData struct {
	req            *services.FaucetRequest
	resultAmount   *big.Int
	resultOK       bool
	resultErr      error
	expectedAmount *big.Int
	expectedOK     bool
	expectedErr    error
}

func checkTestData(test *testData, t *testing.T) {
	if test.resultOK != test.expectedOK {
		t.Errorf("Invalid OK: expected %v but got %v", test.expectedOK, test.resultOK)
	}

	if test.resultAmount.Cmp(test.expectedAmount) != 0 {
		t.Errorf("Invalid amount credited: expecte %v but got %v", test.expectedAmount, test.resultAmount)
	}

	if test.resultErr == nil && test.expectedErr != nil || test.resultErr != nil && test.expectedErr == nil {
		t.Errorf("Invalid Error: expected %v but got %v", test.expectedErr, test.resultErr)
	}
}

var (
	chains    = []*big.Int{big.NewInt(10), big.NewInt(2345), big.NewInt(1)}
	addresses = []common.Address{common.HexToAddress("0xab"), common.HexToAddress("0xcd"), common.HexToAddress("0xef")}
	values    = []*big.Int{big.NewInt(9), big.NewInt(11), big.NewInt(10)}
)

func TestBlackList(t *testing.T) {
	// Create BlackList controlled credit
	bl := NewBlackList(chains, addresses)
	credit := bl.Control(MockCredit)

	// Prepare test data
	rounds := 600
	tests := make([]*testData, 0)
	for i := 0; i < rounds; i++ {
		var expectedAmount *big.Int
		if i%2 == 0 {
			expectedAmount = big.NewInt(0)
		} else {
			expectedAmount = big.NewInt(10)
		}
		tests = append(
			tests,
			&testData{
				req: &services.FaucetRequest{
					ChainID: chains[i%3],
					Address: addresses[(i+i%2)%3],
				},
				expectedOK:     i%2 == 1,
				expectedAmount: expectedAmount,
				expectedErr:    nil,
			},
		)
	}

	// Apply tests
	wg := &sync.WaitGroup{}
	for _, test := range tests {
		wg.Add(1)
		go func(test *testData) {
			defer wg.Done()
			amount, ok, err := credit(test.req)
			test.resultAmount, test.resultOK, test.resultErr = amount, ok, err
		}(test)
	}
	wg.Wait()

	// Ensure results are correct
	for _, test := range tests {
		checkTestData(test, t)
	}
}

func TestCoolDown(t *testing.T) {
	// Create CoolDown controlled credit
	cd := NewCoolDown(time.Duration(10*time.Millisecond), 2)
	credit := cd.Control(MockCredit)

	// Prepare test data
	rounds := 600
	tests := make([]*testData, 0)
	for i := 0; i < rounds; i++ {
		var expectedAmount *big.Int
		if i%6 < 3 {
			expectedAmount = big.NewInt(10)
		} else {
			expectedAmount = big.NewInt(0)
		}
		tests = append(
			tests,
			&testData{
				req: &services.FaucetRequest{
					ChainID: chains[i%3],
					Address: addresses[i%3],
				},
				expectedOK:     i%6 < 3,
				expectedAmount: expectedAmount,
				expectedErr:    nil,
			},
		)
	}

	// Apply tests
	wg := &sync.WaitGroup{}
	for i, test := range tests {
		wg.Add(1)
		go func(test *testData) {
			defer wg.Done()
			amount, ok, err := credit(test.req)
			test.resultAmount, test.resultOK, test.resultErr = amount, ok, err
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
		checkTestData(test, t)
	}
}

var errTest = fmt.Errorf("Could not connect")

func MockBalanceAt(chainID *big.Int, a common.Address) (*big.Int, error) {
	if chainID.Cmp(chains[2]) == 0 {
		// Simulate error
		return nil, errTest
	}
	return big.NewInt(10), nil
}

func TestMaxBalance(t *testing.T) {
	// Create CoolDown controlled credit
	mb := NewMaxBalance(big.NewInt(20), MockBalanceAt)
	credit := mb.Control(MockCredit)

	// Prepare test data
	rounds := 600
	tests := make([]*testData, 0)
	for i := 0; i < rounds; i++ {
		var expectedAmount *big.Int
		var expectedErr error
		switch i % 3 {
		case 0:
			expectedAmount = big.NewInt(10)
		case 1:
			expectedAmount = big.NewInt(0)
		case 2:
			expectedAmount = big.NewInt(0)
			expectedErr = errTest
		}

		tests = append(
			tests,
			&testData{
				req: &services.FaucetRequest{
					ChainID: chains[i%3],
					Value:   values[i%3],
				},
				expectedOK:     i%3 == 0,
				expectedAmount: expectedAmount,
				expectedErr:    expectedErr,
			},
		)
	}

	// Apply tests
	wg := &sync.WaitGroup{}
	for _, test := range tests {
		wg.Add(1)
		go func(test *testData) {
			defer wg.Done()
			amount, ok, err := credit(test.req)
			test.resultAmount, test.resultOK, test.resultErr = amount, ok, err
		}(test)
	}
	wg.Wait()

	// Ensure results are correct
	for _, test := range tests {
		checkTestData(test, t)
	}
}
