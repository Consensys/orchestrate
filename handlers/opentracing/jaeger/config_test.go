package jaeger

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestHost(t *testing.T) {
	name := "jaeger.agent.host"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Host(flgs)
	expected := "localhost"
	assert.Equal(t, expected, viper.GetString(name), "Default")

	os.Setenv("JAEGER_AGENT_HOST", "env-jaeger")
	expected = "env-jaeger"
	assert.Equal(t, expected, viper.GetString(name), "From Environment Variable")
	os.Unsetenv("JAEGER_AGENT_HOST")

	args := []string{
		"--jaeger-host=flag-jaeger",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "flag-jaeger"
	assert.Equal(t, expected, viper.GetString(name), "From Flag")
}

func TestPort(t *testing.T) {
	name := "jaeger.agent.port"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Port(flgs)
	expected := 6831
	assert.Equal(t, expected, viper.GetInt(name), "Default")

	os.Setenv("JAEGER_AGENT_PORT", "5778")
	expected = 5778
	assert.Equal(t, expected, viper.GetInt(name), "From Environment Variable")
	os.Unsetenv("JAEGER_AGENT_PORT")

	args := []string{
		"--jaeger-port=5779",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = 5779
	assert.Equal(t, expected, viper.GetInt(name), "From Flag")
}

func TestSamplerParam(t *testing.T) {
	name := "jaeger.sampler.param"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	SamplerParam(flgs)
	expected := 1
	assert.Equal(t, expected, viper.GetInt(name), "Default")

	os.Setenv("JAEGER_SAMPLER_PARAM", "0")
	expected = 0
	assert.Equal(t, expected, viper.GetInt(name), "From Environment Variable")
	os.Unsetenv("JAEGER_HOST")

	args := []string{
		"--jaeger-sampler-param=0",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = 0
	assert.Equal(t, expected, viper.GetInt(name), "From Flag")
}
