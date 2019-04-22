package registry

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestABIs(t *testing.T) {
	name := "abis"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	ABIs(flgs)

	// Test default
	expected := abiDefault
	assert.Equal(t, expected, viper.GetStringSlice(name), "Default config should match")

	// Test environment variable
	os.Setenv("ABI", "ERC20[v0.1.3]:[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}] ERC1400[v0.1.3]:[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]")

	expected = []string{
		"ERC20[v0.1.3]:[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
		"ERC1400[v0.1.3]:[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
	}
	assert.Equal(t, expected, viper.GetStringSlice(name), "Changing env var should change ABIs")
	os.Unsetenv("ABI")

	// Test flags
	args := []string{
		"--abi=MyContract[v1]:[ABI1]",
		"--abi=MyContract[v2]:[ABI2]",
	}
	flgs.Parse(args)

	t.Logf("Flags: %v", len(viper.GetStringSlice(name)))
	expected = []string{
		"MyContract[v1]:[ABI1]",
		"MyContract[v2]:[ABI2]",
	}
	assert.Equal(t, expected, viper.GetStringSlice(name), "Changing flags should change ABIs")
}

func TestFromABIConfig(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	ABIs(flgs)

	contracts, err := FromABIConfig()

	assert.Nil(t, err, "Should parse default properly")
	assert.Len(t, contracts, 1, "Expected 1 contract")
}
