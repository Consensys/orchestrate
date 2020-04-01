// +build unit

package key

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAuthKey(t *testing.T) {
	name := "auth.api-key"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	APIKey(flgs)

	_ = os.Setenv("AUTH_API_KEY", "test-key")
	assert.Equal(t, "test-key", viper.GetString(name), "TenancyEnable #1")
}
