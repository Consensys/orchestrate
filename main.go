package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	var opts Config
	LoadConfig(&opts)
	ConfigureLogger(opts.Log)
	log.Info("Start worker...")
}
