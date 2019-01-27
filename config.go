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
	Vault struct {
		Accounts []string `short:"a" long:"vault-accounts" env:"VAULT_ACCOUNTS" default:"56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E" default:"5FBB50BFF6DFAD35C4A374C9237BA2F7EAED9C6868E0108CB259B62D68029B1A" default:"86B021CCB810F26A30445B85F71E4C1596A11A97DDF9B9E348AC93D1DA6735BC" default:"DD614C3B343E1B6DBD1B2811D4F146CC90337DEEF96AB97C353578E871B19D5E" default:"425D92F63A836F890F1690B34B6A25C2971EF8D035CD8EA8592FD1069BD151C6" default:"C4B172E72033581BC41C36FA0448FCF031E9A31C4A3E300E541802DFB7248307" default:"706CC0876DA4D52B6DCE6F5A0FF210AEFCD51DE9F9CFE7D1BF7B385C82A06B8C" default:"1476C66DE79A57E8AB4CADCECCBE858C99E5EDF3BFFEA5404B15322B5421E18C" default:"A2426FE76ECA2AA7852B95A2CE9CC5CC2BC6C05BB98FDA267F2849A7130CF50D" default:"41B9C5E497CFE6A1C641EFCA314FF84D22036D1480AF5EC54558A5EDD2FEAC03"`
	}
	WorkerSlots uint `short:"w" long:"worker-slots" env:"WORKER" default:"50"`
}

// KafkaConfig is the config part concerning kafka
type KafkaConfig struct {
	ConsumerGroup string `short:"c" long:"consumer-group" env:"CONSUMER_GROUP" default:"tx-signer-group"`
	InTopic       string `short:"i" long:"in-topic" env:"KAFKA_TOPIC_TX_SIGNER" default:"topic-tx-signer"`
	OutTopic      string `short:"o" long:"out-topic" env:"KAFKA_TOPIC_TX_SENDER" default:"topic-tx-sender"`
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
