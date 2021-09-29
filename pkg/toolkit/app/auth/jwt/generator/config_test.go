package generator

import (
	"os"
	"testing"

	tlstestutils "github.com/consensys/orchestrate/pkg/toolkit/tls/testutils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAuthServicePrivateKey(t *testing.T) {
	name := "auth.jwt.private.key"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	PrivateKey(flgs)

	_ = os.Setenv("AUTH_JWT_PRIVATE_KEY", tlstestutils.OneLineRSAKeyPEMA)
	assert.Equal(t, tlstestutils.OneLineRSAKeyPEMA, viper.GetString(name), "TenancyEnable #1")
}
