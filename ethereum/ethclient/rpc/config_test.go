package rpc

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestEthClientURLs(t *testing.T) {
	name := "eth.clients"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	URLs(flgs)

	var expected []string
	if len(expected) != len(viper.GetStringSlice(name)) {
		t.Errorf("EthClientURLs #1: expected %v but got %v", expected, viper.GetStringSlice(name))
	} else {
		for i, url := range viper.GetStringSlice(name) {
			if url != expected[i] {
				t.Errorf("EthClientURLs #1: expected %v but got %v", expected, viper.GetStringSlice(name))
			}
		}
	}

	_ = os.Setenv("ETH_CLIENT_URL", "http://localhost:7546 http://localhost:8546")
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
	err := flgs.Parse(args)
	assert.Nil(t, err)

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
