// +build unit

package keystore

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSecretPkeys(t *testing.T) {
	name := "secret.pkeys"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	SecretPkeys(flgs)
	var expected []string
	assert.Equal(t, expected, viper.GetStringSlice(name), "Default")

	_ = os.Setenv("SECRET_PKEY", "56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E 5FBB50BFF6DFAD35C4A374C9237BA2F7EAED9C6868E0108CB259B62D68029B1A")
	expected = []string{
		"56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E",
		"5FBB50BFF6DFAD35C4A374C9237BA2F7EAED9C6868E0108CB259B62D68029B1A",
	}
	assert.Equal(t, expected, viper.GetStringSlice(name), "From Environment Variable")
	_ = os.Unsetenv("SECRET_PKEY")

	args := []string{
		"--secret-pkey=86B021CCB810F26A30445B85F71E4C1596A11A97DDF9B9E348AC93D1DA6735BC",
	}
	_ = flgs.Parse(args)
	expected = []string{
		"86B021CCB810F26A30445B85F71E4C1596A11A97DDF9B9E348AC93D1DA6735BC",
	}
	assert.Equal(t, expected, viper.GetStringSlice(name), "From Flag")
}
