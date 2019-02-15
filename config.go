package main

import (
	"fmt"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

// LoggerConfig logger configuration
type LoggerConfig struct {
	Level  string `long:"log-level" env:"LOG_LEVEL" default:"debug" description:"Log level, one of panic, fatal, error, warn, info, debug, trace."`
	Format string `long:"log-format" env:"LOG_FORMAT" default:"text" description:"Log formatter, one of text, json."`
}

// KafkaConfig is the configuration of application dealing with Kafka
type KafkaConfig struct {
	Address  []string `short:"k" long:"kafka-address" env:"KAFKA_ADDRESS" default:"localhost:9092" description:"Address of Kafka server to connect to"`
	OutTopic string   `short:"o" long:"kafka-out-topic" env:"KAFKA_TOPIC_TX_DECODER" default:"topic-tx-decoder" description:"Kafka topic to send message after processing"`
}

// WorkerConfig application configuration
type WorkerConfig struct {
	Slots uint `short:"w" long:"worker-slots" env:"WORKER_SLOTS" default:"100" description:"Number of messages that can be treat in parallel."`
}

// EthConfig is the configuration of application dealing with Ethereum
type EthConfig struct {
	URLs []string `short:"e" long:"eth-client" env:"ETH_CLIENT_URL" default:"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7" default:"https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c" default:"https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c" default:"https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c"`
}

// ListenerConfig is listener configuration
type ListenerConfig struct {
	Block struct {
		Backoff string `long:"listener-block-backoff" env:"LISTENER_BLOCK_BACKOFF" default:"1s" description:"Backoff time to wait before retrying after failing to find a mined block. Valid time units are ns, us (or µs), ms, s, m, h"`
		Limit   uint64 `long:"listener-block-limit" env:"LISTENER_BLOCK_LIMIT" default:"40" description:"Limit number of block that can be prefetched while listening"`
	}

	Tracker struct {
		Depth uint64 `long:"listener-tracker-depth" env:"LISTENER_TRACKER_DEPTH" default:"5" description:"Depth at which we consider a block final"`
	}

	Start struct {
		Default  string            `long:"listener-start-default" env:"LISTENER_START_DEFAULT" default:"oldest" description:"Block position the listener should start listening from 'latest', 'oldest', 'genesis'"`
		Specific map[string]string `long:"listener-start" env:"LISTENER_START" description:"Position listener should start listening from (format <chainID>:<blockNumber>-<txIndex> e.g. 0x2a:2348721-5)"`
	}
}

// Config worker configuration
type Config struct {
	Log      LoggerConfig
	Worker   WorkerConfig
	Kafka    KafkaConfig
	Eth      EthConfig
	Listener ListenerConfig
}

// TranslateBlockNumber translate a starting block number into its integer value
func TranslateBlockNumber(blockNumber string) (int64, error) {
	switch blockNumber {
	case "genesis":
		return 0, nil
	case "latest":
		return -1, nil
	case "oldest":
		return -2, nil
	default:
		res, err := strconv.ParseInt(blockNumber, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%q is an invalid starting blockNumber expected 'latest', 'oldest', 'genesis' or an integer", blockNumber)
		}
		return res, nil
	}
}

// ParseStartingPosition parse a starting position
func ParseStartingPosition(specific string) (int64, int64, error) {
	split := strings.Split(specific, "-")
	if len(split) != 2 {
		return 0, 0, fmt.Errorf("<%q is an invalid starting position (should match <blockNumber>-<txIndex>)", specific)
	}

	txIndex, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("<%q is an invalid starting position (should match <blockNumber>-<txIndex>)", specific)
	}

	blockNumber, err := TranslateBlockNumber(split[0])
	if err != nil {
		return 0, 0, err
	}

	return blockNumber, txIndex, nil

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
