package multitenancy

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(EnabledViperKey, enabledDefault)
	_ = viper.BindEnv(EnabledViperKey, enabledEnv)
	viper.SetDefault(AuthServiceCertificateViperKey, authServiceCertificateDefault)
	_ = viper.BindEnv(AuthServiceCertificateViperKey, authServiceCertificateEnv)
}

// Enable or disable the multi-tenancy support process
const (
	enabledFlag     = "multi-tenancy-enabled"
	EnabledViperKey = "multi.tenancy.enabled"
	enabledDefault  = false
	enabledEnv      = "MULTI_TENANCY_ENABLED"
)

// TenancyEnable register flag for Enable Multi-Tenancy
func Enabled(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Whether or not to use Multi Tenancy (one of %q).
Environment variable: %q`, []string{"false", "true"}, enabledEnv)
	f.Bool(enabledFlag, enabledDefault, desc)
	_ = viper.BindPFlag(EnabledViperKey, f.Lookup(enabledFlag))
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
