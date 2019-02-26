package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/cmd"
)

func main() {
	command := cmd.NewCommand()

	if err := command.Execute(); err != nil {
		log.Errorf("%v\n", err)
		os.Exit(1)
	}
}
