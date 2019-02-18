# Common

Implements common elements used through core stack such as handlers and infrastructures.

## Installation

To install Core-Stack Core package, you need to install Go and set your Go workspace first.

1. Download and install it:

```sh
$ go get -u gitlab.com/ConsenSys/client/fr/core-stack/core.git
```

2. Import it in your code:

```go
import "gitlab.com/ConsenSys/client/fr/core-stack/common.git"
```

## Prerequisite

Core-Stack requires Go 1.11

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
	"gitlab.com/ConsenSys/client/fr/core-stack/common.git/config"
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
