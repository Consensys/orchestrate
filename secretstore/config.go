package secretstore

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	vault "github.com/hashicorp/vault/api"
)

func init() {
	viper.SetDefault("vault.uri", "https://127.0.0.1:8200")
	viper.SetDefault("vault.tokenName", "YOU FORGOT TO SPECIFY THE NAME OF THE VAULT TOKEN ON AWS IN YOUR CONFIG")
}

var (
	vaultURIFlag     = "vault-uri"
	vaultURIViperKey = "vault.uri"
	vaultURIDefault  = "https://127.0.0.1:8200"
	vaultURIEnv      = "VAULT_URI"

	vaultTokenNameFlag = "vault-token-name"
	vaultTokenNameViperKey = "vault.token.name"
	vaultTokenNameDefault = "YOU FORGOT TO SPECIFY THE NAME OF THE VAULT TOKEN ON AWS IN YOUR CONFIG"
	vaultTokenNameEnv = "VAULT_TOKEN_NAME"
)

// VaultURI register a flag for vault server address
func VaultURI(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret vault URI Environment variable: %q`, vaultURIEnv)
	f.String(vaultURIFlag, vaultURIDefault, desc)
	viper.SetDefault(vaultURIViperKey, vaultURIDefault)
	viper.BindPFlag(vaultURIViperKey, f.Lookup(vaultURIFlag))
	viper.BindEnv(vaultURIViperKey, vaultURIEnv)
}

// VaultTokenName register a flag for vault server address
func VaultTokenName(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret vault URI Environment variable: %q`, vaultTokenNameEnv)
	f.String(vaultTokenNameFlag, vaultTokenNameDefault, desc)
	viper.SetDefault(vaultTokenNameViperKey, vaultTokenNameDefault)
	viper.BindPFlag(vaultTokenNameViperKey, f.Lookup(vaultURIFlag))
	viper.BindEnv(vaultTokenNameViperKey, vaultTokenNameEnv)
}

// VaultConfigFromViper import vault configuration from viper
func VaultConfigFromViper() *vault.Config {
	config := vault.DefaultConfig()
	config.Address = viper.GetString("vault.uri")
	return config
}

// VaultTokenFromViper imports the vault token secret name on AWS
func VaultTokenFromViper() string {
	return viper.GetString("vault.token.name")
}
