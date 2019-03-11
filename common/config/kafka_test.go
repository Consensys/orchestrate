package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

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
