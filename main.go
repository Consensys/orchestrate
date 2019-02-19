package main

import (
	"fmt"
	"os"

	"gitlab.com/ConsenSys/client/fr/core-stack/boilerplate-worker.git/cmd"
)

func main() {
	command := cmd.NewCommand()

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
