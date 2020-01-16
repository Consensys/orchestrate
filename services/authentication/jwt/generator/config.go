package generator

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(PrivateKeyViperKey, privateKeyDefault)
	_ = viper.BindEnv(PrivateKeyViperKey, privateKeyEnv)
}

// Provision trusted certificate of the authentication service (base64 encoded)
const (
	privateKeyFlag     = "auth-jwt-private-key"
	PrivateKeyViperKey = "auth.jwt.private.key"
	privateKeyDefault  = ""
	privateKeyEnv      = "AUTH_JWT_PRIVATE_KEY"
)

// PrivateKey register flag for Authentication service Certificate
func PrivateKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Private Key who sign the JWT encoded in base64
Environment variable: %q`, privateKeyEnv)
	f.String(privateKeyFlag, privateKeyDefault, desc)
	_ = viper.BindPFlag(PrivateKeyViperKey, f.Lookup(privateKeyFlag))
}
