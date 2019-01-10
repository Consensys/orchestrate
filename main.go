package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/infra"
)

func newEthClient(rawurl string) *infra.EthClient {
	ec, err := infra.Dial(rawurl)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to Ethereum client")
	return ec
}

func main() {
	var chainURL = "https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7"

	ec := newEthClient(chainURL)
	block, err := ec.BlockByNumber(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(block.Transactions())

	b, err := rlp.EncodeToBytes(&block)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(b))
}
