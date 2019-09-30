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
	ContractRegistryGRPCTarget(flgs)
	expected := "localhost:8080"
	assert.Equal(t, expected, viper.GetString(grpcTargetContractRegistryViperKey), "Default")

	_ = os.Setenv(grpcTargetContractRegistryEnv, "env-grpc-contract-registry")
	expected = "env-grpc-contract-registry"
	assert.Equal(t, expected, viper.GetString(grpcTargetContractRegistryViperKey), "From Environment Variable")
	_ = os.Unsetenv(grpcTargetContractRegistryEnv)

	args := []string{
		"--grpc-target-contract-registry=flag-grpc-contract-registry",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "Parse Contract Registry flags should not error")
	expected = "flag-grpc-contract-registry"
	assert.Equal(t, expected, viper.GetString(grpcTargetContractRegistryViperKey), "From Flag")
}
