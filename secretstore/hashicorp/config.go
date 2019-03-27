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
	vaultURIFlag     = "vault-uri"
	vaultURIViperKey = "vault.uri"
	vaultURIDefault  = "http://127.0.0.1:8200"
	vaultURIEnv      = "VAULT_URI"
)

// InitFlags register flags for hashicorp vault
func InitFlags(f *pflag.FlagSet) {
	VaultURI(f)
}

// VaultURI register a flag for vault server address
func VaultURI(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret vault URI 
Environment variable: %q`, vaultURIEnv)
	f.String(vaultURIFlag, vaultURIDefault, desc)
	viper.BindPFlag(vaultURIViperKey, f.Lookup(vaultURIFlag))
}

// NewConfig icreates vault configuration from viper
func NewConfig() *vault.Config {
	config := vault.DefaultConfig()
	config.Address = viper.GetString(vaultURIViperKey)
	return config
}


