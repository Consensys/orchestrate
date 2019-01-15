package main

import (
	flags "github.com/jessevdk/go-flags"

	log "github.com/sirupsen/logrus"
)

var opts struct {
	Log struct {
		Level  string `short:"l" long:"log-level" env:"LOG_LEVEL" default:"debug" description:"Log level, one of panic, fatal, error, warn, info, debug, trace."`
		Format string `long:"log-format" env:"LOG_FORMAT" default:"text" description:"Log formatter, one of text, json."`
	}
}

func init() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	configureLogger(opts.Log.Level, opts.Log.Format)
}

func configureLogger(level string, format string) {
	switch format {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}
	if logLevel, err := log.ParseLevel(level); err != nil {
		log.Fatalf("Invalid log level, %v", err)
	} else {
		log.New()
		log.SetLevel(logLevel)
	}
	log.Debugf("%+v", opts)
}
