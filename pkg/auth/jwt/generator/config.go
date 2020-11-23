package generator

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/jwt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tls/certificate"
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
	desc := fmt.Sprintf(`Private key to sign generated JWT tokens
Environment variable: %q`, privateKeyEnv)
	f.String(privateKeyFlag, privateKeyDefault, desc)
	_ = viper.BindPFlag(PrivateKeyViperKey, f.Lookup(privateKeyFlag))
}

type Config struct {
	KeyPair         *certificate.KeyPair
	ClaimsNamespace string
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		KeyPair: &certificate.KeyPair{
			Cert: []byte(vipr.GetString(jwt.CertificateViperKey)),
			Key:  []byte(vipr.GetString(PrivateKeyViperKey)),
		},
		ClaimsNamespace: vipr.GetString(jwt.ClaimsNamespaceViperKey),
	}
}
