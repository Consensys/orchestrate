package txdecoder

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestDisableExternalTx(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	InitFlags(flgs)

	expected := false
	assert.Equal(t, expected, ExternalTxDisabled(), "Default")

	_ = os.Setenv("DISABLE_EXTERNAL_TX", "true")
	expected = true
	assert.Equal(t, expected, ExternalTxDisabled(), "Env")
	_ = os.Unsetenv("DISABLE_EXTERNAL_TX")

	args := []string{
		"--disable-external-tx=true",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = true
	assert.Equal(t, expected, ExternalTxDisabled(), "Flag")
}
