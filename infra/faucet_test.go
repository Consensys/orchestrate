package infra

import (
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

var testControllerCfg = &SimpleCreditControllerConfig{
	BalanceAt: func(chainID *big.Int, a common.Address) (*big.Int, error) {
		if a.Hex() == "0xdbb881a51CD4023E4400CEF3ef73046743f08da3" {
			return big.NewInt(2000), nil
		}
		return big.NewInt(1000), nil
	},
	CreditAmount: big.NewInt(1000),
	MaxBalance:   big.NewInt(1999),
	CreditDelay:  time.Duration(10 * time.Millisecond),
	BlackList:    map[string]struct{}{"0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff": struct{}{}},
}

func TestSimpleCreditController(t *testing.T) {

	c := NewSimpleCreditController(testControllerCfg, 10)

	// Black listed address should not be credited
	amount, ok := c.ShouldCredit(big.NewInt(10), common.HexToAddress("0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff"), big.NewInt(100))
	if ok || amount != nil {
		t.Errorf("SimpleCreditController: should not credit black listed account")
	}

	// Black listed address should not be credited
	amount, ok = c.ShouldCredit(big.NewInt(10), common.HexToAddress("0xdbb881a51CD4023E4400CEF3ef73046743f08da3"), big.NewInt(100))
	if ok || amount != nil {
		t.Errorf("SimpleCreditController: should not credit account with to high balance")
	}

	// Should credit in nominal case
	amount, ok = c.ShouldCredit(big.NewInt(10), common.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684"), big.NewInt(100))
	if !ok || amount.Int64() != 1000 {
		t.Errorf("SimpleCreditController: should credit in nominal case")
	}

	// Should not credit befor delay
	amount, ok = c.ShouldCredit(big.NewInt(10), common.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684"), big.NewInt(100))
	if ok || amount != nil {
		t.Errorf("SimpleCreditController: should not credit account that has not cooldown")
	}
}

var testCreditControllerData = []struct {
	chainID *big.Int
	a       common.Address
}{
	{big.NewInt(1), common.HexToAddress("0xdbb881a51CD4023E4400CEF3ef73046743f08da3")},
	{big.NewInt(1), common.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684")},
	{big.NewInt(1), common.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684")},
	{big.NewInt(2), common.HexToAddress("0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff")},
	{big.NewInt(1), common.HexToAddress("0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff")},
	{big.NewInt(2), common.HexToAddress("0xdbb881a51CD4023E4400CEF3ef73046743f08da3")},
}

func TestSimpleCreditControllerConcurrent(t *testing.T) {
	c := NewSimpleCreditController(testControllerCfg, 10)

	rounds := 600
	credits := make(chan bool, rounds)
	wg := &sync.WaitGroup{}
	for i := 1; i <= rounds; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, ok := c.ShouldCredit(testCreditControllerData[i%6].chainID, testCreditControllerData[i%6].a, big.NewInt(100))
			// Test as been designed such as 1 out of 6 entry are valid for a credit
			if ok {
				credits <- ok
			}
		}(i)
		if i%6 == 0 {
			// Sleep to cooldown delay on controller
			time.Sleep(10 * time.Millisecond)
		}
	}
	wg.Wait()
	if len(credits) != rounds/6 {
		t.Errorf("SimpleController: expected %v credits but got %v", rounds/6, len(credits))
	}
}
