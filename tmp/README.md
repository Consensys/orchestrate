# Core

Core is the lower level package in Core-Stack, it implements building blocks that are shared in all core stack *infra* and *microservice*. In particular it implements

- *types* which define all main types used in core-stack (such as ``types.Context`` which is the central object manipulated by all Core-Stack workers)
- *services* which define *interfaces* for microservices to speak together
- *protobuf* schemes 

## Installation

To install Core-Stack Core package, you need to install Go and set your Go workspace first.

1. Download and install it:

```sh
$ go get -u gitlab.com/ConsenSys/client/fr/core-stack/core.git
```

2. Import it in your code:

```go
import "gitlab.com/ConsenSys/client/fr/core-stack/core.git"
```

## Prerequisite

Core-Stack requires Go 1.11

## Worker

### Quick Start

```sh
$ cat examples/quick-start/main.go
```

```go
package main

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
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
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
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

	"gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
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
