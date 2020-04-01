// +build unit

package secretstore

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSecretStore(t *testing.T) {
	name := "secret.store"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	SecFlag(flgs)
	expected := "in-memory"
	assert.Equal(t, expected, viper.GetString(name), "Default")

	_ = os.Setenv("SECRET_STORE", "env-store")
	expected = "env-store"
	assert.Equal(t, expected, viper.GetString(name), "From Environment Variable")
	_ = os.Unsetenv("SECRET_STORE")

	args := []string{
		"--secret-store=flag-store",
	}
	_ = flgs.Parse(args)
	expected = "flag-store"
	assert.Equal(t, expected, viper.GetString(name), "From Flag")
}
