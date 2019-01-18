# Core

Core is the lower level package in Core-Stack, it implements building blocks that are shared in all core stack *infra* and *microservice*. In particular it implements


- *types* which define all main types used in core-stack (such as ``types.Context`` which is the central object manipulated by all Core-Stack workers)
- *services* which define *interfaces* for microservices to speak together
- *protobuf* schemes 

## Contents

- [Core](#core)
  - [Contents](#contents)
  - [Installation](#installation)
  - [Prerequisite](#prerequisite)
  - [Worker](#worker)
    - [Quick Start](#quick-start)
    - [Handlers](#handlers)
      - [Pipeline/Middleware](#pipelinemiddleware)

## Installation

To install Core-Stack COre package, you need to install Go and set your Go workspace first.

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



Context-Manager is a go package part of ConsenSys France Core-Stack. Context-manager allows to manipulate a context that is transmitted between microservices in Core-Stack.

It bases on proto-buffer protocol to serialize and deserialize message that can be transmitted from a microservice to another.

## Worker

### Quick Start

```sh
$ cat examples/quick-start/main.go
```

```go
package main

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Define a handler method 
func handler(ctx *types.Context) {
	fmt.Printf("* Handling %v\n", ctx.Msg.(string))
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

* Handling Message-1
* Handling Message-2
* Handling Message-3
```

### Handlers

Handler functions are the building blocks for workers, they match the interface

```go
type HandlerFunc func(ctx *Context)
```

When creating a worker you must register a sequence of handlers by using ``worker.Use(handler)``. When running, worker will apply handlers sequence on every message feeded to the worker.

#### Pipeline/Middleware

Handlers can be either 

- *pipeline* meaning it proceed to its execution then closes
- *middleware* meaning it proceeds to part of its execution, execution pending handlers then finishes execution

Middleware is a very common pattern that permit to maintain a scope open while executing unknown functions.

```sh
$ cat examples/pipeline-middleware/main.go
```

```go
package main

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

// Define a pipeline handler
func pipeline(ctx *types.Context) {
	fmt.Printf("* Pipeline handling %v\n", ctx.Msg.(string))
}

// Define a middleware handler
func middleware(ctx *types.Context) {
	// Start middleware execution
	fmt.Printf("* Middleware starts handling %v\n", ctx.Msg.(string))

	// Trigger execution of pending handlers
	ctx.Next()

	// Executed after pending handlers have executed
	fmt.Printf("* Middleware finishes handling %v\n", ctx.Msg.(string))
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

* Middleware starts handling Message-1
* Pipeline handling Message-1
* Middleware finishes handling Message-1
* Middleware starts handling Message-2
* Pipeline handling Message-2
* Middleware finishes handling Message-2
```