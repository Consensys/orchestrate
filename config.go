package main

import (
	"net/url"

	flags "github.com/jessevdk/go-flags"

	log "github.com/sirupsen/logrus"
)

// LoggerConfig logger configuration
type LoggerConfig struct {
	Level  string `long:"log-level" env:"LOG_LEVEL" default:"debug" description:"Log level, one of panic, fatal, error, warn, info, debug, trace."`
	Format string `long:"log-format" env:"LOG_FORMAT" default:"text" description:"Log formatter, one of text, json."`
}

// AppConfig application configuration
type AppConfig struct {
	InTopic       string `short:"i" long:"in-topic" env:"IN_TOPIC" default:"tx-signer-in"`
	OutTopic      string `short:"o" long:"out-topic" env:"OUT_TOPIC" default:"tx-signer-out"`
	ConsumerGroup string `short:"c" long:"consumer-group" env:"CONSUMER_GROUP" default:"tx-signer-group"`
}

// ConnConfig connection configuration
type ConnConfig struct {
	Kafka struct {
		URL  string
		Host string `long:"kafka-host" env:"KAFKA_HOST" default:"localhost"`
		Port string `long:"kafka-port" env:"KAFKA_PORT" default:"9092"`
	}
	ETHClient struct {
		URL string `short:"e" long:"eth-client" env:"ETH_CLIENT" default:"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7"`
	}
}

// checkURL ensures that the given URL is well formed
func checkURL(name, u string) {
	_, err := url.ParseRequestURI(u)
	if err != nil {
		log.Error(name, " URL (", u, ") is invalid")
		panic(err)
	}

	log.Info(name, ": ", u)
}

// generateURL generates URL from host and port (<host>:<port>)
func generateURL(name, host, port string) string {
	u := host + ":" + port

	checkURL(name, u)

	return u
}

// configureURLs configure the URLs of the services that the worker connects to
func configureURLs() {
	opts.Conn.Kafka.URL = generateURL("Kafka", opts.Conn.Kafka.Host, opts.Conn.Kafka.Port)

	checkURL("ETH Client", opts.Conn.ETHClient.URL)
}

// Config worker configuration
type Config struct {
	Log  LoggerConfig
	App  AppConfig
	Conn ConnConfig
}

// LoadConfig load configuration
func LoadConfig(opts interface{}) {
	_, err := flags.Parse(opts)
	if err != nil {
		panic(err)
	}

	configureURLs()
}

// ConfigureLogger configure logger
func ConfigureLogger(opts LoggerConfig) {
	switch opts.Format {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}
	if logLevel, err := log.ParseLevel(opts.Level); err != nil {
		log.Fatalf("Invalid log level, %v", err)
	} else {
		log.New()
		log.SetLevel(logLevel)
	}
}
