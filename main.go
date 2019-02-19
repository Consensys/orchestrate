package main

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	var opts Config
	LoadConfig(&opts)
	ConfigureLogger(opts.Log)
	go http.ListenAndServe(opts.HTTP.Hostname, prepareHTTPRouter(context.Background()))
	log.Info("Start worker...")
}
