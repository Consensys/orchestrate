package logger

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLogLevel(t *testing.T) {
	name := "log.level"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	LogLevel(flgs)

	expected := "info"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #1: expected %q but got %q", expected, viper.GetString(name))
	}

	_ = os.Setenv(logLevelEnv, "fatal")
	expected = "fatal"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #2: expected %q but got %q", expected, viper.GetString(name))
	}
	_ = os.Unsetenv(logLevelEnv)

	args := []string{
		"--log-level=debug",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "debug"
	if viper.GetString(name) != expected {
		t.Errorf("LogLevel #3: expected %q but got %q", expected, viper.GetString(name))
	}

	name = "log.format"
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	LogFormat(f)
	expected = "text"
	if viper.GetString(name) != expected {
		t.Errorf("LogFormat #1: expected %q but got %q", expected, viper.GetString(name))
	}

	_ = os.Setenv(logFormatEnv, "json")
	expected = "json"
	if viper.GetString(name) != expected {
		t.Errorf("LogFormat #2: expected %q but got %q", expected, viper.GetString(name))
	}
	_ = os.Unsetenv(logFormatEnv)

	args = []string{
		"--log-format=xml",
	}
	err = f.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "xml"
	if viper.GetString(name) != expected {
		t.Errorf("LogFormat #3: expected %q but got %q", expected, viper.GetString(name))
	}

	viper.Set(LogLevelViperKey, "info")
	viper.Set(LogFormatViperKey, "text")
	InitLogger()

	l := log.StandardLogger()
	assert.Equal(t, l.Formatter, &log.TextFormatter{}, "Should be the same formatter")
	assert.Equal(t, l.GetLevel(), log.InfoLevel, "Should be the same log level")

	viper.Set(LogLevelViperKey, "debug")
	viper.Set(LogFormatViperKey, "json")

	InitLogger()
	assert.Equal(t, l.Formatter, &log.JSONFormatter{}, "Should be the same formatter")
	assert.Equal(t, l.GetLevel(), log.DebugLevel, "Should be the same log level")

}
