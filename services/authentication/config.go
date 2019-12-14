package authentication

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(AuthServiceCertificateViperKey, authServiceCertificateDefault)
	_ = viper.BindEnv(AuthServiceCertificateViperKey, authServiceCertificateEnv)
}

// Provision trusted certificate of the authentication service (base64 encoded)
const (
	authServiceCertificateFlag     = "auth-service-certificate"
	AuthServiceCertificateViperKey = "auth.service.certificate"
	authServiceCertificateDefault  = ""
	authServiceCertificateEnv      = "AUTH_SERVICE_CERTIFICATE"
)

// AuthServiceCertificate register flag for Authentication service Certificate
func AuthServiceCertificate(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Certificate of the authentication service encoded in base64
Environment variable: %q`, authServiceCertificateEnv)
	f.String(authServiceCertificateFlag, authServiceCertificateDefault, desc)
	_ = viper.BindPFlag(AuthServiceCertificateViperKey, f.Lookup(authServiceCertificateFlag))
}
