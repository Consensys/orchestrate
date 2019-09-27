package engine

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSlots(t *testing.T) {
	name := "engine.slots"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Slots(flgs)
	expected := 20
	assert.Equal(t, expected, viper.GetInt(name), "Default")

	os.Setenv("ENGINE_SLOTS", "125")
	expected = 125
	assert.Equal(t, expected, viper.GetInt(name), "From Environment Variable")
	os.Unsetenv("ENGINE_SLOTS")

	args := []string{
		"--engine-slots=150",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = 150
	assert.Equal(t, expected, viper.GetInt(name), "From Flag")
}
