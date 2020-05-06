// +build unit

package log

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLevel(t *testing.T) {
	name := "log.level"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Level(flgs)

	expected := "info"
	if viper.GetString(name) != expected {
		t.Errorf("Level #1: expected %q but got %q", expected, viper.GetString(name))
	}

	_ = os.Setenv("LOG_LEVEL", "fatal")
	expected = "fatal"
	if viper.GetString(name) != expected {
		t.Errorf("Level #2: expected %q but got %q", expected, viper.GetString(name))
	}
	_ = os.Unsetenv("LOG_LEVEL")

	args := []string{
		"--log-level=debug",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "debug"
	if viper.GetString(name) != expected {
		t.Errorf("Level #3: expected %q but got %q", expected, viper.GetString(name))
	}

	name = "log.format"
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Format(f)
	expected = "text"
	if viper.GetString(name) != expected {
		t.Errorf("Format #1: expected %q but got %q", expected, viper.GetString(name))
	}

	_ = os.Setenv("LOG_FORMAT", "json")
	expected = "json"
	if viper.GetString(name) != expected {
		t.Errorf("Format #2: expected %q but got %q", expected, viper.GetString(name))
	}
	_ = os.Unsetenv("LOG_FORMAT")

	args = []string{
		"--log-format=xml",
	}
	err = f.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "xml"
	if viper.GetString(name) != expected {
		t.Errorf("Format #3: expected %q but got %q", expected, viper.GetString(name))
	}
}
