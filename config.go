package main

import (
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
	WorkerSlots uint `short:"w" long:"worker-slots" env:"WORKER" default:"50"`
}

// KafkaConfig is the config part concerning kafka
type KafkaConfig struct {
	ConsumerGroup string `short:"c" long:"consumer-group" env:"CONSUMER_GROUP" default:"tx-decoder-group"`
	InTopic       string `short:"i" long:"in-topic" env:"KAFKA_TOPIC_TX_DECODER" default:"topic-tx-decoder"`
	OutTopic      string `short:"o" long:"out-topic" env:"KAFKA_TOPIC_TX_DECODED" default:"topic-tx-decoded"`
	Address       string `long:"kafka-address" env:"KAFKA_ADDRESS" default:"localhost:9092"`
}

// EthConfig is the config part concerning the ethereum environment
type EthConfig struct {
	URL string `short:"e" long:"eth-client" env:"ETH_CLIENT_URL" default:"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7"`
}

// Config worker configuration
type Config struct {
	Log   LoggerConfig
	App   AppConfig
	Kafka KafkaConfig
	Eth   EthConfig
}

// LoadConfig load configuration
func LoadConfig(opts interface{}) {
	_, err := flags.Parse(opts)
	if err != nil {
		panic(err)
	}
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
