package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	command := NewRunCommand()
	if err := command.Execute(); err != nil {
		log.WithError(err).Fatalf("main: execution failed")
	}
	log.Infof("main: execution completed")
}
