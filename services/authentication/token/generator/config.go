package generator

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(AuthServicePrivateKeyViperKey, authServicePrivateKeyDefault)
	_ = viper.BindEnv(AuthServicePrivateKeyViperKey, authServicePrivateKeyEnv)
}

// Provision trusted certificate of the authentication service (base64 encoded)
const (
	authServicePrivateKeyFlag     = "auth-service-private-key"
	AuthServicePrivateKeyViperKey = "auth.service.private.key"
	authServicePrivateKeyDefault  = ""
	authServicePrivateKeyEnv      = "AUTH_SERVICE_PRIVATE_KEY"
)

// AuthServicePrivateKey register flag for Authentication service Certificate
func AuthServicePrivateKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Private Key who sign the JWT encoded in base64
Environment variable: %q`, authServicePrivateKeyEnv)
	f.String(authServicePrivateKeyFlag, authServicePrivateKeyDefault, desc)
	_ = viper.BindPFlag(AuthServicePrivateKeyViperKey, f.Lookup(authServicePrivateKeyFlag))
}
