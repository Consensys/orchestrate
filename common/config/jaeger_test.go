package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestJaegerHost(t *testing.T) {
	name := "jaeger.host"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	JaegerHost(flgs)
	expected := "jaeger"
	assert.Equal(t, expected, viper.GetString(name), "Default")

	os.Setenv("JAEGER_HOST", "env-jaeger")
	expected = "env-jaeger"
	assert.Equal(t, expected, viper.GetString(name), "From Environment Variable")
	os.Unsetenv("JAEGER_HOST")

	args := []string{
		"--jaeger-host=flag-jaeger",
	}
	flgs.Parse(args)
	expected = "flag-jaeger"
	assert.Equal(t, expected, viper.GetString(name), "From Flag")
}

func TestJaegerPort(t *testing.T) {
	name := "jaeger.port"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	JaegerPort(flgs)
	expected := 5775
	assert.Equal(t, expected, viper.GetInt(name), "Default")

	os.Setenv("JAEGER_PORT", "5778")
	expected = 5778
	assert.Equal(t, expected, viper.GetInt(name), "From Environment Variable")
	os.Unsetenv("JAEGER_HOST")

	args := []string{
		"--jaeger-port=5779",
	}
	flgs.Parse(args)
	expected = 5779
	assert.Equal(t, expected, viper.GetInt(name), "From Flag")
}

func TestJaegerSampler(t *testing.T) {
	name := "jaeger.sampler"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	JaegerSampler(flgs)
	expected := 0.5
	assert.Equal(t, expected, viper.GetFloat64(name), "Default")

	os.Setenv("JAEGER_SAMPLER", "1.6")
	expected = 1.6
	assert.Equal(t, expected, viper.GetFloat64(name), "From Environment Variable")
	os.Unsetenv("JAEGER_HOST")

	args := []string{
		"--jaeger-sampler=0.12345",
	}
	flgs.Parse(args)
	expected = 0.12345
	assert.Equal(t, expected, viper.GetFloat64(name), "From Flag")
}
