# Ethereum

Implement elements of infrastructure based on go-ethereum such as

- *Client* to connect to Ethereum clients
- *Tx Encoding*/*Decoding* resources
- *Tx-Listener*
- *Signer* to sign transaction

## Installation

To install Core-Stack Core package, you need to install Go and set your Go workspace first.

1. Download and install it:

```sh
$ go get -u gitlab.com/ConsenSys/client/fr/core-stack/core.git
```

2. Import it in your code:

```go
import "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git"
```

## Prerequisite

Core-Stack requires Go 1.11