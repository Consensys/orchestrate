package client

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGRPCStoreTarget(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	EnvelopeStoreGRPCTarget(flgs)
	expected := "localhost:8080"
	assert.Equal(t, expected, viper.GetString(grpcTargetEnvelopeStoreViperKey), "Default")

	_ = os.Setenv(grpcTargetEnvelopeStoreEnv, "env-grpc-store")
	expected = "env-grpc-store"
	assert.Equal(t, expected, viper.GetString(grpcTargetEnvelopeStoreViperKey), "From Environment Variable")
	_ = os.Unsetenv(grpcTargetEnvelopeStoreEnv)

	args := []string{
		"--grpc-target-envelope-store=flag-grpc-store",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "Parse Store flags should not error")
	expected = "flag-grpc-store"
	assert.Equal(t, expected, viper.GetString(grpcTargetEnvelopeStoreViperKey), "From Flag")
}
