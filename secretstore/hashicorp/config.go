package hashicorp

import (
	"fmt"

	vault "github.com/hashicorp/vault/api"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(vaultURIViperKey, vaultURIDefault)
	viper.BindEnv(vaultURIViperKey, vaultURIEnv)
}

var (
	vaultURIFlag     		= "vault-uri"
	vaultURIViperKey 		= "vault.uri"
	vaultURIDefault  		= "http://127.0.0.1:8200"
	vaultURIEnv      		= "VAULT_URI"

	vaultSecretPathFlag		= "vault-secret-path"
	vaultSecretPathViperKey	= "vault.secret.path"
	vaultSecretPathDefault	= "/secret/orchestra"
	vaultSecretPathEnv 		= "VAULT_SECRET_PATH"
)

// InitFlags register flags for hashicorp vault
func InitFlags(f *pflag.FlagSet) {
	VaultURI(f)
	VaultSecretPath(f)
}

// VaultURI register a flag for vault server address
func VaultURI(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret vault URI 
Environment variable: %q`, vaultURIEnv)
	f.String(vaultURIFlag, vaultURIDefault, desc)
	viper.BindPFlag(vaultURIViperKey, f.Lookup(vaultURIFlag))
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
	config.Address = viper.GetString(vaultURIViperKey)
	return config
}

// GetSecretPath returns the secret path set in deployment by vault
func GetSecretPath() string {
	return viper.GetString(vaultSecretPathViperKey)
}


