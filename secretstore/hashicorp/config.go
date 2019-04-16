package hashicorp

import (
	"fmt"
	"os"

	vault "github.com/hashicorp/vault/api"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	fmt.Printf("Vault_token: %v \n", os.Getenv("VAULT_TOKEN"))
	viper.SetDefault(vaultSecretPathViperKey, vaultSecretPathDefault)
	viper.BindEnv(vaultSecretPathViperKey, vaultSecretPathEnv)
}

/*
	The following variables are directly made available to vault

	Defined in https://github.com/hashicorp/vault/blob/master/api/client.go#L42

	const EnvRateLimit 				= "VAULT_RATE_LIMIT"
	const EnvVaultAddress 			= "VAULT_ADDR"
	const EnvVaultAgentAddr 		= "VAULT_AGENT_ADDR"
	const EnvVaultCACert 			= "VAULT_CACERT"
	const EnvVaultCAPath 			= "VAULT_CAPATH"
	const EnvVaultClientCert 		= "VAULT_CLIENT_CERT"
	const EnvVaultClientKey 		= "VAULT_CLIENT_KEY"
	const EnvVaultClientTimeout 	= "VAULT_CLIENT_TIMEOUT"
	const EnvVaultMFA 				= "VAULT_MFA"
	const EnvVaultMaxRetries 		= "VAULT_MAX_RETRIES"
	const EnvVaultNamespace 		= "VAULT_NAMESPACE"
	const EnvVaultSkipVerify 		= "VAULT_SKIP_VERIFY"
	const EnvVaultTLSServerName 	= "VAULT_TLS_SERVER_NAME"
	const EnvVaultToken 			= "VAULT_TOKEN"
	const EnvVaultWrapTTL 			= "VAULT_WRAP_TTL"
*/

var (
	vaultSecretPathFlag		= "vault-secret-path"
	vaultSecretPathViperKey	= "vault.secret.path"
	vaultSecretPathDefault	= "/secret"
	vaultSecretPathEnv 		= "VAULT_SECRET_PATH"
)

// InitFlags register flags for hashicorp vault
func InitFlags(f *pflag.FlagSet) {
	VaultSecretPath(f)
}

// VaultSecretPath registers a flag for the path used by vault secret engine
func VaultSecretPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path
Environment variable: %q`, vaultSecretPathEnv)
	f.String(vaultSecretPathFlag, vaultSecretPathDefault, desc)
	viper.BindPFlag(vaultSecretPathViperKey, f.Lookup(vaultSecretPathFlag))
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