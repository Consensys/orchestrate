package secretstore

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	vault "github.com/hashicorp/vault/api"
)

func init() {
	viper.SetDefault("vault.uri", "https://127.0.0.1:8200")
	viper.SetDefault("vault.tokenName", "NO TOKEN NAME SPECIFIED")
}

var (
	vaultURIFlag     = "vault-uri"
	vaultURIViperKey = "vault.uri"
	vaultURIDefault  = "https://127.0.0.1:8200"
	vaultURIEnv      = "VAULT_URI"

	vaultTokenNameFlag = "vault-token-name"
	vaultTokenNameViperKey = "vault.token.name"
	vaultTokenNameDefault = "NO TOKEN NAME SPECIFIED"
	vaultTokenNameEnv = "VAULT_TOKEN_NAME"

	vaultTokenFlag = "vault-token"
	vaultTokenViperKey = "vault.token.value"
	vaultTokenDefault = ""
	vaultTokenEnv = "VAULT_TOKEN"

	vaultUnsealKeyFlag = "vault-unseal-key"
	vaultUnsealKeyViperKey = "vault.unseal.key"
	vaultUnsealKeyDefault = ""
	vaultUnsealKeyEnv = "VAULT_UNSEAL_KEY"
)

// VaultURI register a flag for vault server address
func VaultURI(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret vault URI Environment variable: %q`, vaultURIEnv)
	f.String(vaultURIFlag, vaultURIDefault, desc)
	viper.SetDefault(vaultURIViperKey, vaultURIDefault)
	viper.BindPFlag(vaultURIViperKey, f.Lookup(vaultURIFlag))
	viper.BindEnv(vaultURIViperKey, vaultURIEnv)
}

// VaultToken register a flag for vault server address
func VaultToken(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret vault Token Environment variable: %q`, vaultTokenEnv)
	f.String(vaultTokenFlag, vaultTokenDefault, desc)
	viper.SetDefault(vaultTokenViperKey, vaultTokenDefault)
	viper.BindPFlag(vaultTokenViperKey, f.Lookup(vaultTokenFlag))
	viper.BindEnv(vaultTokenViperKey, vaultTokenEnv)
}

// VaultUnsealKey registers a flag for the value of the vault unseal key
func VaultUnsealKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret vault unseal key Environment variable: %q`, vaultTokenEnv)
	f.String(vaultUnsealKeyFlag, vaultUnsealKeyDefault, desc)
	viper.SetDefault(vaultUnsealKeyViperKey, vaultUnsealKeyDefault)
	viper.BindPFlag(vaultUnsealKeyViperKey, f.Lookup(vaultUnsealKeyFlag))
	viper.BindEnv(vaultUnsealKeyViperKey, vaultUnsealKeyEnv)

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

// VaultTokenNameFromViper imports the vault token secret name on AWS
func VaultTokenNameFromViper() string {
	return viper.GetString("vault.token.name")
}

// VaultTokenFromViper imports the vault token from viper
func VaultTokenFromViper() string {
	return viper.GetString("vault.token.value")
}

// VaultUnsealKeyFromViper imports the vault unseal key from viper
func VaultUnsealKeyFromViper() string {
	return viper.GetString("vault.unseal.key")
}
