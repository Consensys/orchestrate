package infra

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	flags "github.com/jessevdk/go-flags"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
)

func MockCrediter(ctx context.Context, r *services.FaucetRequest) (*big.Int, bool, error) {
	return r.Value, true, nil
}

func MockBalanceAt(ctx context.Context, chainID *big.Int, a common.Address) (*big.Int, error) {
	return big.NewInt(100000000000000000), nil
}

func TestFaucet(t *testing.T) {
	// Set default configuration
	opts := FaucetConfig{}
	flags.ParseArgs(&opts, []string{})

	faucet, err := CreateFaucet(opts, MockBalanceAt, MockCrediter)
	if err != nil {
		t.Errorf("Faucet should be created from default option")
	}

	// Valid Credit
	req := services.FaucetRequest{
		ChainID: big.NewInt(1),
		Address: common.HexToAddress("0x7E654d251Da770A068413677967F6d3Ea2FeA9E4"),
		Value:   big.NewInt(50000000000000000),
	}
	amount, ok, _ := faucet.Credit(context.Background(), &req)
	if !ok || amount.Uint64() != 50000000000000000 {
		t.Errorf("Expected valid transfer but got %v %v", ok, amount)
	}

	// Invalid Credit
	req = services.FaucetRequest{
		ChainID: big.NewInt(3),
		Address: common.HexToAddress("0x7E654d251Da770A068413677967F6d3Ea2FeA9E4"),
		Value:   big.NewInt(50000000000000000),
	}
	amount, ok, _ = faucet.Credit(context.Background(), &req)
	if ok || amount.Uint64() != 0 {
		t.Errorf("Expected invalid transfer")
	}
}
