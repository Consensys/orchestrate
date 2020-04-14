// +build unit

package client

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestClientFlags(t *testing.T) {
	flgs := pflag.NewFlagSet("TestClientFlags", pflag.ContinueOnError)
	Flags(flgs)
	expected := "localhost:8080"
	assert.Equal(t, expected, viper.GetString(EnvelopeStoreURLViperKey), "Default")

	_ = os.Setenv(envelopeStoreURLEnv, "env-grpc-store")
	expected = "env-grpc-store"
	assert.Equal(t, expected, viper.GetString(EnvelopeStoreURLViperKey), "From Environment Variable")
	_ = os.Unsetenv(envelopeStoreURLEnv)

	args := []string{
		"--envelope-store-url=flag-grpc-store",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "Parse Store flags should not error")
	expected = "flag-grpc-store"
	assert.Equal(t, expected, viper.GetString(EnvelopeStoreURLViperKey), "From Flag")
}

func TestClientDefaultConfig(t *testing.T) {
	flgs := pflag.NewFlagSet("TestClientDefaultConfig", pflag.ContinueOnError)
	Flags(flgs)
	
	cfg := NewConfigFromViper(viper.GetViper())
	
	assert.Equal(t, envelopeStoreURLDefault, cfg.envelopeStoreURL, "Default store url")
	assert.Equal(t, "jaeger", cfg.serviceName, "Default service name")
}

func TestClientNewConfig(t *testing.T) {
	flgs := pflag.NewFlagSet("TestClientNewConfig", pflag.ContinueOnError)
	Flags(flgs)
	
	expectedStoreUrl := "flag-grpc-store:8080"
	args := []string{
		"--envelope-store-url=" + expectedStoreUrl,
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "Parse Store flags should not error")
	
	cfg := NewConfigFromViper(viper.GetViper())
	assert.Equal(t, expectedStoreUrl, cfg.envelopeStoreURL, "Default store url")
}
