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
	// Instantiate Engine
	cfg := engine.NewConfig()
	engine := engine.NewEngine(&cfg)

	// Register an handler
	engine.Register(handler)

	// Create an input channel of messages
	in := make(chan interface{})

	// Run Engine on input channel
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

	// Close channel & wait for Engine to treat all messages
	close(in)
	wg.Wait()

	// CleanUp Engine to avoid memory leak
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

When creating an engine you must register a chain of handlers by using ``engine.Register(handler)``. When running, each time a new message is feeded to the engine, the engine generates a ``engine.TxContext`` and apply handlers sequence on this context object.

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

func aborter(txctx *engine.TxContext) {
	txctx.Logger.Infof("Aborting %v\n", txctx.Msg.(string))
	txctx.AbortWithError()
}

// Define a pipeline handler
func pipeline(txctx *engine.TxContext) {
	txctx.Logger.Infof("Pipeline handling %v\n", txctx.Msg.(string))
}

// Define a middleware handler
func middleware(txctx *engine.TxContext) {
	// Start middleware execution
	txctx.Logger.Infof("Middleware starts handling %v\n", txctx.Msg.(string))

	// Trigger execution of pending handlers
	txctx.Next()

	// Executed after pending handlers have executed
	txctx.Logger.Infof("Middleware finishes handling %v\n", txctx.Msg.(string))
}

func main() {
	// Register handlers
	engine.Init(context.Background())
	engine.Register(middleware)
	engine.Register(pipeline)
	engine.Register(aborter)

	// Create an input channel of messages
	in := make(chan interface{})

	// Run Engine on input channel
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		engine.Run(context.Background(), in)
		wg.Done()
	}()

	// Feed channel
	in <- "Message-1"
	in <- "Message-2"

	// Close channel & wait for Engine to treat all messages
	close(in)
	wg.Wait()

	// CleanUp Engine to avoid memory leak
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

### Composite handlers

It is possible to compose handlers together in order to create more complex handlers. 

`CombineHandlers(handlers ...HandlerFunc)` function is here for this purpose.

```sh
$ cat examples/composite-handlers/main.go
```

```go
package main

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/examples"
)

func aborter(txctx *engine.TxContext) {
	txctx.Logger.Infof("Aborting %v\n", txctx.Msg.(examples.Msg))
	txctx.Abort()
}

// Define a pipeline handler
func pipeline(name string) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger.Infof("Pipeline-%v handling %v\n", name, txctx.Msg.(examples.Msg))
	}
}

// Define a middleware handler
func middleware(name string) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Start middleware execution
		txctx.Logger.Infof("Middleware-%v starts handling %v\n", name, txctx.Msg.(examples.Msg))

		// Trigger execution of pending handlers
		txctx.Next()

		// Executed after pending handlers have executed
		txctx.Logger.Infof("Middleware-%v finishes handling %v\n", name, txctx.Msg.(examples.Msg))
	}
}

func main() {
	// Register handlers
	engine.Init(context.Background())

	// Declare individual handlers
	pipepline1 := pipeline("1")
	pipepline2 := pipeline("2")
	pipepline3 := pipeline("3")
	pipepline4 := pipeline("4")
	middleware1 := middleware("1")
	middleware2 := middleware("2")
	middleware3 := middleware("3")

	// Declare 2 composite handlers
	left := engine.CombineHandlers(middleware2, pipepline2)
	right := engine.CombineHandlers(middleware3, aborter, pipepline3)

	// Declare a forked handler
	fork := func(txctx *engine.TxContext) {
		switch txctx.Msg.Entrypoint() {
		case "left":
			left(txctx)
		case "right":
			right(txctx)
		}
	}

	// Declare overall composite handler
	handler := engine.CombineHandlers(pipepline1, middleware1, fork, pipepline4)

	// Register composite handler
	engine.Register(handler)

	// Create an input channel of messages
	in := make(chan engine.Msg)

	// Run myEngine on input channel
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		engine.Run(context.Background(), in)
		wg.Done()
	}()

	// Feed channel
	in <- examples.Msg("left")
	in <- examples.Msg("right")

	// Close channel & wait for myEngine to treat all messages
	close(in)
	wg.Wait()

	// CleanUp Engine to avoid memory leak
	engine.CleanUp()
}
```

```sh
# Run example (note that handlers sequence is not the same when applying left message or right message)
$ go run examples/composite-handlers/main.go

INFO[0000] Pipeline-1 handling left
INFO[0000] Middleware-1 starts handling left
INFO[0000] Middleware-2 starts handling left
INFO[0000] Pipeline-2 handling left
INFO[0000] Middleware-2 finishes handling left
INFO[0000] Pipeline-4 handling left
INFO[0000] Middleware-1 finishes handling left
INFO[0000] Pipeline-1 handling right
INFO[0000] Middleware-1 starts handling right
INFO[0000] Middleware-3 starts handling right
INFO[0000] Aborting right
INFO[0000] Middleware-3 finishes handling right
INFO[0000] Middleware-1 finishes handling right
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
	// Instantiate Engine that can treat 100 message concurrently
	// Instantiate an Engine that can treat 100 message concurrently in 100 distinct partitions
	cfg := engine.NewConfig()
	cfg.Slots = 100
	engine := engine.NewEngine(&cfg)

	// Register handler
	h := ExampleHandler{0, 0}
	engine.Register(h.handleSafe)
	engine.Register(h.handleUnsafe)

	// Run Engine on 100 distinct input channel
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

	// Feed 10000 to the Engine
	for i := 0; i < 100; i++ {
		for j, in := range inputs {
			in <- fmt.Sprintf("Message %v-%v", j, i)
		}
	}

	// Close all channels & wait for Engine to treat all messages
	for _, in := range inputs {
		close(in)
	}
	wg.Wait()

	// CleanUp Engine to avoid memory leak
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
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/logger"
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
	},
}

func init() {
	rootCmd.AddCommand(cmdExample)
	logger.LogLevel(cmdExample.Flags())
	logger.LogFormat(cmdExample.Flags())
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

## Errors

Internal errors are built on top of protobuf and contain

- `string` **message** describing the error
- `uint64` **error code** that maps to a five hex character code (e.g. `F000A`)
- `string` **component** indicating in system the error has been raised
- `map<string, string>` **extra** holding extra information to diagnose the error

### Error codes

Error codes are `uint64` that maps to five hex character code

| Class | Subclass | Error Code | Condition                     | Comment                                   |
|-------|----------|------------|-------------------------------|-------------------------------------------|
| 01XXX |          |    01000   | warning                       | Raised to indicate a warning              |
| 01XXX |   011XX  |    01100   | retry_warning                 | Error occured system retries              |
| 01XXX |   012XX  |    01200   | faucet_warning                | Faucet credit has been denied             |
| 08XXX |          |    08000   | connection_exception          | Failed connecting to an external service  |
| 08XXX |   081XX  |    08100   | kafka_connection_exception    | Failed connecting to Kafka                |
| 08XXX |   082XX  |    08200   | http_connection_exception     | Failed connecting to an HTTP service      |
| 08XXX |   083XX  |    08300   | ethereum_connection_exception | Failed connecting to Ethereum jsonRPC API |
| 0AXXX |          |    0A000   | feature_not_supported         | Feature is not supported                  |
| 42XXX |          |    42000   | invalid_data                  | Failed to process data                    |
| 42XXX |   421XX  |    42100   | invalid_encoding              | Failed to decode a message                |
| 42XXX |   422XX  |    42200   | invalid_solidity_data         | Failed to process Solidity related data   |
| 42XXX |   422XX  |    42201   | invalid_method_signature      | Invalid Solidity method signature         |
| 42XXX |   422XX  |    42202   | invalid_args_count            | Invalid args count provided               |
| 42XXX |   422XX  |    42203   | invalid_arg                   | Invalid arg provided                      |
| 42XXX |   422XX  |    42204   | invalid_topics_count          | Invalid topics count in event log         |
| 42XXX |   422XX  |    42205   | invalid_event_data            | Invalid data in event log                 |
| 42XXX |   423XX  |    42300   | invalid_format                | Data does not match expected format       |
| DBXXX |          |    DB000   | storage_exception             | Failed accessing stored data              |
| DBXXX |   DB1XX  |    DB100   | constraint_violated           | Data constraint violated                  |
| DBXXX |   DB2XX  |    DB200   | not_found                     | No data found for given parameters        |
| DBXXX |   DB3XX  |    DB300   | data_corrupted                | Data is corrupted                         |
| F0XXX |          |    F0000   | invalid_config                | Invalid configuration                     |
| FFXXX |          |    FF000   | internal_error                | Internal error                            |