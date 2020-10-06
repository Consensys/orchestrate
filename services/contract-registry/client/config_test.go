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
	assert.Equal(t, expected, viper.GetString(GRPCURLViperKey), "Default")

	_ = os.Setenv(grpcURLEnv, "env-grpc-contract-registry")
	expected = "env-grpc-contract-registry"
	assert.Equal(t, expected, viper.GetString(GRPCURLViperKey), "From Environment Variable")
	_ = os.Unsetenv(grpcURLEnv)

	args := []string{
		"--contract-registry-url=flag-grpc-contract-registry",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "Parse Contract Registry flags should not error")
	expected = "flag-grpc-contract-registry"
	assert.Equal(t, expected, viper.GetString(GRPCURLViperKey), "From Flag")
}
