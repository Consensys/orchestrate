package authentication

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(AuthServiceCertificateViperKey, authServiceCertificateDefault)
	_ = viper.BindEnv(AuthServiceCertificateViperKey, authServiceCertificateEnv)
	viper.SetDefault(TenantNamespaceViperKey, tenantNamespaceDefault)
	_ = viper.BindEnv(TenantNamespaceViperKey, tenantNamespaceEnv)
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

// Provision tenant namespace to retrieve the tenant id in the OpenId or Access Token (JWT)
const (
	tenantNamespaceFlag     = "tenant-namespace"
	TenantNamespaceViperKey = "tenant.namespace"
	tenantNamespaceDefault  = "http://tenant.info"
	tenantNamespaceEnv      = "TENANT_NAMESPACE"
)

// TenantNamespace register flag for tenant namespace
func TenantNamespace(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Tenant Namespace to retrieve the tenant id in the OpenId or Access Token (JWT)
Environment variable: %q`, tenantNamespaceEnv)
	f.String(tenantNamespaceFlag, tenantNamespaceDefault, desc)
	_ = viper.BindPFlag(TenantNamespaceViperKey, f.Lookup(tenantNamespaceFlag))
}
