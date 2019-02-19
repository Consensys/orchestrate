package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

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
