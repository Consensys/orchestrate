package logger

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLogLevel(t *testing.T) {
	name := "log.level"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	LogLevel(flgs)

	expected := "debug"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #1: expected %q but got %q", expected, viper.GetString(name))
	}

	os.Setenv("LOG_LEVEL", "fatal")
	expected = "fatal"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--log-level=text",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "text"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #3: expected %q but got %q", expected, viper.GetString(name))
	}
}

func TestLogFormat(t *testing.T) {
	name := "log.format"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	LogFormat(flgs)
	expected := "text"
	if viper.GetString(name) != expected {
		t.Errorf("LogFormat #1: expected %q but got %q", expected, viper.GetString(name))
	}

	os.Setenv("LOG_FORMAT", "json")
	expected = "json"
	if viper.GetString(name) != expected {
		t.Errorf("LogFormat #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--log-format=xml",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "xml"
	if viper.GetString(name) != expected {
		t.Errorf("LogFormat #3: expected %q but got %q", expected, viper.GetString(name))
	}
}
