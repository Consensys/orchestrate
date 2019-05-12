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
import "gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git"
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
	"math/big"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	handler "gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tx-listener/handler/base"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tx-listener/listener"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Handler is a engine HandlerFunc
func Handler(txctx *engine.TxContext) {
	// Cast message into sarama.ConsumerMessage
	r, ok := txctx.Msg.(*types.TxListenerReceipt)
	if !ok {
		panic("loader: expected a types.TxListenerReceipt")
	}

	fmt.Printf("* New receipt ChainID=%v BlockNumber=%v Txindex=%v TxHash=%v\n", r.ChainID.Text(10), r.BlockNumber, r.TxIndex, r.TxHash.Hex())
}

func main() {
	// Set log configuration
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)

	// Initialize Listener
	viper.Set("eth.clients", []string{
		"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7",
		"https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
	})

	// Initialize listener
	listener.Init(context.Background())

	// Initialize engine and register handlers
	engine.Init(context.Background())
	engine.Register(Handler)

	// Create handler
	conf, err := handler.NewConfig()
	if err != nil {
		log.WithError(err).Fatalf("listener: could not create config")
	}
	h := handler.NewHandler(engine.GlobalEngine(), conf)

	// Start listening
	_ = listener.Listen(
		context.Background(),
		[]*big.Int{
			big.NewInt(1),
			big.NewInt(3),
		},
		h,
	)
}
```

```sh
# Run example
$ go run examples/quick-start/main.go

INFO[0001] ethereum: multi-client ready                  chains="[3 1]"
INFO[0001] tx-listener: ready to listen to chains        chains="[1 3]"
INFO[0002] tx-listener: start listening                  chain.id=3 start.block=5562232 start.tx-index=0
DEBU[0002] engine: start running loop                    loops.count=1
DEBU[0002] listener: start loop                          chain.id=3 start.block=5562232 start.tx-index=0
INFO[0002] tx-listener: start listening                  chain.id=1 start.block=7725464 start.tx-index=0
DEBU[0002] listener: start loop                          chain=1 start.block=7725464 start.tx-index=0
DEBU[0002] engine: start running loop                    loops.count=2
* New receipt ChainID=3 BlockNumber=5562232 Txindex=0 TxHash=0x8a2747cabca2917d0f8f5a18b53cc75091abd5995aaa3d9b4b6d5c6c82438d2c
* New receipt ChainID=3 BlockNumber=5562232 Txindex=1 TxHash=0x7a3393f49d6f51ec9117eb62624ed9ceda29b5879940627d874f0462fa2c6d05
* New receipt ChainID=3 BlockNumber=5562232 Txindex=2 TxHash=0xfbc889ad20a13299cd665b72daa68ec123f9bb2a940f4fb37a86a2bc35678207
* New receipt ChainID=3 BlockNumber=5562232 Txindex=3 TxHash=0x0bdc18d7cc0447f8d17a498408a90f4fdcd0a66ec57a6efc74bed1355af0ee33
* New receipt ChainID=3 BlockNumber=5562232 Txindex=4 TxHash=0x91188e789bdbfa81047681e168c475d91b5baa6abdf4679dbe9ee26c2d5a3255
* New receipt ChainID=3 BlockNumber=5562232 Txindex=5 TxHash=0x65fcae4e6344ad2fd3ab1c3f181d7c6bdb3c4ca7770811a3c2e5a88d81a57031
```
