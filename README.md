# Nonce

Nonce implement infrastructure elements to smoothly build a Nonce.

## Installation

To install Core-Stack Core package, you need to install Go and set your Go workspace first.

1. Download and install it:

```sh
$ go get -u "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"

```

2. Import it in your code:

```go
import "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
```

## Prerequisite

Core-Stack requires Go 1.12

## Create a Nonce

```sh
$ cat examples/simple/main.go
```

```go
package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/nonce.git"
)

// TODO
```

```sh
# Run example
$ go run examples/simple/main.go

```
