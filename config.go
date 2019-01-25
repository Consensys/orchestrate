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

// KafkaConfig is the configuration of application dealing with Kafka
type KafkaConfig struct {
	ConsumerGroup string `short:"c" long:"consumer-group" env:"KAFKA_CONSUMER_GROUP" default:"tx-nonce-group"`
	OutTopic      string `short:"o" long:"out-topic" env:"KAFKA_TOPIC_TX_SENDER" default:"topic-tx-sender"`
	InTopic       string `short:"i" long:"in-topic" env:"KAFKA_TOPIC_TX_NONCE" default:"topic-tx-nonce"`
	Address       string `long:"kafka-address" env:"KAFKA_ADDRESS" default:"localhost:9092"`
}

// RedisConfig is the configuration of application dealing with Redis
type RedisConfig struct {
	Address     string `long:"redis-address" env:"REDIS_ADDRESS" default:"localhost:6379"`
	LockTimeout int    `long:"redis-lock-timeout" env:"REDIS_LOCKTIMEOUT" default:"1500"`
}

// EthConfig is the configuration of application dealing with Ethereum
type EthConfig struct {
	URL string `short:"e" long:"eth-client" env:"ETH_CLIENT" default:"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7"`
}

// Config worker configuration
type Config struct {
	Log   LoggerConfig
	App   AppConfig
	Kafka KafkaConfig
	Redis RedisConfig
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
