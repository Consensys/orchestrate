# Ethereum

Implement elements of infrastructure based on go-ethereum such as

- *Client* to connect to Ethereum clients
- *Tx Encoding*/*Decoding* resources
- *Tx-Listener* that is a blockchain tx-listener connected to multiple chains
- *Signer* to sign transaction

## Installation

To install Core-Stack Core package, you need to install Go and set your Go workspace first.

1. Download and install it:

```sh
$ go get -u gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core
```

2. Import it in your code:

```go
import "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git"
```

## Prerequisite

Core-Stack requires Go 1.11

## Multi-Chain Tx-Listener

### Quick Start

```sh
$ cat examples/tx-listener/main.go
```

```go
package main

import (
	"context"
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener"
)

func main() {
	ethURLs := []string{
		"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7",
		"https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
	}
	mec, err := ethclient.MultiDial(ethURLs)
	if err != nil {
		fmt.Printf("Error Dialing client: %v", err)
		return
	}

	// Create listener
	listenerCfg := listener.NewConfig()
	txlistener := listener.NewTxListener(listener.NewEthClient(mec, listenerCfg))

	// Start listening on every chain starting from last block
	for _, chainID := range mec.Networks(context.Background()) {
		txlistener.Listen(chainID, -1, 0, listenerCfg)
	}

	// Consume receipts
	for r := range txlistener.Receipts() {
		fmt.Printf("* New receipt ChainID=%v BlockNumber=%v TxHash=%v\n", r.ChainID.Text(16), r.BlockNumber, r.TxHash.Hex())
	}
}

```

```sh
# Run example
$ go run examples/quick-start/main.go

INFO[0002] tx-listener: start listening from block=5018648 tx=0  Chain=3
INFO[0002] tx-listener: start listening from block=3869057 tx=0  Chain=4
INFO[0002] tx-listener: start listening from block=10361214 tx=0  Chain=2a
INFO[0002] tx-listener: start listening from block=7220203 tx=0  Chain=1
* New receipt ChainID=4 BlockNumber=3869057 TxHash=0xeccc4b6d97b98030d6c0a829c18307b3875bd9ece17a514e733b432b2f37ebdb
* New receipt ChainID=4 BlockNumber=3869057 TxHash=0x2becf34f2dbeb2024b9692bc8fd0b97778b2e380e32ee3a8a7757f68139644b3
* New receipt ChainID=4 BlockNumber=3869057 TxHash=0x93d089e31f60957ebe9f9697362812e970265df77954535d25f5c2c049da91c8
```

