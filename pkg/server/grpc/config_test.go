package grpcserver

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestHostname(t *testing.T) {
	name := "grpc.hostname"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Hostname(flgs)
	expected := ""
	if viper.GetString(name) != expected {
		t.Errorf("GRPC Hostname #1: expected %q but got %q", expected, viper.GetString(name))
	}

	_ = os.Setenv("GRPC_HOSTNAME", "localhost")
	expected = "localhost"
	if viper.GetString(name) != expected {
		t.Errorf("GRPC Hostname #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--grpc-hostname=127.0.0.1",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "127.0.0.1"
	if viper.GetString(name) != expected {
		t.Errorf("GRPC Hostname #3: expected %q but got %q", expected, viper.GetString(name))
	}
}

func TestMetricsPort(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Port(flgs)
	assert.Equal(t, uint64(8080), viper.GetUint64("grpc.port"), "Default Port should be correct")
}
