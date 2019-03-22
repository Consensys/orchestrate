package handlers

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

type MockGasEstimator struct {
	t *testing.T
}

func (e *MockGasEstimator) EstimateGas(ctx context.Context, chainID *big.Int, call ethereum.CallMsg) (uint64, error) {
	if chainID.Text(10) == "0" {
		return 0, fmt.Errorf("Could not estimate gas")
	}
	return 18, nil
}

func makeGasEstimatorContext(i int) *worker.Context {
	ctx := worker.NewContext()
	ctx.Reset()
	ctx.Logger = log.NewEntry(log.StandardLogger())
	ctx.T.Sender = &common.Account{}
	ctx.T.Tx = &ethpb.Transaction{TxData: &ethpb.TxData{}}

	switch i % 2 {
	case 0:
		ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(0))
		ctx.Keys["errors"] = 1
		ctx.Keys["result"] = uint64(0)
	case 1:
		ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(1))
		ctx.Keys["errors"] = 0
		ctx.Keys["result"] = uint64(18)
	}
	return ctx
}

func TestGasEstimator(t *testing.T) {
	me := MockGasEstimator{t: t}
	estimator := GasEstimator(&me)

	rounds := 100
	outs := make(chan *worker.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeGasEstimatorContext(i)
		go func(ctx *worker.Context) {
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

		if out.T.Tx.TxData.GetGas() != result {
			t.Errorf("GasEstimator: expected gas limit %v but got %v", result, out.T.Tx.TxData.GetGas())
		}
	}
}

type MockGasPricer struct {
	t *testing.T
}

func (e *MockGasPricer) SuggestGasPrice(ctx context.Context, chainID *big.Int) (*big.Int, error) {
	if chainID.Text(10) == "0" {
		return big.NewInt(0), fmt.Errorf("Could not estimate gas")
	}
	return big.NewInt(10), nil
}

func makeGasPricerContext(i int) *worker.Context {
	ctx := worker.NewContext()
	ctx.Reset()
	ctx.Logger = log.NewEntry(log.StandardLogger())
	ctx.T.Tx = &ethpb.Transaction{TxData: &ethpb.TxData{}}

	switch i % 2 {
	case 0:
		ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(0))
		ctx.Keys["errors"] = 1
		ctx.Keys["result"] = ""
	case 1:
		ctx.T.Chain = (&common.Chain{}).SetID(big.NewInt(1))
		ctx.Keys["errors"] = 0
		ctx.Keys["result"] = "0xa"
	}
	return ctx
}

func TestGasPricer(t *testing.T) {
	mp := MockGasPricer{t: t}
	pricer := GasPricer(&mp)

	rounds := 100
	outs := make(chan *worker.Context, rounds)
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		ctx := makeGasPricerContext(i)
		go func(ctx *worker.Context) {
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
		errCount, result := out.Keys["errors"].(int), out.Keys["result"]
		if len(out.T.Errors) != errCount {
			t.Errorf("GasPricer: expected %v errors but got %v", errCount, out.T.Errors)
		}

		if out.T.Tx.TxData.GetGasPrice() != result {
			t.Errorf("GasPricer: expected gas price %v but got %v", result, out.T.Tx.TxData.GetGasPrice())
		}
	}
}
