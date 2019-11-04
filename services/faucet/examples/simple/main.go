package main

// import (
// 	"context"
// 	"fmt"
// 	"math/big"
// 	"time"

// 	"github.com/ethereum/go-ethereum/common"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/core/services"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet"
// )

// func mockCredit(ctx context.Context, r *services.FaucetRequest) (*big.Int, bool, error) {
// 	return r.Value, true, nil
// }

// func createFaucet() *faucet.ControlledFaucet {
// 	// Set BlackList controller
// 	chains := []*big.Int{
// 		big.NewInt(3), //
// 	}
// 	addresses := []common.Address{
// 		common.HexToAddress("0x7E654d251Da770A068413677967F6d3Ea2FeA9E4"),
// 	}
// 	bl := faucet.NewBlackList(chains, addresses)

// 	// Set Cooldown controller that requires a 60 sec interval between credits
// 	cd := faucet.NewCoolDown(time.Duration(60*time.Second), 50)

// 	// Set MaxBalance controller (it connects to ropsten to retrieve balance)
// 	chainURL := "https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7"
// 	ec, err := ethereum.Dial(chainURL)
// 	if err != nil {
// 		panic(err)
// 	}
// 	balanceAt := func(ctx context.Context, chainID *big.Int, a common.Address) (*big.Int, error) {
// 		return ec.BalanceAt(ctx, a, nil)
// 	}
// 	mb := faucet.NewMaxBalance(
// 		big.NewInt(200000000000000000), // MaxBalance authorized 0.2 ETH
// 		balanceAt,
// 	)

// 	// Create Faucet
// 	return faucet.NewControlledFaucet(mockCredit, bl.Control, cd.Control, mb.Control)
// }

// func main() {
// 	// Create faucet
// 	f := createFaucet()

// 	// Credit a random ethereum address with a value of ETH over MaxBalance
// 	amount, ok, _ := f.Credit(
// 		context.Background(),
// 		&services.FaucetRequest{
// 			ChainID: big.NewInt(3),
// 			Address: common.HexToAddress("0xd048EB6e9B7031f4fcfE264736A26b2A2268154B"),
// 			Value: big.NewInt(300000000000000000), // 0.3 ETH
// 		},
// 	)
// 	fmt.Printf("* 1. Amount credited=%v (credited=%v)", amount, ok)

// 	// Credit a random ethereum address
// 	amount, ok, _  = f.Credit(
// 		context.Background(),
// 		&services.FaucetRequest{
// 			ChainID: big.NewInt(3),
// 			Address: common.HexToAddress("0xd048EB6e9B7031f4fcfE264736A26b2A2268154B"),
// 			Value: big.NewInt(100000000000000000), // 0.1 ETH
// 		},
// 	)
// 	fmt.Printf("* 2. Amount credited=%v (credited=%v)", amount, ok)

// 	// Credit address again (expected to fail due to CoolDown)
// 	amount, ok, _ = f.Credit(
// 		context.Background(),
// 		&services.FaucetRequest{
// 			ChainID: big.NewInt(3),
// 			Address: common.HexToAddress("0xd048EB6e9B7031f4fcfE264736A26b2A2268154B"),
// 			Value: big.NewInt(100000000000000000), // 0.1 ETH
// 		},
// 	)
// 	fmt.Printf("* 3. Amount credited=%v (credited=%v)n", amount, ok)

// 	// Credit black list address again (expected to fail due to BLackList)
// 	amount, ok, _ =	f.Credit(
// 		context.Background(),
// 		&services.FaucetRequest{
// 			ChainID: big.NewInt(3),
// 			Address: common.HexToAddress("0x7E654d251Da770A068413677967F6d3Ea2FeA9E4"),
// 			Value: big.NewInt(100000000000000000), // 0.1 ETH
// 		},
// 	)
// 	fmt.Printf("* 4. Amount credited=%v (credited=%v)", amount, ok)
// }
