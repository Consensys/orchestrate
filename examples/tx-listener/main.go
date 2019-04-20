package main

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	listener "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener"
)

func main() {
	// Set Ethereum client URLs to connect to
	viper.Set("eth.clients", []string{
		"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7",
		"https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
	})

	// Initialize Multi-Client
	ethclient.Init(context.Background())

	// Create listener
	config := listener.NewConfig()
	txlistener := listener.NewTxListener(ethclient.GlobalMultiClient(), config)

	// Start listening on every chain starting from last block
	for _, chainID := range ethclient.GlobalMultiClient().Networks(context.Background()) {
		txlistener.Listen(chainID, -1, 0)
	}

	// Consume receipts
	for r := range txlistener.Receipts() {
		fmt.Printf("* New receipt ChainID=%v BlockNumber=%v TxHash=%v\n", r.ChainID.Text(16), r.BlockNumber, r.TxHash.Hex())
	}
}
