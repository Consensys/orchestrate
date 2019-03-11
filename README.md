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
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/types"
)

// Define a handler method
func handler(ctx *types.Context) {
	ctx.Logger.Infof("Handling %v\n", ctx.Msg.(string))
}

func main() {
	// Instantiate worker (limited to 1 message processed at a time)
	worker := core.NewWorker(1)

	// Register handler
	worker.Use(handler)

	// Create an input channel of messages
	in := make(chan interface{})

	// Run worker on input channel
	go func() { worker.Run(in) }()

	// Feed channel
	in <- "Message-1"
	in <- "Message-2"
	in <- "Message-3"

	// Close channel & wiat for worker to treat all messages
	close(in)
	<-worker.Done()
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

Handler functions are the building blocks for workers, they match the interface

```go
type HandlerFunc func(ctx *Context)
```

When creating a worker you must register a sequence of handlers by using ``worker.Use(handler)``. When running, each time a new message is feeded to the worker, the worker generates a ``types.Context`` and apply handlers sequence on this context object.

#### Pipeline/Middleware

Handlers can be either 

- *pipeline* meaning it proceed to its execution then closes
- *middleware* meaning it proceeds to beginning of its own execution, then execute pending handlers then finishes own execution

Middleware is a common pattern that permit to maintain a scope of variables open while executing unknown functions.

```sh
$ cat examples/pipeline-middleware/main.go
```

```go
package main

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/types"
)

// Define a pipeline handler
func pipeline(ctx *types.Context) {
	ctx.Logger.Infof("Pipeline handling %v\n", ctx.Msg.(string))
}

// Define a middleware handler
func middleware(ctx *types.Context) {
	// Start middleware execution
	ctx.Logger.Infof("Middleware starts handling %v\n", ctx.Msg.(string))

	// Trigger execution of pending handlers
	ctx.Next()

	// Executed after pending handlers have executed
	ctx.Logger.Infof("Middleware finishes handling %v\n", ctx.Msg.(string))
}

func main() {
	// Instantiate worker (limited to 1 message processed at a time)
	worker := core.NewWorker(1)

	// Register handlers
	worker.Use(middleware)
	worker.Use(pipeline)

	// Create an input channel of messages
	in := make(chan interface{})

	// Run worker on input channel
	go func() { worker.Run(in) }()

	// Feed channel
	in <- "Message-1"
	in <- "Message-2"

	// Close channel & wiat for worker to treat all messages
	close(in)
	<-worker.Done()
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

A worker can handle multiple message at once in parallel goroutines, therefore handler functions must manage there own resource in concurrenct safe manner. Note that while multiple contexts can be handled in parallel, a given context is never handled by more than one handler function at a time.

```sh
$ cat examples/concurrency/main.go
```

```go
package main

import (
	"fmt"
	"sync/atomic"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/types"
)

// ExampleHandler is an handler that increment counters
type ExampleHandler struct {
	safeCounter   uint32
	unsafeCounter uint32
}

func (h *ExampleHandler) handleSafe(ctx *types.Context) {
	// Increment counter using atomic
	atomic.AddUint32(&h.safeCounter, 1)
}

func (h *ExampleHandler) handleUnsafe(ctx *types.Context) {
	// Increment counter with no concurrent protection
	h.unsafeCounter++
}

func main() {
	// Instantiate a worker that can treat 1000 messages in parallelW
	worker := core.NewWorker(1000)

	// Register handler
	h := ExampleHandler{0, 0}
	worker.Use(h.handleSafe)
	worker.Use(h.handleUnsafe)

	// Start worker
	in := make(chan interface{})
	go func() { worker.Run(in) }()

	// Feed 10000 to the worker
	for i := 0; i < 10000; i++ {
		in <- "Message"
	}

	// Close channel
	close(in)
	<-worker.Done()

	// Print counters
	fmt.Printf("* Safe counter: %v\n", h.safeCounter)
	fmt.Printf("* Unsafe counter: %v\n", h.unsafeCounter)
}
```

```sh
# Run example (note that unsafe counter output is not deterministic)
$ go run examples/concurrency/main.go

* Safe counter: 10000
* Unsafe counter: 9989
```

## Cobra CLI

### Quick Start

```sh
$ cat examples/command/main.go
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
	Use:              "worker",
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
  worker example [OPTIONS] [flags]

Flags:
      --eth-client strings   Ethereum client URLs.
                             Environment variable: "ETH_CLIENT_URL" (default [https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7,https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c,https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c,https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c])
  -h, --help                 help for example
      --log-format string    Log formatter (one of ["text" "json"]).
                             Environment variable: "LOG_FORMAT" (default "text")
      --log-level string     Log level (one of ["panic" "fatal" "error" "warn" "info" "debug" "trace"]).
                             Environment variable: "LOG_LEVEL" (default "debug")
```
