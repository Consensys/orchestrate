// +build unit

package http

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.NotNil(t, cfg)
}

func TestHostname(t *testing.T) {
	name := "rest.hostname"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Hostname(flgs)
	expected := ""
	if viper.GetString(name) != expected {
		t.Errorf("Hostname #1: expected %q but got %q", expected, viper.GetString(name))
	}

	_ = os.Setenv("REST_HOSTNAME", "localhost")
	expected = "localhost"
	if viper.GetString(name) != expected {
		t.Errorf("Hostname #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--rest-hostname=127.0.0.1",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "127.0.0.1"
	if viper.GetString(name) != expected {
		t.Errorf("Hostname #3: expected %q but got %q", expected, viper.GetString(name))
	}
}

func TestHTTPPort(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Port(flgs)
	assert.Equal(t, uint64(8081), viper.GetUint64("rest.port"), "Default Port should be correct")
}

func TestMetricsHostname(t *testing.T) {
	name := "metrics.hostname"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	MetricsHostname(flgs)
	expected := ""
	if viper.GetString(name) != expected {
		t.Errorf("Metrics Hostname #1: expected %q but got %q", expected, viper.GetString(name))
	}

	_ = os.Setenv("METRICS_HOSTNAME", "localhost")
	expected = "localhost"
	if viper.GetString(name) != expected {
		t.Errorf("Metrics Hostname #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--metrics-hostname=127.0.0.1",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "127.0.0.1"
	if viper.GetString(name) != expected {
		t.Errorf("Metrics Hostname #3: expected %q but got %q", expected, viper.GetString(name))
	}
}

func TestMetricsPort(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	MetricsPort(flgs)
	assert.Equal(t, uint64(8082), viper.GetUint64("metrics.port"), "Default Port should be correct")
}

func TestFlags(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(f)
}
