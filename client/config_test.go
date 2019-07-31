package client

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGRPCStoreTarget(t *testing.T) {
	name := "grpc.target.envelope.store"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	EnvelopeStoreGRPCTarget(flgs)
	expected := "localhost:8080"
	assert.Equal(t, expected, viper.GetString(name), "Default")

	os.Setenv("GRPC_TARGET_ENVELOPE_STORE", "env-grpc-store")
	expected = "env-grpc-store"
	assert.Equal(t, expected, viper.GetString(name), "From Environment Variable")
	os.Unsetenv("GRPC_TARGET_ENVELOPE_STORE")

	args := []string{
		"--grpc-target-envelope-store=flag-grpc-store",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "Parse Store flags should not error")
	expected = "flag-grpc-store"
	assert.Equal(t, expected, viper.GetString(name), "From Flag")
}
