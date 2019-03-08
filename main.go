package main

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/boilerplate-worker.git/cmd"
)

func main() {
	command := cmd.NewCommand()

	if err := command.Execute(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
