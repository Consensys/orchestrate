package main

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener/listener"
)

func main() {
	// Initialize client
	viper.Set("eth.clients", []string{
		"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7",
		"https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
	})
	ethclient.Init(context.Background())

	// Initialize listener
	listener.Init(context.Background())

	// Consume receipts
	for r := range listener.GlobalListener().Receipts() {
		fmt.Printf("* New receipt ChainID=%v BlockNumber=%v TxHash=%v\n", r.ChainID.Text(10), r.BlockNumber, r.TxHash.Hex())
	}
}
