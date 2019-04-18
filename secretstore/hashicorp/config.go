package hashicorp

import (
	"fmt"
	"os"
	log "github.com/sirupsen/logrus"
	vault "github.com/hashicorp/vault/api"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(vaultSecretPathViperKey, vaultSecretPathDefault)
	viper.BindEnv(vaultSecretPathViperKey, vaultSecretPathEnv)

	log.Infof("Using vault_address: %v \n", os.Getenv("VAULT_ADDR"))
	log.Infof("Using vault_agent_address: %v \n", os.Getenv("VAULT_AGENT_ADDR"))
}

var (
	
	vaultSecretPathEnv 				= "VAULT_SECRET_PATH"
	vaultRateLimitEnv 				= "VAULT_RATE_LIMIT"
	vaultAddressEnv 				= "VAULT_ADDR"
	vaultAgentAddrEnv 				= "VAULT_AGENT_ADDR"
	vaultCACertEnv 					= "VAULT_CACERT"
	vaultCAPathEnv 					= "VAULT_CAPATH"
	vaultClientCertEnv 				= "VAULT_CLIENT_CERT"
	vaultClientKeyEnv 				= "VAULT_CLIENT_KEY"
	vaultClientTimeoutEnv 			= "VAULT_CLIENT_TIMEOUT"
	vaultMFAEnv 					= "VAULT_MFA"
	vaultMaxRetriesEnv 				= "VAULT_MAX_RETRIES"
	vaultNamespaceEnv 				= "VAULT_NAMESPACE"
	vaultSkipVerifyEnv 				= "VAULT_SKIP_VERIFY"
	vaultTLSServerNameEnv		 	= "VAULT_TLS_SERVER_NAME"
	vaultTokenEnv		 			= "VAULT_TOKEN"
	vaultWrapTTLEnv 				= "VAULT_WRAP_TTL"

	vaultSecretPathFlag				= "vault-secret-path"
	vaultRateLimitFlag 				= "vault-rate-limit"
	vaultAddressFlag 				= "vault-addr"
	vaultAgentAddrFlag 				= "vault-agent-addr"
	vaultCACertFlag 				= "vault-cacert"
	vaultCAPathFlag 				= "vault-capath"
	vaultClientCertFlag 			= "vault-client-cert"
	vaultClientKeyFlag 				= "vault-client-key"
	vaultClientTimeoutFlag 			= "vault-client-timeout"
	vaultMFAFlag 					= "vault-mfa"
	vaultMaxRetriesFlag 			= "vault-max-retries"
	vaultNamespaceFlag 				= "vault-namespace"
	vaultSkipVerifyFlag 			= "vault-skip-verify"
	vaultTLSServerNameFlag		 	= "vault-tls-server-name"
	vaultTokenFlag		 			= "vault-token"
	vaultWrapTTLFlag 				= "vault-wrap-ttl"

	vaultSecretPathViperKey			= "vault.secret.path"
	vaultRateLimitViperKey 			= "vault.rate.limit"
	vaultAddressViperKey 			= "vault.addr"
	vaultAgentAddrViperKey 			= "vault.agent.addr"
	vaultCACertViperKey 			= "vault.cacert"
	vaultCAPathViperKey 			= "vault.capath"
	vaultClientCertViperKey 		= "vault.client.cert"
	vaultClientKeyViperKey 			= "vault.client.key"
	vaultClientTimeoutViperKey 		= "vault.client.timeout"
	vaultMFAViperKey 				= "vault.mfa"
	vaultMaxRetriesViperKey 		= "vault.max.retries"
	vaultNamespaceViperKey 			= "vault.namespace"
	vaultSkipVerifyViperKey 		= "vault.skip.verify"
	vaultTLSServerNameViperKey		= "vault.tls.server.name"
	vaultTokenViperKey		 		= "vault.token"
	vaultWrapTTLViperKey 			= "vault.wrap.ttl"

	// No need to redefine the default here
	vaultSecretPathDefault			= "/secret"
	vaultRateLimitDefault 			= ""
	vaultAddressDefault 			= ""
	vaultAgentAddrDefault 			= ""
	vaultCACertDefault 				= ""
	vaultCAPathDefault 				= ""
	vaultClientCertDefault 			= ""
	vaultClientKeyDefault 			= ""
	vaultClientTimeoutDefault 		= ""
	vaultMFADefault 				= ""
	vaultMaxRetriesDefault 			= ""
	vaultNamespaceDefault 			= ""
	vaultSkipVerifyDefault 			= ""
	vaultTLSServerNameDefault		= ""
	vaultTokenDefault		 		= ""
	vaultWrapTTLDefault 			= ""

)

// InitFlags register flags for hashicorp vault
func InitFlags(f *pflag.FlagSet) {
	VaultSecretPath(f)
	VaultRateLimit(f)
	VaultAddress(f)
	VaultAgentAddr(f)
	VaultCACert(f)
	VaultCAPath(f)
	VaultClientCert(f)
	VaultClientKey(f)
	VaultClientTimeout(f)
	VaultMFA(f)
	VaultNamespace(f)
	VaultSkipVerify(f)
	VaultTLSServerName(f)
	VaultToken(f)
	VaultWrapTTL(f)
}

// VaultSecretPath registers a flag for the path used by vault secret engine
func VaultSecretPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultSecretPathEnv)
	f.String(vaultSecretPathFlag, vaultSecretPathDefault, desc)
	viper.BindPFlag(vaultSecretPathViperKey, f.Lookup(vaultSecretPathFlag))
}

// VaultRateLimit registers a flag for the path used by vault secret engine
func VaultRateLimit(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultRateLimitEnv)
	f.String(vaultRateLimitFlag, vaultRateLimitDefault, desc)
	viper.BindPFlag(vaultRateLimitViperKey, f.Lookup(vaultRateLimitFlag))
}

// VaultAddress registers a flag for the path used by vault secret engine
func VaultAddress(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultAddressEnv)
	f.String(vaultAddressFlag, vaultAddressDefault, desc)
	viper.BindPFlag(vaultAddressViperKey, f.Lookup(vaultAddressFlag))
}

// VaultAgentAddr registers a flag for the path used by vault secret engine
func VaultAgentAddr(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultAgentAddrEnv)
	f.String(vaultAgentAddrFlag, vaultAgentAddrDefault, desc)
	viper.BindPFlag(vaultAgentAddrViperKey, f.Lookup(vaultAgentAddrFlag))
}

// VaultCACert registers a flag for the path used by vault secret engine
func VaultCACert(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultCACertEnv)
	f.String(vaultCACertFlag, vaultCACertDefault, desc)
	viper.BindPFlag(vaultCACertViperKey, f.Lookup(vaultCACertFlag))
}

// VaultCAPath registers a flag for the path used by vault secret engine
func VaultCAPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultCAPathEnv)
	f.String(vaultCAPathFlag, vaultCAPathDefault, desc)
	viper.BindPFlag(vaultCAPathViperKey, f.Lookup(vaultCAPathFlag))
}

// VaultClientCert registers a flag for the path used by vault secret engine
func VaultClientCert(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultClientCertEnv)
	f.String(vaultClientCertFlag, vaultClientCertDefault, desc)
	viper.BindPFlag(vaultClientCertViperKey, f.Lookup(vaultClientCertFlag))
}

// VaultClientKey registers a flag for the path used by vault secret engine
func VaultClientKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultClientKeyEnv)
	f.String(vaultClientKeyFlag, vaultClientKeyDefault, desc)
	viper.BindPFlag(vaultClientKeyViperKey, f.Lookup(vaultClientKeyFlag))
}

// VaultClientTimeout registers a flag for the path used by vault secret engine
func VaultClientTimeout(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultClientTimeoutEnv)
	f.String(vaultClientTimeoutFlag, vaultClientTimeoutDefault, desc)
	viper.BindPFlag(vaultClientTimeoutViperKey, f.Lookup(vaultClientTimeoutFlag))
}

// VaultMFA registers a flag for the path used by vault secret engine
func VaultMFA(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultMFAEnv)
	f.String(vaultMFAFlag, vaultMFADefault, desc)
	viper.BindPFlag(vaultMFAViperKey, f.Lookup(vaultMFAFlag))
}

// VaultMaxRetries registers a flag for the path used by vault secret engine
func VaultMaxRetries(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultMaxRetriesEnv)
	f.String(vaultMaxRetriesFlag, vaultMaxRetriesDefault, desc)
	viper.BindPFlag(vaultMaxRetriesViperKey, f.Lookup(vaultMaxRetriesFlag))
}

// VaultNamespace registers a flag for the path used by vault secret engine
func VaultNamespace(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultNamespaceEnv)
	f.String(vaultNamespaceFlag, vaultNamespaceDefault, desc)
	viper.BindPFlag(vaultNamespaceViperKey, f.Lookup(vaultNamespaceFlag))
}

// VaultSkipVerify registers a flag for vault client
func VaultSkipVerify(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultSkipVerifyEnv)
	f.String(vaultSkipVerifyFlag, vaultSkipVerifyDefault, desc)
	viper.BindPFlag(vaultSkipVerifyViperKey, f.Lookup(vaultSkipVerifyFlag))
}

// VaultTLSServerName registers a flag for vault client
func VaultTLSServerName(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultTLSServerNameEnv)
	f.String(vaultTLSServerNameFlag, vaultTLSServerNameDefault, desc)
	viper.BindPFlag(vaultTLSServerNameViperKey, f.Lookup(vaultTLSServerNameFlag))
}

// VaultToken registers a flag for vault client
func VaultToken(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultTokenEnv)
	f.String(vaultTokenFlag, vaultTokenDefault, desc)
	viper.BindPFlag(vaultTokenViperKey, f.Lookup(vaultTokenFlag))
}

// VaultWrapTTL registers a flag for vault client
func VaultWrapTTL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultWrapTTLEnv)
	f.String(vaultWrapTTLFlag, vaultWrapTTLDefault, desc)
	viper.BindPFlag(vaultWrapTTLViperKey, f.Lookup(vaultWrapTTLFlag))
}

// NewConfig creates vault configuration from viper
func NewConfig() *vault.Config {
	config := vault.DefaultConfig()
	return config
}

// GetSecretPath returns the secret path set in deployment by vault
func GetSecretPath() string {
	return viper.GetString(vaultSecretPathViperKey)
}