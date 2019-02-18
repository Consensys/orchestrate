package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func TestLogLevel(t *testing.T) {
	name := "log.level"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	LogLevel(flgs)

	expected := "debug"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #1: expected %q but got %q", expected, viper.GetString(name))
	}

	os.Setenv("LOG_LEVEL", "fatal")
	expected = "fatal"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--log-level=text",
	}
	flgs.Parse(args)
	expected = "text"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #3: expected %q but got %q", expected, viper.GetString(name))
	}

}

func TestLogFormat(t *testing.T) {
	name := "log.format"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	LogFormat(flgs)
	expected := "text"
	if viper.GetString(name) != expected {
		t.Errorf("LogFormat #1: expected %q but got %q", expected, viper.GetString(name))
	}

	os.Setenv("LOG_FORMAT", "json")
	expected = "json"
	if viper.GetString(name) != expected {
		t.Errorf("LogFormat #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--log-format=xml",
	}
	flgs.Parse(args)
	expected = "xml"
	if viper.GetString(name) != expected {
		t.Errorf("LogFormat #3: expected %q but got %q", expected, viper.GetString(name))
	}
}

func TestEthClientURLs(t *testing.T) {
	name := "eth.clients"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	EthClientURLs(flgs)

	expected := []string{
		"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7",
		"https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
	}
	if len(expected) != len(viper.GetStringSlice(name)) {
		t.Errorf("EthClientURLs #1: expected %v but got %v", expected, viper.GetStringSlice(name))
	} else {
		for i, url := range viper.GetStringSlice(name) {
			if url != expected[i] {
				t.Errorf("EthClientURLs #1: expected %v but got %v", expected, viper.GetStringSlice(name))
			}
		}
	}

	os.Setenv("ETH_CLIENT_URL", "http://localhost:7546 http://localhost:8546")
	expected = []string{
		"http://localhost:7546",
		"http://localhost:8546",
	}

	if len(expected) != len(viper.GetStringSlice(name)) {
		t.Errorf("EthClientURLs #2: expect %v but got %v", expected, viper.GetStringSlice(name))
	} else {
		for i, url := range viper.GetStringSlice(name) {
			if url != expected[i] {
				t.Errorf("EthClientURLs #2: expect %v but got %v", expected, viper.GetStringSlice(name))
			}
		}
	}

	args := []string{
		"--eth-client=http://localhost:6546",
		"--eth-client=http://localhost:7546,http://localhost:8646",
	}
	flgs.Parse(args)

	expected = []string{
		"http://localhost:6546",
		"http://localhost:7546",
		"http://localhost:8646",
	}
	if len(expected) != len(viper.GetStringSlice(name)) {
		t.Errorf("EthClientURLs #2: expect %v but got %v", expected, viper.GetStringSlice(name))
	} else {
		for i, url := range viper.GetStringSlice(name) {
			if url != expected[i] {
				t.Errorf("EthClientURLs #2: expect %v but got %v", expected, viper.GetStringSlice(name))
			}
		}
	}
}

func TestKafkaAddresses(t *testing.T) {
	name := "kafka.addresses"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	KafkaAddresses(flgs)

	expected := []string{
		"localhost:9092",
	}
	if len(expected) != len(viper.GetStringSlice(name)) {
		t.Errorf("KafkaAddresses #1: expected %v but got %v", expected, viper.GetStringSlice(name))
	} else {
		for i, url := range viper.GetStringSlice(name) {
			if url != expected[i] {
				t.Errorf("KafkaAddresses #1: expected %v but got %v", expected, viper.GetStringSlice(name))
			}
		}
	}

	os.Setenv("KAFKA_ADDRESS", "localhost:9192")
	expected = []string{
		"localhost:9192",
	}
	if len(expected) != len(viper.GetStringSlice(name)) {
		t.Errorf("EthClientURLs #2: expect %v but got %v", expected, viper.GetStringSlice(name))
	} else {
		for i, url := range viper.GetStringSlice(name) {
			if url != expected[i] {
				t.Errorf("EthClientURLs #2: expect %v but got %v", expected, viper.GetStringSlice(name))
			}
		}
	}

	args := []string{
		"--kafka-address=127.0.0.1:9091",
		"--kafka-address=127.0.0.2:9091",
	}
	flgs.Parse(args)

	expected = []string{
		"127.0.0.1:9091",
		"127.0.0.2:9091",
	}
	if len(expected) != len(viper.GetStringSlice(name)) {
		t.Errorf("EthClientURLs #3: expect %v but got %v", expected, viper.GetStringSlice(name))
	} else {
		for i, url := range viper.GetStringSlice(name) {
			if url != expected[i] {
				t.Errorf("EthClientURLs #3: expect %v but got %v", expected, viper.GetStringSlice(name))
			}
		}
	}
}

func TestWorkerInTopic(t *testing.T) {
	name := "worker.in"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	WorkerInTopic(flgs, "TOPIC_IN", "test-in-topic")
	expected := "test-in-topic"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #1: expected %q but got %q", expected, viper.GetString(name))
	}

	os.Setenv("TOPIC_IN", "test-env-in-topic")
	expected = "test-env-in-topic"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--worker-in=test-flag-in-topic",
	}
	flgs.Parse(args)
	expected = "test-flag-in-topic"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #3: expected %q but got %q", expected, viper.GetString(name))
	}
}

func TestWorkerOutTopic(t *testing.T) {
	name := "worker.out"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	WorkerOutTopic(flgs, "TOPIC_OUT", "test-out-topic")
	expected := "test-out-topic"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #1: expected %q but got %q", expected, viper.GetString(name))
	}

	os.Setenv("TOPIC_OUT", "test-env-out-topic")
	expected = "test-env-out-topic"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--worker-out=test-flag-out-topic",
	}
	flgs.Parse(args)
	expected = "test-flag-out-topic"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #3: expected %q but got %q", expected, viper.GetString(name))
	}
}

func TestWorkerConsumerGroup(t *testing.T) {
	name := "worker.group"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	WorkerConsumerGroup(flgs, "CONSUMER_GROUP", "test-consumer-group")
	expected := "test-consumer-group"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #1: expected %q but got %q", expected, viper.GetString(name))
	}

	os.Setenv("CONSUMER_GROUP", "test-env-consumer-group")
	expected = "test-env-consumer-group"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--worker-group=test-flag-consumer-group",
	}
	flgs.Parse(args)
	expected = "test-flag-consumer-group"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #3: expected %q but got %q", expected, viper.GetString(name))
	}
}

func TestWorkerSlots(t *testing.T) {
	name := "worker.slots"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	WorkerSlots(flgs)
	expected := 100
	if viper.GetInt(name) != expected {
		t.Errorf("LogLevel #1: expected %q but got %q", expected, viper.GetString(name))
	}

	os.Setenv("WORKER_SLOTS", "125")
	expected = 125
	if viper.GetInt(name) != expected {
		t.Errorf("LogLevel #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--worker-slots=150",
	}
	flgs.Parse(args)
	expected = 150
	if viper.GetInt(name) != expected {
		t.Errorf("LogLevel #3: expected %q but got %q", expected, viper.GetString(name))
	}
}

func TestRedisAddress(t *testing.T) {
	name := "redis.address"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	RedisAddress(flgs)
	expected := "localhost:6379"
	if viper.GetString(name) != expected {
		t.Errorf("RedisAddress #1: expected %q but got %q", expected, viper.GetString(name))
	}

	os.Setenv("REDIS_ADDRESS", "127.0.0.1:6378")
	expected = "127.0.0.1:6378"
	if viper.GetString(name) != expected {
		t.Errorf("RedisAddress #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--redis-address=127.0.0.1:6379",
	}
	flgs.Parse(args)
	expected = "127.0.0.1:6379"
	if viper.GetString(name) != expected {
		t.Errorf("RedisAddress #3: expected %q but got %q", expected, viper.GetString(name))
	}
}
