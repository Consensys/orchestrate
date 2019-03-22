package secretstore

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	vault "github.com/hashicorp/vault/api"
)

func init() {
	viper.SetDefault(vaultURIViperKey, vaultURIDefault)
	viper.BindEnv(vaultURIViperKey, vaultURIEnv)
	viper.SetDefault(vaultTokenNameViperKey, vaultTokenNameEnv)
	viper.BindEnv(vaultTokenNameViperKey, vaultTokenNameEnv)
}

var (
	vaultURIFlag     = "vault-uri"
	vaultURIViperKey = "vault.uri"
	vaultURIDefault  = "http://127.0.0.1:8200"
	vaultURIEnv      = "VAULT_URI"

	vaultTokenNameFlag     = "vault-token-name"
	vaultTokenNameViperKey = "vault.token.name"
	vaultTokenNameDefault  = ""
	vaultTokenNameEnv      = "VAULT_TOKEN_NAME"

	vaultTokenFlag     = "vault-token"
	vaultTokenViperKey = "vault.token.value"
	vaultTokenDefault  = ""
	vaultTokenEnv      = "VAULT_TOKEN"

	vaultUnsealKeyFlag     = "vault-unseal-key"
	vaultUnsealKeyViperKey = "vault.unseal.key"
	vaultUnsealKeyDefault  = ""
	vaultUnsealKeyEnv      = "VAULT_UNSEAL_KEY"
)

// InitFlags register flags for hashicorp vault
func InitFlags(f *pflag.FlagSet) {
	VaultURI(f)
	VaultToken(f)
	VaultUnsealKey(f)
	VaultTokenName(f)
}

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

// NewConfig icreates vault configuration from viper
func NewConfig() *vault.Config {
	config := vault.DefaultConfig()
	config.Address = viper.GetString(vaultURIViperKey)
	return config
}

// AutoInit will try to Init the vault directly or FetchFromAws
func AutoInit(hashicorps *Hashicorps) (err error) {
	tokenName := viper.GetString("vault.token.name")
	awsSS := NewAWS(7)
	err = hashicorps.InitVault()
	if err != nil {
		// Probably Vault is already unsealed so we retrieve credentials from AWS
		err = hashicorps.InitFromAWS(awsSS, tokenName)
		if err != nil {
			return fmt.Errorf("Could not retrieve credentials from AWS: %v", err)
		}
	} else {
		// Vault has been properly unsealed so we push credentials on AWS
		err = hashicorps.SendToCredStore(awsSS, tokenName)
		if err != nil {
			return fmt.Errorf("Could not send credentials to AWS: %v", err)
		}

	}
	return nil
}
