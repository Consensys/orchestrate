# Package

Core-Stack is a blockchain **Transaction Orchestration** system that can operate multiple chains simultaneously.
It provides production grade and agnostic mechanisms for transaction crafting, nonce management, transaction signing, transaction receipt listening, transaction receipt decoding, faucet and more.

Core-Stack is a **microservices architecture** composed of APIs & Workers. 
Workers communicate following **publish-subscribe** pattern using *Apache Kafka* as message broker. 
All messages are standardized using **protobuf** and **GRPC**.

Core-Stack is Plug & Play, a user only needs to send a business protobuf message to Core-Stack input topic,
Core-Stack then manages the full lifecycle of the transaction from crafting the transaction to decoding event logs data.

## Goal

Package is a low level library in Core-Stack dependency tree. In particular it implements

- *protobuf* containing all protobuf schemes
- *core* that defines core structural elements of Core-Stack (such as ``types.Context``) 
- *common* which are resources that are shared between multiple Core-Stack services

## Installation

To install Core-Stack Core package, you need to install Go and set your Go workspace first.

1. Download and install it:

```sh
$ go get -u gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core
```

2. Import it in your code:

```go
import "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git"
```

## Prerequisite

Core-Stack requires Go 1.12 or upper

## Worker

### Quick Start

```sh
$ cat examples/quick-start/main.go
```

```go
package main

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Define a handler method
func handler(ctx *engine.TxContext) {
	ctx.Logger.Infof("Handling %v\n", ctx.Msg.(string))
}

func main() {
	// Instantiate worker
	cfg := engine.NewConfig()
	engine := engine.NewEngine(&cfg)

	// Register an handler
	engine.Use(handler)

	// Create an input channel of messages
	in := make(chan interface{})

	// Run worker on input channel
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		engine.Run(context.Background(), in)
		wg.Done()
	}()

	// Feed channel
	in <- "Message-1"
	in <- "Message-2"
	in <- "Message-3"

	// Close channel & wait for worker to treat all messages
	close(in)
	wg.Wait()

	// CleanUp worker to avoid memory leak
	engine.CleanUp()
}
```

```sh
# Run example
$ go run examples/quick-start/main.go

INFO[0000] Handling Message-1
INFO[0000] Handling Message-2
INFO[0000] Handling Message-3
```

### Handlers

Handler functions are the building blocks for engines, they match the interface

```go
type HandlerFunc func(txctx *engine.TxContext)
```

When creating an engine you must register a chain of handlers by using ``engine.Use(handler)``. When running, each time a new message is feeded to the engine, the engine generates a ``engine.TxContext`` and apply handlers sequence on this context object.

#### Pipeline/Middleware

Handlers can be either 

- *pipeline* meaning handlers execution is linear
- *middleware* which has a 3 phases execution, it starts with its own execution then executes pending handlers and finally finish its own execution 

Middleware is a common pattern that allows to maintain a scope of variables open while executing unknown functions.

```sh
$ cat examples/pipeline-middleware/main.go
```

```go
package main

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Define a pipeline handler
func pipeline(ctx *engine.TxContext) {
	ctx.Logger.Infof("Pipeline handling %v\n", ctx.Msg.(string))
}

// Define a middleware handler
func middleware(ctx *engine.TxContext) {
	// Start middleware execution
	ctx.Logger.Infof("Middleware starts handling %v\n", ctx.Msg.(string))

	// Trigger execution of pending handlers
	ctx.Next()

	// Executed after pending handlers have executed
	ctx.Logger.Infof("Middleware finishes handling %v\n", ctx.Msg.(string))
}

func main() {
	cfg := engine.NewConfig()
	engine := engine.NewEngine(&cfg)

	// Register handlers
	engine.Use(middleware)
	engine.Use(pipeline)

	// Create an input channel of messages
	in := make(chan interface{})

	// Run worker on input channel
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		engine.Run(context.Background(), in)
		wg.Done()
	}()

	// Feed channel
	in <- "Message-1"
	in <- "Message-2"

	// Close channel & wait for worker to treat all messages
	close(in)
	wg.Wait()

	// CleanUp worker to avoid memory leak
	engine.CleanUp()
}
```

```sh
# Run example
$ go run examples/pipeline-middleware/main.go

INFO[0000] * Middleware starts handling Message-1
INFO[0000] * Pipeline handling Message-1
INFO[0000] * Middleware finishes handling Message-1
INFO[0000] * Middleware starts handling Message-2
INFO[0000] * Pipeline handling Message-2
INFO[0000] * Middleware finishes handling Message-2
```

#### Concurrency

A engine can handle multiple message concurrently in parallel goroutines. When declaring an handler you can configure

- ```slots``` which is the maximum count of messages that can be treated concurrently

**WARNING** ```HandlerFunc``` **MUST** be concurrency safe. 

Note that while multiple contexts can be handled in parallel, a given context is never handled by more than one ```HandlerFunc``` at a time.

```sh
$ cat examples/concurrency/main.go
```

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// ExampleHandler is an handler that increment counters
type ExampleHandler struct {
	safeCounter   uint32
	unsafeCounter uint32
}

func (h *ExampleHandler) handleSafe(ctx *engine.TxContext) {
	// Increment counter using atomic
	atomic.AddUint32(&h.safeCounter, 1)
}

func (h *ExampleHandler) handleUnsafe(ctx *engine.TxContext) {
	// Increment counter with no concurrent protection
	h.unsafeCounter++
}

func main() {
	// Instantiate worker that can treat 100 message concurrently
	// Instantiate an Engine that can treat 100 message concurrently in 100 distinct partitions
	cfg := engine.NewConfig()
	cfg.Slots = 100
	engine := engine.NewEngine(&cfg)

	// Register handler
	h := ExampleHandler{0, 0}
	engine.Use(h.handleSafe)
	engine.Use(h.handleUnsafe)

	// Run worker on 100 distinct input channel
	wg := &sync.WaitGroup{}
	inputs := make([]chan interface{}, 0)
	for i := 0; i < 100; i++ {
		inputs = append(inputs, make(chan interface{}, 100))
		wg.Add(1)
		go func(in chan interface{}) {
			engine.Run(context.Background(), in)
			wg.Done()
		}(inputs[i])
	}

	// Feed 10000 to the worker
	for i := 0; i < 100; i++ {
		for j, in := range inputs {
			in <- fmt.Sprintf("Message %v-%v", j, i)
		}
	}

	// Close all channels & wait for worker to treat all messages
	for _, in := range inputs {
		close(in)
	}
	wg.Wait()

	// CleanUp worker to avoid memory leak
	engine.CleanUp()

	// Print counters
	fmt.Printf("* Safe counter: %v\n", h.safeCounter)
	fmt.Printf("* Unsafe counter: %v\n", h.unsafeCounter)
}
```

```sh
# Run example (note that unsafe counter output is not deterministic)
$ go run examples/concurrency/main.go

* Safe counter: 10000
* Unsafe counter: 9994
```

## Cobra CLI

### Quick Start

```sh
$ cat examples/cobra/main.go
```

```go
package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common/config"
)

var rootCmd = &cobra.Command{
	Use:              "engine",
	TraverseChildren: true,
	Version:          "v0.1.0",
}

var cmdExample = &cobra.Command{
	Use:   "example [OPTIONS]",
	Short: "An example command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Log-Level:", viper.GetString("log.level"))
		fmt.Println("Log-Format:", viper.GetString("log.format"))
		fmt.Println("Eth-Clients:", viper.GetStringSlice("eth.clients"))
	},
}

func init() {
	rootCmd.AddCommand(cmdExample)
	config.LogLevel(cmdExample.Flags())
	config.LogFormat(cmdExample.Flags())
	config.EthClientURLs(cmdExample.Flags())
}

func main() {
	rootCmd.Execute()
}
```

```sh
# Run example
$ ETH_CLIENT_URL="http://localhost:8545 http://localhost:7545" go run examples/command/main.go  example --log-level fatal

Log-Level: fatal
Log-Format: text
Eth-Clients: [http://localhost:8545 http://localhost:7545]
```

```sh
# Run help command
$ go run examples/command/main.go help example

An example command

Usage:
  engine example [OPTIONS] [flags]

Flags:
      --eth-client strings   Ethereum client URLs.
                             Environment variable: "ETH_CLIENT_URL" (default [https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7,https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c,https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c,https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c])
  -h, --help                 help for example
      --log-format string    Log formatter (one of ["text" "json"]).
                             Environment variable: "LOG_FORMAT" (default "text")
      --log-level string     Log level (one of ["panic" "fatal" "error" "warn" "info" "debug" "trace"]).
                             Environment variable: "LOG_LEVEL" (default "debug")
```
