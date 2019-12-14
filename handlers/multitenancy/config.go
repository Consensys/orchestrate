package multitenancy

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(EnabledViperKey, enabledDefault)
	_ = viper.BindEnv(EnabledViperKey, enabledEnv)
	viper.SetDefault(TenantNamespaceViperKey, tenantNamespaceDefault)
	_ = viper.BindEnv(TenantNamespaceViperKey, tenantNamespaceEnv)
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
