package http

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func TestHTTPHostname(t *testing.T) {
	name := "http.hostname"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	HTTPHostname(flgs)
	expected := ":8080"
	if viper.GetString(name) != expected {
		t.Errorf("HTTPHostname #1: expected %q but got %q", expected, viper.GetString(name))
	}

	os.Setenv("HTTP_HOSTNAME", "localhost:3000")
	expected = "localhost:3000"
	if viper.GetString(name) != expected {
		t.Errorf("HTTPHostname #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--http-hostname=127.0.0.1:3000",
	}
	flgs.Parse(args)
	expected = "127.0.0.1:3000"
	if viper.GetString(name) != expected {
		t.Errorf("HTTPHostname #3: expected %q but got %q", expected, viper.GetString(name))
	}
}
