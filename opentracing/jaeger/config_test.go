package jaeger

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
	expected := 6831
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

func TestJaegerSamplerParam(t *testing.T) {
	name := "jaeger.sampler.param"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	JaegerSamplerParam(flgs)
	expected := 1
	assert.Equal(t, expected, viper.GetInt(name), "Default")

	os.Setenv("JAEGER_SAMPLER_PARAM", "0")
	expected = 0
	assert.Equal(t, expected, viper.GetInt(name), "From Environment Variable")
	os.Unsetenv("JAEGER_HOST")

	args := []string{
		"--jaeger-sampler-param=0",
	}
	flgs.Parse(args)
	expected = 0
	assert.Equal(t, expected, viper.GetInt(name), "From Flag")
}
