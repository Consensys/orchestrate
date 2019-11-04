package tessera

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func TestFlags(t *testing.T) {
	name := "tessera.urls"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	InitFlags(flgs)

	expected := map[string]string{}
	if len(expected) != len(viper.GetStringMapString(name)) {
		t.Errorf("TesseraURLs #1: expected %v but got %v", expected, viper.GetStringMapString(name))
	} else {
		for _, chainID := range viper.GetStringMapString(name) {
			if viper.GetStringMapString(name)[chainID] != expected[chainID] {
				t.Errorf("TesseraURLs #1: expected %v but got %v", expected, viper.GetStringMapString(name))
			}
		}
	}

	_ = os.Setenv("TESSERA_URL",
		"{\"10\": \"http://tessera1:9080\", \"22\": \"Somewhere over the rainbow\", \"888\": \"http://localhost:80\"}",
	)

	expected = map[string]string{
		"10":  "http://tessera1:9080",
		"22":  "Somewhere over the rainbow",
		"888": "http://localhost:80",
	}

	if len(expected) != len(viper.GetStringMapString(name)) {
		t.Errorf("TesseraURLs #2: expect %v but got %v", expected, viper.GetStringMapString(name))
	} else {
		for _, chainID := range viper.GetStringMapString(name) {
			if viper.GetStringMapString(name)[chainID] != expected[chainID] {
				t.Errorf("TesseraURLs #2: expect %v but got %v", expected, viper.GetStringMapString(name))
			}
		}
	}
}
