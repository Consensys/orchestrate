// +build unit

package client

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGRPCContractRegistryTarget(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	ContractRegistryURL(flgs)
	expected := "localhost:8080"
	assert.Equal(t, expected, viper.GetString(ContractRegistryURLViperKey), "Default")

	_ = os.Setenv(contractRegistryURLEnv, "env-grpc-contract-registry")
	expected = "env-grpc-contract-registry"
	assert.Equal(t, expected, viper.GetString(ContractRegistryURLViperKey), "From Environment Variable")
	_ = os.Unsetenv(contractRegistryURLEnv)

	args := []string{
		"--contract-registry-url=flag-grpc-contract-registry",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "Parse Contract Registry flags should not error")
	expected = "flag-grpc-contract-registry"
	assert.Equal(t, expected, viper.GetString(ContractRegistryURLViperKey), "From Flag")
}
