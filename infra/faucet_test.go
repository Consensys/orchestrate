package infra

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
)

func MockCrediter(ctx context.Context, r *services.FaucetRequest) (*big.Int, bool, error) {
	return r.Value, true, nil
}

func MockBalanceAt(ctx context.Context, chainID *big.Int, a common.Address) (*big.Int, error) {
	return big.NewInt(100000000000000000), nil
}

func TestFaucet(t *testing.T) {
	// Set configuration for test
	viper.Set("faucet.blacklist", []string{"0x7E654d251Da770A068413677967F6d3Ea2FeA9E4@3"})
	viper.Set("faucet.addresses", []string{"0x7E654d251Da770A068413677967F6d3Ea2FeA9E4@3"})
	viper.Set("faucet.cooldown", 60*time.Second)
	viper.Set("faucet.max", "200000000000000000")

	faucet, err := CreateFaucet(MockBalanceAt, MockCrediter)
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
