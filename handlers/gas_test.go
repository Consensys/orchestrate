package handlers

import (
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

type MockGasEstimator struct {
	t *testing.T
}

func (e *MockGasEstimator) EstimateGas(chainID *big.Int, call ethereum.CallMsg) (uint64, error) {
	if chainID.Text(10) == "0" {
		return 0, fmt.Errorf("Could not estimate gas")
	}
	return 18, nil
}

func makeGasEstimatorContext(i int) *types.Context {
	ctx := types.NewContext()
	ctx.Reset()
	switch i % 2 {
	case 0:
		ctx.T.Chain().ID = big.NewInt(0)
		ctx.Keys["errors"] = 1
		ctx.Keys["result"] = uint64(0)
	case 1:
		ctx.T.Chain().ID = big.NewInt(1)
		ctx.Keys["errors"] = 0
		ctx.Keys["result"] = uint64(18)
	}
	return ctx
}

func TestGasEstimator(t *testing.T) {
	me := MockGasEstimator{t: t}
	estimator := GasEstimator(&me)

	rounds := 100
	outs := make(chan *types.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeGasEstimatorContext(i)
		go func(ctx *types.Context) {
			defer wg.Done()
			estimator(ctx)
			outs <- ctx
		}(ctx)
	}
	wg.Wait()
	close(outs)
	if len(outs) != rounds {
		t.Errorf("GasEstimator: expected %v outs but got %v", rounds, len(outs))
	}

	for out := range outs {
		errCount, result := out.Keys["errors"].(int), out.Keys["result"].(uint64)
		if len(out.T.Errors) != errCount {
			t.Errorf("GasEstimator: expected %v errors but got %v", errCount, out.T.Errors)
		}

		if out.T.Tx().GasLimit() != result {
			t.Errorf("GasEstimator: expected gas limit %v but got %v", result, out.T.Tx().GasLimit())
		}
	}
}

type MockGasPricer struct {
	t *testing.T
}

func (e *MockGasPricer) SuggestGasPrice(chainID *big.Int) (*big.Int, error) {
	if chainID.Text(10) == "0" {
		return big.NewInt(0), fmt.Errorf("Could not estimate gas")
	}
	return big.NewInt(10), nil
}

func makeGasPricerContext(i int) *types.Context {
	ctx := types.NewContext()
	ctx.Reset()
	switch i % 2 {
	case 0:
		ctx.T.Chain().ID = big.NewInt(0)
		ctx.Keys["errors"] = 1
		ctx.Keys["result"] = int64(0)
	case 1:
		ctx.T.Chain().ID = big.NewInt(1)
		ctx.Keys["errors"] = 0
		ctx.Keys["result"] = int64(10)
	}
	return ctx
}

func TestGasPricer(t *testing.T) {
	mp := MockGasPricer{t: t}
	pricer := GasPricer(&mp)

	rounds := 100
	outs := make(chan *types.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeGasPricerContext(i)
		go func(ctx *types.Context) {
			defer wg.Done()
			pricer(ctx)
			outs <- ctx
		}(ctx)
	}
	wg.Wait()
	close(outs)
	if len(outs) != rounds {
		t.Errorf("GasEstimator: expected %v outs but got %v", rounds, len(outs))
	}

	for out := range outs {
		errCount, result := out.Keys["errors"].(int), out.Keys["result"].(int64)
		if len(out.T.Errors) != errCount {
			t.Errorf("GasPricer: expected %v errors but got %v", errCount, out.T.Errors)
		}

		if out.T.Tx().GasPrice().Int64() != result {
			t.Errorf("GasPricer: expected gas price %v but got %v", result, out.T.Tx().GasLimit())
		}
	}
}
